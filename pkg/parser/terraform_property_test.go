package parser

import (
	"fmt"
	"strings"
	"testing"

	"pgregory.net/rapid"
)

// Feature: ack-scanner, Property 11: Terraform JSON phrase detection
// **Validates: Requirements 4.2, 4.3**

// jsonPhrasesForTest are the phrases that should trigger JSON field detection.
var jsonPhrasesForTest = []string{
	"JSON",
	"JSON-encoded",
	"JSON string",
	"JSON formatted",
	"JSON schema",
	"policy document",
}

// genSafeWord generates a word that does not contain any JSON indicator phrases.
// It uses lowercase letters and avoids "json", "policy", and "document" substrings.
func genSafeWord() *rapid.Generator[string] {
	return rapid.Custom(func(t *rapid.T) string {
		// Use a limited alphabet that cannot spell "json", "policy", or "document"
		// by using only letters: a, b, c, f, g, h, k, r, w, x, z
		safeChars := []byte("abcfghkrwxz")
		length := rapid.IntRange(1, 10).Draw(t, "wordLen")
		var sb strings.Builder
		for i := 0; i < length; i++ {
			idx := rapid.IntRange(0, len(safeChars)-1).Draw(t, "charIdx")
			sb.WriteByte(safeChars[idx])
		}
		return sb.String()
	})
}

// genSafeDescription generates a description string that does NOT contain
// any of the JSON indicator phrases (case-insensitive).
func genSafeDescription() *rapid.Generator[string] {
	return rapid.Custom(func(t *rapid.T) string {
		wordCount := rapid.IntRange(1, 8).Draw(t, "wordCount")
		words := make([]string, wordCount)
		for i := 0; i < wordCount; i++ {
			words[i] = genSafeWord().Draw(t, "safeWord")
		}
		return strings.Join(words, " ")
	})
}

// genDescriptionWithPhrase generates a description that contains exactly one
// of the JSON phrases injected at a random position within the string.
func genDescriptionWithPhrase() *rapid.Generator[string] {
	return rapid.Custom(func(t *rapid.T) string {
		// Pick a phrase to inject
		phraseIdx := rapid.IntRange(0, len(jsonPhrasesForTest)-1).Draw(t, "phraseIdx")
		phrase := jsonPhrasesForTest[phraseIdx]

		// Optionally randomize casing of the phrase
		useMixed := rapid.Bool().Draw(t, "mixedCase")
		if useMixed {
			phrase = strings.ToUpper(phrase[:1]) + strings.ToLower(phrase[1:])
		}

		// Generate prefix and suffix words
		prefixCount := rapid.IntRange(0, 4).Draw(t, "prefixCount")
		suffixCount := rapid.IntRange(0, 4).Draw(t, "suffixCount")

		var parts []string
		for i := 0; i < prefixCount; i++ {
			parts = append(parts, genSafeWord().Draw(t, "prefix"))
		}
		parts = append(parts, phrase)
		for i := 0; i < suffixCount; i++ {
			parts = append(parts, genSafeWord().Draw(t, "suffix"))
		}
		return strings.Join(parts, " ")
	})
}

// TestContainsJSONPhrase_TrueWhenPhrasePresent verifies that containsJSONPhrase
// returns true when the description contains at least one of the JSON phrases.
func TestContainsJSONPhrase_TrueWhenPhrasePresent(t *testing.T) {
	// Feature: ack-scanner, Property 11: Terraform JSON phrase detection
	rapid.Check(t, func(t *rapid.T) {
		desc := genDescriptionWithPhrase().Draw(t, "description")
		if !containsJSONPhrase(desc) {
			t.Fatalf("containsJSONPhrase(%q) = false, want true (description contains a JSON phrase)", desc)
		}
	})
}

// TestContainsJSONPhrase_FalseWhenNoPhrase verifies that containsJSONPhrase
// returns false when the description does not contain any of the JSON phrases.
func TestContainsJSONPhrase_FalseWhenNoPhrase(t *testing.T) {
	// Feature: ack-scanner, Property 11: Terraform JSON phrase detection
	rapid.Check(t, func(t *rapid.T) {
		desc := genSafeDescription().Draw(t, "description")
		if containsJSONPhrase(desc) {
			t.Fatalf("containsJSONPhrase(%q) = true, want false (description has no JSON phrase)", desc)
		}
	})
}

