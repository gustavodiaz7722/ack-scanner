package parser

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/aws-controllers-k8s/ack-scanner/pkg/types"
)

// argumentEntry represents a parsed field from the Terraform Argument Reference section.
type argumentEntry struct {
	FieldName   string
	Description string
	Required    bool
}

// TerraformParser parses Terraform provider documentation.
// The docs live in hashicorp/terraform-provider-aws under website/docs/r/
// with ~1,666 files named {service}_{resource}.html.markdown.
type TerraformParser struct{}

// jsonPhrases are the case-insensitive phrases that indicate a JSON field.
var jsonPhrases = []string{
	"json-encoded",
	"json string",
	"json formatted",
	"json schema",
	"policy document",
	"json",
}

// argPattern matches lines like: * `field_name` - (Required) description...
// or: * `field_name` - (Optional) description...
var argPattern = regexp.MustCompile("^\\*\\s+`([^`]+)`\\s+-\\s+\\((Required|Optional)\\)\\s+(.*)")

// jsonencodePattern matches lines like: field_name = jsonencode(
var jsonencodePattern = regexp.MustCompile(`^\s*(\w+)\s*=\s*jsonencode\s*\(`)

// policyDocPattern matches lines like: field_name = data.aws_iam_policy_document.foo.json
var policyDocPattern = regexp.MustCompile(`^\s*(\w+)\s*=\s*data\.aws_iam_policy_document\.\w+\.json`)

// ExtractArguments parses the "## Argument Reference" section and returns
// field entries matching the pattern: * `field_name` - (Required|Optional) description
func (p *TerraformParser) ExtractArguments(content string) []argumentEntry {
	var entries []argumentEntry

	inArgSection := false
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()

		// Detect start of Argument Reference section
		if strings.HasPrefix(line, "## Argument Reference") {
			inArgSection = true
			continue
		}

		// Stop at next ## section
		if inArgSection && strings.HasPrefix(line, "## ") {
			break
		}

		if !inArgSection {
			continue
		}

		matches := argPattern.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		entry := argumentEntry{
			FieldName:   matches[1],
			Description: matches[3],
			Required:    matches[2] == "Required",
		}
		entries = append(entries, entry)
	}

	return entries
}

// resourceBlockPattern matches HCL resource block declarations like:
//
//	resource "aws_eks_addon" "example" {
var resourceBlockPattern = regexp.MustCompile(`^\s*resource\s+"([^"]+)"\s+"[^"]+"\s*\{`)

// ExtractJsonEncodeFields parses "## Example Usage" code blocks and returns
// field names assigned via jsonencode(...) or data.aws_iam_policy_document.*.json.
// It only extracts fields from resource blocks matching the expected resource type
// (derived from the doc filename), preventing false positives from supporting
// resources like aws_iam_role that commonly appear in examples.
func (p *TerraformParser) ExtractJsonEncodeFields(content string, expectedResourceType string) []string {
	fieldSet := make(map[string]struct{})

	inExampleSection := false
	inCodeBlock := false
	inMatchingResource := false
	braceDepth := 0
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()

		// Detect start of Example Usage section
		if strings.HasPrefix(line, "## Example Usage") {
			inExampleSection = true
			continue
		}

		// Stop at next ## section (but not sub-sections like ### within Example Usage)
		if inExampleSection && strings.HasPrefix(line, "## ") && !strings.HasPrefix(line, "## Example Usage") {
			break
		}

		if !inExampleSection {
			continue
		}

		// Track code block boundaries
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") {
			inCodeBlock = !inCodeBlock
			if !inCodeBlock {
				// Reset resource tracking when code block ends
				inMatchingResource = false
				braceDepth = 0
			}
			continue
		}

		if !inCodeBlock {
			continue
		}

		// Check for resource block declaration
		if matches := resourceBlockPattern.FindStringSubmatch(line); matches != nil {
			resourceType := matches[1] // e.g., "aws_eks_addon"
			inMatchingResource = (resourceType == expectedResourceType)
			braceDepth = 1
			continue
		}

		// Track brace depth to know when we exit a resource block
		if inMatchingResource || braceDepth > 0 {
			for _, ch := range line {
				if ch == '{' {
					braceDepth++
				} else if ch == '}' {
					braceDepth--
					if braceDepth <= 0 {
						inMatchingResource = false
						braceDepth = 0
						break
					}
				}
			}
		}

		// Only extract from the matching resource block
		if !inMatchingResource {
			continue
		}

		// Check for jsonencode pattern
		if matches := jsonencodePattern.FindStringSubmatch(line); matches != nil {
			fieldSet[matches[1]] = struct{}{}
		}

		// Check for data.aws_iam_policy_document.*.json pattern
		if matches := policyDocPattern.FindStringSubmatch(line); matches != nil {
			fieldSet[matches[1]] = struct{}{}
		}
	}

	fields := make([]string, 0, len(fieldSet))
	for f := range fieldSet {
		fields = append(fields, f)
	}
	return fields
}

