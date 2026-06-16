package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aws-controllers-k8s/ack-scanner/pkg/types"
	"pgregory.net/rapid"
)

// Feature: ack-scanner, Property 6: Go parser heuristic field detection
// **Validates: Requirements 3.1, 3.3**

// heuristicKeywords are the keywords used by the parser to identify document-string fields.
var heuristicKeywords = []string{
	"Policy",
	"Document",
	"Configuration",
	"Template",
	"Schema",
	"Definition",
}

// nonHeuristicSuffixes are suffixes that do NOT contain any heuristic keyword.
var nonHeuristicSuffixes = []string{
	"Name",
	"ARN",
	"ID",
	"Endpoint",
	"Status",
	"Type",
	"Region",
	"Value",
	"Key",
	"Tag",
	"Owner",
	"Account",
	"Bucket",
	"Table",
	"Queue",
	"Cluster",
	"Instance",
	"Role",
	"Group",
	"User",
	"Action",
	"Resource",
}

// genUpperCamelWord generates a CamelCase word: one uppercase letter followed by 1-6 lowercase letters.
func genUpperCamelWord() *rapid.Generator[string] {
	return rapid.Custom(func(t *rapid.T) string {
		upper := rapid.ByteRange('A', 'Z').Draw(t, "upper")
		lowerCount := rapid.IntRange(1, 6).Draw(t, "lowerCount")
		var sb strings.Builder
		sb.WriteByte(upper)
		for i := 0; i < lowerCount; i++ {
			lower := rapid.ByteRange('a', 'z').Draw(t, "lower")
			sb.WriteByte(lower)
		}
		return sb.String()
	})
}

// genHeuristicFieldName generates a valid Go exported field name that contains
// one of the heuristic keywords (case-preserving).
func genHeuristicFieldName() *rapid.Generator[string] {
	return rapid.Custom(func(t *rapid.T) string {
		keyword := rapid.SampledFrom(heuristicKeywords).Draw(t, "keyword")
		// Optionally add a prefix word (e.g., "Assume" + "Role" + "Policy" + "Document")
		prefixCount := rapid.IntRange(0, 2).Draw(t, "prefixCount")
		var sb strings.Builder
		for i := 0; i < prefixCount; i++ {
			word := genUpperCamelWord().Draw(t, "prefix")
			sb.WriteString(word)
		}
		sb.WriteString(keyword)
		// Optionally add a suffix word
		suffixCount := rapid.IntRange(0, 2).Draw(t, "suffixCount")
		for i := 0; i < suffixCount; i++ {
			word := genUpperCamelWord().Draw(t, "suffix")
			sb.WriteString(word)
		}
		return sb.String()
	})
}

// genNonHeuristicFieldName generates a valid Go exported field name that does NOT
// contain any heuristic keyword (case-insensitive).
func genNonHeuristicFieldName() *rapid.Generator[string] {
	return rapid.Custom(func(t *rapid.T) string {
		// Use a safe suffix that doesn't contain any heuristic keyword
		suffix := rapid.SampledFrom(nonHeuristicSuffixes).Draw(t, "suffix")
		// Optionally add a prefix word
		prefixCount := rapid.IntRange(0, 2).Draw(t, "prefixCount")
		var sb strings.Builder
		for i := 0; i < prefixCount; i++ {
			word := genUpperCamelWord().Draw(t, "prefix")
			sb.WriteString(word)
		}
		sb.WriteString(suffix)
		result := sb.String()
		// Verify it doesn't accidentally contain a heuristic keyword
		lower := strings.ToLower(result)
		for _, kw := range heuristicKeywords {
			if strings.Contains(lower, strings.ToLower(kw)) {
				// Fallback to a safe name
				return "SafeFieldName"
			}
		}
		return result
	})
}

// genStructName generates a valid Go exported struct name.
func genStructName() *rapid.Generator[string] {
	return rapid.Custom(func(t *rapid.T) string {
		word := genUpperCamelWord().Draw(t, "structWord")
		suffix := rapid.SampledFrom([]string{"Spec", "Status", "Input", "Output", "Params"}).Draw(t, "structSuffix")
		return word + suffix
	})
}

// fieldDef represents a field definition in a generated struct.
type fieldDef struct {
	Name        string
	Type        string
	IsHeuristic bool // true if name contains a heuristic keyword
	IsPtrString bool // true if type is *string
	IsExported  bool // true if field starts with uppercase
}