// TestContainsJSONPhrase_CaseInsensitive verifies that phrase detection is
// case-insensitive by testing all-uppercase, all-lowercase, and mixed-case variants.
func TestContainsJSONPhrase_CaseInsensitive(t *testing.T) {
	// Feature: ack-scanner, Property 11: Terraform JSON phrase detection
	rapid.Check(t, func(t *rapid.T) {
		phraseIdx := rapid.IntRange(0, len(jsonPhrasesForTest)-1).Draw(t, "phraseIdx")
		phrase := jsonPhrasesForTest[phraseIdx]

		// Test all-uppercase
		upper := strings.ToUpper(phrase)
		desc := "the field accepts " + upper + " input"
		if !containsJSONPhrase(desc) {
			t.Fatalf("containsJSONPhrase(%q) = false, want true for uppercase phrase %q", desc, upper)
		}

		// Test all-lowercase
		lower := strings.ToLower(phrase)
		desc = "the field accepts " + lower + " input"
		if !containsJSONPhrase(desc) {
			t.Fatalf("containsJSONPhrase(%q) = false, want true for lowercase phrase %q", desc, lower)
		}
	})
}

// Feature: ack-scanner, Property 11b: Terraform jsonencode detection
// **Validates: Requirements 4.4**

// genSnakeCaseFieldName generates a valid snake_case field name (lowercase letters
// and underscores, must start with a letter, 1-3 segments).
func genSnakeCaseFieldName() *rapid.Generator[string] {
	return rapid.Custom(func(t *rapid.T) string {
		segCount := rapid.IntRange(1, 3).Draw(t, "segCount")
		segments := make([]string, segCount)
		for i := 0; i < segCount; i++ {
			length := rapid.IntRange(2, 8).Draw(t, "segLen")
			var sb strings.Builder
			for j := 0; j < length; j++ {
				ch := rapid.ByteRange('a', 'z').Draw(t, "ch")
				sb.WriteByte(ch)
			}
			segments[i] = sb.String()
		}
		return strings.Join(segments, "_")
	})
}

