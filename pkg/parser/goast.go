// Package parser provides source code parsers for Go types, generator.yaml,
// and Terraform documentation files.
package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strings"

	"github.com/aws-controllers-k8s/ack-scanner/pkg/types"
)

// heuristicSubstrings are the case-insensitive substrings used to identify
// document-string fields. A field name must contain at least one of these
// to be considered a candidate for is_document or is_iam_policy annotation.
var heuristicSubstrings = []string{
	"policy",
	"document",
	"configuration",
	"template",
	"schema",
	"definition",
}

// GoASTParser parses Go source files to extract *string fields.
type GoASTParser struct{}

// ParseTypesFile parses a types.go file and returns all *string fields
// matching the document-string heuristics.
func (p *GoASTParser) ParseTypesFile(filePath string) ([]types.GoField, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var fields []types.GoField

	ast.Inspect(file, func(n ast.Node) bool {
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return true
		}

		structName := typeSpec.Name.Name

		for _, field := range structType.Fields.List {
			// Skip anonymous/embedded fields
			if len(field.Names) == 0 {
				continue
			}

			// Skip unexported fields
			fieldName := field.Names[0].Name
			if !ast.IsExported(fieldName) {
				continue
			}

			// Check if field type is *string
			if !isPointerToString(field.Type) {
				continue
			}

			// Check if field name matches heuristic
			if !p.MatchesHeuristic(fieldName) {
				continue
			}

			jsonTag := extractJSONTag(field.Tag)

			fields = append(fields, types.GoField{
				StructName: structName,
				FieldName:  fieldName,
				FieldPath:  structName + "." + fieldName,
				GoType:     "*string",
				JSONTag:    jsonTag,
			})
		}

		return true
	})

	return fields, nil
}

// MatchesHeuristic returns true if a field name matches document-string patterns.
// It performs case-insensitive substring matching against the known heuristic list:
// "Policy", "Document", "Configuration", "Template", "Schema", "Definition".
func (p *GoASTParser) MatchesHeuristic(fieldName string) bool {
	lower := strings.ToLower(fieldName)
	for _, substr := range heuristicSubstrings {
		if strings.Contains(lower, substr) {
			return true
		}
	}
	return false
}

// isPointerToString checks whether an AST expression represents the type *string.
func isPointerToString(expr ast.Expr) bool {
	starExpr, ok := expr.(*ast.StarExpr)
	if !ok {
		return false
	}
	ident, ok := starExpr.X.(*ast.Ident)
	if !ok {
		return false
	}
	return ident.Name == "string"
}

// extractJSONTag extracts the JSON tag value from a struct field tag.
// For example, given `json:"policy_document,omitempty"`, it returns "policy_document".
// Returns an empty string if no JSON tag is present.
func extractJSONTag(tag *ast.BasicLit) string {
	if tag == nil {
		return ""
	}

	// tag.Value includes the backticks, e.g. `json:"name,omitempty"`
	tagValue := tag.Value
	if len(tagValue) < 2 {
		return ""
	}
	// Strip the backticks
	tagValue = tagValue[1 : len(tagValue)-1]

	// Use reflect.StructTag to parse
	st := reflect.StructTag(tagValue)
	jsonTag := st.Get("json")
	if jsonTag == "" || jsonTag == "-" {
		return ""
	}

	// Strip options like ",omitempty"
	if idx := strings.Index(jsonTag, ","); idx != -1 {
		jsonTag = jsonTag[:idx]
	}

	return jsonTag
}