// genFieldDefs generates a slice of field definitions with a controlled mix of:
// - *string fields that match heuristics (should be returned)
// - *string fields that don't match heuristics (should NOT be returned)
// - non-*string fields (should NOT be returned)
func genFieldDefs() *rapid.Generator[[]fieldDef] {
	return rapid.Custom(func(t *rapid.T) []fieldDef {
		// At least 1 field, up to 10
		count := rapid.IntRange(1, 10).Draw(t, "fieldCount")
		fields := make([]fieldDef, 0, count)
		usedNames := make(map[string]bool)

		for i := 0; i < count; i++ {
			var fd fieldDef
			// Choose what kind of field to generate
			kind := rapid.IntRange(0, 3).Draw(t, "fieldKind")
			switch kind {
			case 0:
				// *string field matching heuristic (should be returned)
				fd.Name = genHeuristicFieldName().Draw(t, "heuristicField")
				fd.Type = "*string"
				fd.IsHeuristic = true
				fd.IsPtrString = true
				fd.IsExported = true
			case 1:
				// *string field NOT matching heuristic (should NOT be returned)
				fd.Name = genNonHeuristicFieldName().Draw(t, "nonHeuristicField")
				fd.Type = "*string"
				fd.IsHeuristic = false
				fd.IsPtrString = true
				fd.IsExported = true
			case 2:
				// Non-*string field matching heuristic (should NOT be returned)
				fd.Name = genHeuristicFieldName().Draw(t, "heuristicNonStr")
				nonStrType := rapid.SampledFrom([]string{"*int64", "*bool", "string", "[]*string", "*MyStruct"}).Draw(t, "nonStrType")
				fd.Type = nonStrType
				fd.IsHeuristic = true
				fd.IsPtrString = false
				fd.IsExported = true
			case 3:
				// Non-*string field NOT matching heuristic (should NOT be returned)
				fd.Name = genNonHeuristicFieldName().Draw(t, "nonHeuristicNonStr")
				nonStrType := rapid.SampledFrom([]string{"*int64", "*bool", "string", "[]*string", "*MyStruct"}).Draw(t, "nonStrType")
				fd.Type = nonStrType
				fd.IsHeuristic = false
				fd.IsPtrString = false
				fd.IsExported = true
			}
			// Avoid duplicate field names within the struct
			if usedNames[fd.Name] {
				fd.Name = fd.Name + fmt.Sprintf("V%d", i)
			}
			usedNames[fd.Name] = true
			fields = append(fields, fd)
		}
		return fields
	})
}

// buildGoSource constructs a syntactically valid Go source file from struct definitions.
func buildGoSource(structName string, fields []fieldDef) string {
	var sb strings.Builder
	sb.WriteString("package types\n\n")
	sb.WriteString(fmt.Sprintf("type %s struct {\n", structName))
	for _, f := range fields {
		if f.IsExported {
			sb.WriteString(fmt.Sprintf("\t%s %s\n", f.Name, f.Type))
		} else {
			sb.WriteString(fmt.Sprintf("\t%s %s\n", strings.ToLower(f.Name[:1])+f.Name[1:], f.Type))
		}
	}
	sb.WriteString("}\n")
	return sb.String()
}

// TestProperty6_ParserReturnsExactlyHeuristicPtrStringFields verifies that
// ParseTypesFile returns exactly those fields that are both typed *string AND
// have names matching the heuristic keywords — no more, no less.
func TestProperty6_ParserReturnsExactlyHeuristicPtrStringFields(t *testing.T) {
	// Feature: ack-scanner, Property 6: Go parser heuristic field detection
	rapid.Check(t, func(t *rapid.T) {
		structName := genStructName().Draw(t, "structName")
		fields := genFieldDefs().Draw(t, "fields")

		// Build Go source
		src := buildGoSource(structName, fields)

		// Write to temp file
		tmpDir := t.Name()
		dir, err := os.MkdirTemp("", tmpDir)
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(dir)

		filePath := filepath.Join(dir, "types.go")
		if err := os.WriteFile(filePath, []byte(src), 0644); err != nil {
			t.Fatal(err)
		}

		// Parse the file
		parser := &GoASTParser{}
		result, err := parser.ParseTypesFile(filePath)
		if err != nil {
			t.Fatalf("ParseTypesFile failed on valid Go source: %v\nSource:\n%s", err, src)
		}

		// Compute expected: fields that are *string AND exported AND match heuristic
		var expectedNames []string
		for _, f := range fields {
			if f.IsPtrString && f.IsExported && f.IsHeuristic {
				expectedNames = append(expectedNames, f.Name)
			}
		}

		// Verify: no false positives — every returned field must be *string and match heuristic
		returnedNames := make(map[string]bool)
		for _, rf := range result {
			returnedNames[rf.FieldName] = true

			// Check type is *string
			if rf.GoType != "*string" {
				t.Fatalf("Returned field %q has GoType=%q, expected *string", rf.FieldName, rf.GoType)
			}

			// Check it matches the heuristic
			if !parser.MatchesHeuristic(rf.FieldName) {
				t.Fatalf("Returned field %q does not match heuristic but was included", rf.FieldName)
			}

			// Check struct name matches
			if rf.StructName != structName {
				t.Fatalf("Returned field %q has StructName=%q, expected %q", rf.FieldName, rf.StructName, structName)
			}
		}

		// Verify: no false negatives — every expected field must be present
		for _, name := range expectedNames {
			if !returnedNames[name] {
				t.Fatalf("Expected field %q (matches heuristic and is *string) was not returned.\nSource:\n%s", name, src)
			}
		}

		// Verify count matches
		if len(result) != len(expectedNames) {
			t.Fatalf("Expected %d fields, got %d.\nExpected: %v\nGot: %v\nSource:\n%s",
				len(expectedNames), len(result), expectedNames, resultFieldNames(result), src)
		}
	})
}