// containsJSONPhrase returns true if the description contains any of the
// case-insensitive JSON indicator phrases.
func containsJSONPhrase(description string) bool {
	lower := strings.ToLower(description)
	for _, phrase := range jsonPhrases {
		if strings.Contains(lower, phrase) {
			return true
		}
	}
	return false
}

// ParseResourceDoc parses a single Terraform resource markdown file using
// two detection signals:
// 1. Description phrases in ## Argument Reference (e.g., "JSON String", "JSON formatted")
// 2. jsonencode() usage in ## Example Usage code blocks
func (p *TerraformParser) ParseResourceDoc(filePath string) ([]types.TerraformField, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading terraform doc %s: %w", filePath, err)
	}
	content := string(data)

	fileName := filepath.Base(filePath)
	service := ServiceFromFileName(fileName)
	resource := ResourceFromFileName(fileName)

	// Extract arguments from Argument Reference section
	args := p.ExtractArguments(content)

	// Extract fields used with jsonencode in examples — only from the
	// resource block matching this doc's resource type
	expectedResourceType := "aws_" + service + "_" + resource
	jsonencodeFields := p.ExtractJsonEncodeFields(content, expectedResourceType)
	jsonencodeSet := make(map[string]struct{}, len(jsonencodeFields))
	for _, f := range jsonencodeFields {
		jsonencodeSet[f] = struct{}{}
	}

	// Build result by combining both signals
	fieldMap := make(map[string]*types.TerraformField)

	// Check description phrases in arguments
	for _, arg := range args {
		if containsJSONPhrase(arg.Description) {
			fieldMap[arg.FieldName] = &types.TerraformField{
				ServiceName:     service,
				ResourceType:    resource,
				FieldName:       arg.FieldName,
				Description:     arg.Description,
				DetectionMethod: types.DetectDescriptionPhrase,
			}
		}
	}

	// Check jsonencode usage in examples
	for _, fieldName := range jsonencodeFields {
		if existing, ok := fieldMap[fieldName]; ok {
			// Found in both description and jsonencode
			existing.DetectionMethod = types.DetectBoth
		} else {
			fieldMap[fieldName] = &types.TerraformField{
				ServiceName:     service,
				ResourceType:    resource,
				FieldName:       fieldName,
				Description:     descriptionForField(args, fieldName),
				DetectionMethod: types.DetectJsonEncodeExample,
			}
		}
	}

	results := make([]types.TerraformField, 0, len(fieldMap))
	for _, field := range fieldMap {
		results = append(results, *field)
	}
	return results, nil
}

// descriptionForField looks up the description for a field name from the argument list.
func descriptionForField(args []argumentEntry, fieldName string) string {
	for _, arg := range args {
		if arg.FieldName == fieldName {
			return arg.Description
		}
	}
	return ""
}

// ParseAllDocs parses all *.html.markdown files in a directory.
func (p *TerraformParser) ParseAllDocs(docsDir string) ([]types.TerraformField, error) {
	pattern := filepath.Join(docsDir, "*.html.markdown")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("globbing terraform docs in %s: %w", docsDir, err)
	}

	var allFields []types.TerraformField
	for _, file := range files {
		fields, err := p.ParseResourceDoc(file)
		if err != nil {
			// Log warning and continue processing remaining files
			log.Printf("Warning: skipping malformed terraform doc %s: %v", filepath.Base(file), err)
			continue
		}
		allFields = append(allFields, fields...)
	}

	return allFields, nil
}

// ServiceFromFileName extracts the service name from a Terraform resource doc file.
// For example, "sns_topic_subscription.html.markdown" → "sns".
func ServiceFromFileName(fileName string) string {
	// Strip the .html.markdown extension
	name := strings.TrimSuffix(fileName, ".html.markdown")
	// Also handle .markdown extension
	name = strings.TrimSuffix(name, ".markdown")

	// Split on first underscore to get service prefix
	idx := strings.Index(name, "_")
	if idx < 0 {
		return name
	}
	return name[:idx]
}

// ResourceFromFileName extracts the resource type from a Terraform resource doc file.
// For example, "sns_topic_subscription.html.markdown" → "topic_subscription".
func ResourceFromFileName(fileName string) string {
	// Strip the .html.markdown extension
	name := strings.TrimSuffix(fileName, ".html.markdown")
	// Also handle .markdown extension
	name = strings.TrimSuffix(name, ".markdown")

	// Take everything after the first underscore
	idx := strings.Index(name, "_")
	if idx < 0 {
		return ""
	}
	return name[idx+1:]
}
