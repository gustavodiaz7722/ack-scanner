package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/aws-controllers-k8s/ack-scanner/pkg/types"
)

// CRDParser parses Kubernetes CRD YAML files to extract string fields
// from the spec section. This is more reliable than parsing Go types.go
// because CRDs capture ALL fields including those defined in per-resource
// Go files (e.g., queue.go, subscription.go), not just types.go.
type CRDParser struct{}

// crdSchema represents the relevant portions of a CRD YAML file.
type crdSchema struct {
	Spec struct {
		Group string `yaml:"group"`
		Names struct {
			Kind string `yaml:"kind"`
		} `yaml:"names"`
		Versions []crdVersion `yaml:"versions"`
	} `yaml:"spec"`
}

type crdVersion struct {
	Name   string `yaml:"name"`
	Schema struct {
		OpenAPIV3Schema struct {
			Properties map[string]crdProperty `yaml:"properties"`
		} `yaml:"openAPIV3Schema"`
	} `yaml:"schema"`
}

type crdProperty struct {
	Type        string                 `yaml:"type"`
	Description string                 `yaml:"description"`
	Properties  map[string]crdProperty `yaml:"properties"`
}

// ParseCRDFile parses a single CRD YAML file and returns all string fields
// found under spec.properties that match the document-string heuristics.
func (p *CRDParser) ParseCRDFile(filePath string) ([]types.GoField, string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, "", fmt.Errorf("reading CRD file %s: %w", filePath, err)
	}

	var crd crdSchema
	if err := yaml.Unmarshal(data, &crd); err != nil {
		return nil, "", fmt.Errorf("parsing CRD YAML %s: %w", filePath, err)
	}

	resourceName := crd.Spec.Names.Kind
	if resourceName == "" {
		return nil, "", fmt.Errorf("CRD %s has no kind defined", filePath)
	}

	// Use the first version's schema (typically v1alpha1)
	if len(crd.Spec.Versions) == 0 {
		return nil, resourceName, nil
	}

	version := crd.Spec.Versions[0]
	topProps := version.Schema.OpenAPIV3Schema.Properties

	// Navigate to spec.properties
	specProp, ok := topProps["spec"]
	if !ok {
		return nil, resourceName, nil
	}

	specProps := specProp.Properties
	if specProps == nil {
		return nil, resourceName, nil
	}

	goParser := &GoASTParser{}
	var fields []types.GoField

	for fieldName, prop := range specProps {
		// Only grab direct string fields under spec (not nested objects)
		if prop.Type != "string" {
			continue
		}

		// Convert JSON field name to CamelCase for matching with generator.yaml
		camelName := jsonToCamel(fieldName)

		// Check if field name matches heuristic
		if !goParser.MatchesHeuristic(camelName) {
			continue
		}

		fields = append(fields, types.GoField{
			StructName: resourceName,
			FieldName:  camelName,
			FieldPath:  resourceName + "." + camelName,
			GoType:     "*string",
			JSONTag:    fieldName,
		})
	}

	return fields, resourceName, nil
}

// ParseAllCRDs parses all CRD YAML files in config/crd/bases/ for a given
// repo directory and returns the combined set of heuristic-matching string fields.
func (p *CRDParser) ParseAllCRDs(repoDir string) ([]types.GoField, error) {
	basesDir := filepath.Join(repoDir, "config", "crd", "bases")
	pattern := filepath.Join(basesDir, "*.yaml")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("globbing CRD files in %s: %w", basesDir, err)
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no CRD YAML files found in %s", basesDir)
	}

	var allFields []types.GoField
	for _, file := range files {
		// Skip common CRDs (fieldexports, iamroleselectors) that aren't service resources
		base := filepath.Base(file)
		if strings.HasPrefix(base, "services.k8s.aws_") {
			continue
		}

		fields, _, err := p.ParseCRDFile(file)
		if err != nil {
			// Log warning but continue
			fmt.Fprintf(os.Stderr, "Warning: failed to parse CRD %s: %v\n", base, err)
			continue
		}
		allFields = append(allFields, fields...)
	}

	return allFields, nil
}

// jsonToCamel converts a JSON camelCase field name to Go CamelCase (exported).
// e.g., "policyDocument" → "PolicyDocument", "redrivePolicy" → "RedrivePolicy"
func jsonToCamel(s string) string {
	if s == "" {
		return ""
	}
	// Capitalize the first letter
	return strings.ToUpper(s[:1]) + s[1:]
}