// TestExtractJsonEncodeFields_FindsJsonencode verifies that ExtractJsonEncodeFields
// correctly extracts field names from code blocks containing jsonencode assignments.
func TestExtractJsonEncodeFields_FindsJsonencode(t *testing.T) {
	// Feature: ack-scanner, Property 11b: Terraform jsonencode detection
	rapid.Check(t, func(t *rapid.T) {
		fieldName := genSnakeCaseFieldName().Draw(t, "fieldName")

		// Construct a Terraform doc with an Example Usage section containing jsonencode
		// The resource type matches the expected type so it gets picked up
		content := fmt.Sprintf(`## Example Usage

`+"```"+`hcl
resource "aws_test_resource" "example" {
  name = "test"
  %s = jsonencode({
    key = "value"
  })
}
`+"```"+`

## Argument Reference
`, fieldName)

		parser := &TerraformParser{}
		fields := parser.ExtractJsonEncodeFields(content, "aws_test_resource")

		found := false
		for _, f := range fields {
			if f == fieldName {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("ExtractJsonEncodeFields did not find field %q in content with jsonencode assignment", fieldName)
		}
	})
}

// TestExtractJsonEncodeFields_FindsPolicyDocRef verifies that ExtractJsonEncodeFields
// correctly extracts field names from code blocks containing data.aws_iam_policy_document references.
func TestExtractJsonEncodeFields_FindsPolicyDocRef(t *testing.T) {
	// Feature: ack-scanner, Property 11b: Terraform jsonencode detection
	rapid.Check(t, func(t *rapid.T) {
		fieldName := genSnakeCaseFieldName().Draw(t, "fieldName")

		// Construct a Terraform doc with a policy document reference
		content := fmt.Sprintf(`## Example Usage

`+"```"+`hcl
resource "aws_test_resource" "example" {
  name = "test"
  %s = data.aws_iam_policy_document.example.json
}
`+"```"+`

## Argument Reference
`, fieldName)

		parser := &TerraformParser{}
		fields := parser.ExtractJsonEncodeFields(content, "aws_test_resource")

		found := false
		for _, f := range fields {
			if f == fieldName {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("ExtractJsonEncodeFields did not find field %q in content with policy document reference", fieldName)
		}
	})
}

// TestExtractJsonEncodeFields_DoesNotMatchOutsideCodeBlocks verifies that
// jsonencode patterns outside of code blocks are not detected.
func TestExtractJsonEncodeFields_DoesNotMatchOutsideCodeBlocks(t *testing.T) {
	// Feature: ack-scanner, Property 11b: Terraform jsonencode detection
	rapid.Check(t, func(t *rapid.T) {
		fieldName := genSnakeCaseFieldName().Draw(t, "fieldName")

		// jsonencode outside code blocks should not be detected
		content := fmt.Sprintf(`## Example Usage

The %s field uses jsonencode to set the value.

%s = jsonencode({})

## Argument Reference
`, fieldName, fieldName)

		parser := &TerraformParser{}
		fields := parser.ExtractJsonEncodeFields(content, "aws_test_resource")

		for _, f := range fields {
			if f == fieldName {
				t.Fatalf("ExtractJsonEncodeFields incorrectly found field %q outside of a code block", fieldName)
			}
		}
	})
}

// Feature: ack-scanner, Property 12: Terraform file name service extraction
// **Validates: Requirements 4.6**

// genServiceName generates a valid service name: lowercase letters only, no underscores,
// length 2-12, matching AWS service prefixes like "sns", "iam", "eks".
func genServiceName() *rapid.Generator[string] {
	return rapid.Custom(func(t *rapid.T) string {
		length := rapid.IntRange(2, 12).Draw(t, "serviceLen")
		var sb strings.Builder
		for i := 0; i < length; i++ {
			ch := rapid.ByteRange('a', 'z').Draw(t, "ch")
			sb.WriteByte(ch)
		}
		return sb.String()
	})
}

// genTFResourceName generates a Terraform-style resource name: lowercase letters
// with underscores allowed, must start with a letter, can have 1-3 segments.
func genTFResourceName() *rapid.Generator[string] {
	return rapid.Custom(func(t *rapid.T) string {
		segCount := rapid.IntRange(1, 3).Draw(t, "segCount")
		segments := make([]string, segCount)
		for i := 0; i < segCount; i++ {
			length := rapid.IntRange(2, 10).Draw(t, "segLen")
			var sb strings.Builder
			for j := 0; j < length; j++ {
				ch := rapid.ByteRange('a', 'z').Draw(t, "ch")
				sb.WriteByte(ch)
			}
			segments[i] = sb.String()
		}
		return strings.Join(segments, "_")
	})
}

// TestServiceFromFileName_ExtractsServicePrefix verifies that ServiceFromFileName
// correctly extracts the service name as the prefix before the first underscore.
func TestServiceFromFileName_ExtractsServicePrefix(t *testing.T) {
	// Feature: ack-scanner, Property 12: Terraform file name service extraction
	rapid.Check(t, func(t *rapid.T) {
		service := genServiceName().Draw(t, "service")
		resource := genTFResourceName().Draw(t, "resource")

		fileName := service + "_" + resource + ".html.markdown"
		got := ServiceFromFileName(fileName)

		if got != service {
			t.Fatalf("ServiceFromFileName(%q) = %q, want %q", fileName, got, service)
		}
	})
}

// TestResourceFromFileName_ExtractsResourceSuffix verifies that ResourceFromFileName
// correctly extracts the resource type as the remainder after the first underscore.
func TestResourceFromFileName_ExtractsResourceSuffix(t *testing.T) {
	// Feature: ack-scanner, Property 12: Terraform file name service extraction
	rapid.Check(t, func(t *rapid.T) {
		service := genServiceName().Draw(t, "service")
		resource := genTFResourceName().Draw(t, "resource")

		fileName := service + "_" + resource + ".html.markdown"
		got := ResourceFromFileName(fileName)

		if got != resource {
			t.Fatalf("ResourceFromFileName(%q) = %q, want %q", fileName, got, resource)
		}
	})
}

// TestFileNameExtraction_ConsistentPair verifies that for any generated filename,
// the service and resource parts together reconstruct the original base name
// (service + "_" + resource).
func TestFileNameExtraction_ConsistentPair(t *testing.T) {
	// Feature: ack-scanner, Property 12: Terraform file name service extraction
	rapid.Check(t, func(t *rapid.T) {
		service := genServiceName().Draw(t, "service")
		resource := genTFResourceName().Draw(t, "resource")

		fileName := service + "_" + resource + ".html.markdown"
		gotService := ServiceFromFileName(fileName)
		gotResource := ResourceFromFileName(fileName)

		reconstructed := gotService + "_" + gotResource
		expected := service + "_" + resource

		if reconstructed != expected {
			t.Fatalf("ServiceFromFileName+ResourceFromFileName(%q): got %q + %q = %q, want %q",
				fileName, gotService, gotResource, reconstructed, expected)
		}
	})
}

// TestServiceFromFileName_NoUnderscore verifies that when there is no underscore
// in the base name, ServiceFromFileName returns the full name without extension.
func TestServiceFromFileName_NoUnderscore(t *testing.T) {
	// Feature: ack-scanner, Property 12: Terraform file name service extraction
	rapid.Check(t, func(t *rapid.T) {
		service := genServiceName().Draw(t, "service")
		fileName := service + ".html.markdown"
		got := ServiceFromFileName(fileName)

		if got != service {
			t.Fatalf("ServiceFromFileName(%q) = %q, want %q (no underscore case)", fileName, got, service)
		}
	})
}

// TestResourceFromFileName_NoUnderscore verifies that when there is no underscore,
// ResourceFromFileName returns an empty string.
func TestResourceFromFileName_NoUnderscore(t *testing.T) {
	// Feature: ack-scanner, Property 12: Terraform file name service extraction
	rapid.Check(t, func(t *rapid.T) {
		service := genServiceName().Draw(t, "service")
		fileName := service + ".html.markdown"
		got := ResourceFromFileName(fileName)

		if got != "" {
			t.Fatalf("ResourceFromFileName(%q) = %q, want empty string (no underscore case)", fileName, got)
		}
	})
}