// TestProperty6_MultipleStructs verifies the parser handles multiple structs
// in the same file correctly — returning heuristic-matching *string fields from
// all structs.
func TestProperty6_MultipleStructs(t *testing.T) {
	// Feature: ack-scanner, Property 6: Go parser heuristic field detection
	rapid.Check(t, func(t *rapid.T) {
		structCount := rapid.IntRange(1, 4).Draw(t, "structCount")

		var src strings.Builder
		src.WriteString("package types\n\n")

		type expectedField struct {
			structName string
			fieldName  string
		}
		var expected []expectedField
		usedStructNames := make(map[string]bool)

		for s := 0; s < structCount; s++ {
			structName := genStructName().Draw(t, "structName")
			// Ensure unique struct names
			if usedStructNames[structName] {
				structName = structName + fmt.Sprintf("%d", s)
			}
			usedStructNames[structName] = true

			fields := genFieldDefs().Draw(t, "fields")

			src.WriteString(fmt.Sprintf("type %s struct {\n", structName))
			for _, f := range fields {
				src.WriteString(fmt.Sprintf("\t%s %s\n", f.Name, f.Type))
				if f.IsPtrString && f.IsExported && f.IsHeuristic {
					expected = append(expected, expectedField{structName: structName, fieldName: f.Name})
				}
			}
			src.WriteString("}\n\n")
		}

		// Write to temp file
		dir, err := os.MkdirTemp("", "property6_multi")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(dir)

		filePath := filepath.Join(dir, "types.go")
		if err := os.WriteFile(filePath, []byte(src.String()), 0644); err != nil {
			t.Fatal(err)
		}

		// Parse the file
		parser := &GoASTParser{}
		result, err := parser.ParseTypesFile(filePath)
		if err != nil {
			t.Fatalf("ParseTypesFile failed: %v\nSource:\n%s", err, src.String())
		}

		// Build set of returned field paths
		returnedPaths := make(map[string]bool)
		for _, rf := range result {
			path := rf.StructName + "." + rf.FieldName
			returnedPaths[path] = true

			if rf.GoType != "*string" {
				t.Fatalf("Returned field %s has GoType=%q, expected *string", path, rf.GoType)
			}
			if !parser.MatchesHeuristic(rf.FieldName) {
				t.Fatalf("Returned field %s does not match heuristic", path)
			}
		}

		// Verify all expected fields are present
		for _, ef := range expected {
			path := ef.structName + "." + ef.fieldName
			if !returnedPaths[path] {
				t.Fatalf("Expected field %s was not returned.\nSource:\n%s", path, src.String())
			}
		}

		// Verify counts match
		if len(result) != len(expected) {
			t.Fatalf("Expected %d fields, got %d.\nSource:\n%s",
				len(expected), len(result), src.String())
		}
	})
}

// TestProperty6_MatchesHeuristicConsistentWithParser verifies that
// MatchesHeuristic and ParseTypesFile agree: a *string field is returned
// by the parser if and only if MatchesHeuristic returns true for its name.
func TestProperty6_MatchesHeuristicConsistentWithParser(t *testing.T) {
	// Feature: ack-scanner, Property 6: Go parser heuristic field detection
	rapid.Check(t, func(t *rapid.T) {
		// Generate a field name — either heuristic or not
		isHeuristic := rapid.Bool().Draw(t, "isHeuristic")
		var fieldName string
		if isHeuristic {
			fieldName = genHeuristicFieldName().Draw(t, "fieldName")
		} else {
			fieldName = genNonHeuristicFieldName().Draw(t, "fieldName")
		}

		structName := "TestStruct"
		src := fmt.Sprintf("package types\n\ntype %s struct {\n\t%s *string\n}\n", structName, fieldName)

		dir, err := os.MkdirTemp("", "property6_heuristic")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(dir)

		filePath := filepath.Join(dir, "types.go")
		if err := os.WriteFile(filePath, []byte(src), 0644); err != nil {
			t.Fatal(err)
		}

		parser := &GoASTParser{}
		result, err := parser.ParseTypesFile(filePath)
		if err != nil {
			t.Fatalf("ParseTypesFile failed: %v\nSource:\n%s", err, src)
		}

		heuristicResult := parser.MatchesHeuristic(fieldName)

		if heuristicResult && len(result) == 0 {
			t.Fatalf("MatchesHeuristic(%q) = true, but ParseTypesFile returned no fields.\nSource:\n%s",
				fieldName, src)
		}
		if !heuristicResult && len(result) > 0 {
			t.Fatalf("MatchesHeuristic(%q) = false, but ParseTypesFile returned %d fields.\nSource:\n%s",
				fieldName, len(result), src)
		}
	})
}

// resultFieldNames extracts field names from parsed results for error messages.
func resultFieldNames(results []types.GoField) []string {
	names := make([]string, len(results))
	for i, r := range results {
		names[i] = r.FieldName
	}
	return names
}
