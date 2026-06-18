package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
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

// ParseCRDFile parses a single CRD YAML file and returns ALL string fields
// found under spec, recursively traversing nested objects. No heuristic
// filtering is applied — the caller (matcher) cross-references against
// Terraform's JSON field list to determine which fields are document-string
// candidates.
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

	var fields []types.GoField
	collectStringFields(specProps, resourceName, resourceName, &fields)

	// Deduplicate by (StructName, FieldName). The same leaf field name can
	// appear at many nesting levels in complex CRDs (e.g., ContentType in
	// sagemaker ModelPackage appears 18 times). For matching against
	// Terraform's flat field names, we only need one entry per unique
	// (resource, field) pair. Keep the shortest path as the representative.
	fields = deduplicateFields(fields)

	return fields, resourceName, nil
}

// deduplicateFields removes duplicate entries that share the same
// (StructName, FieldName) pair, keeping the entry with the shortest FieldPath
// as the canonical representative.
func deduplicateFields(fields []types.GoField) []types.GoField {
	type key struct {
		structName string
		fieldName  string
	}
	seen := make(map[key]int) // maps to index in result slice
	var result []types.GoField

	for _, f := range fields {
		k := key{structName: f.StructName, fieldName: f.FieldName}
		if idx, exists := seen[k]; exists {
			// Keep the shorter path (closer to spec root)
			if len(f.FieldPath) < len(result[idx].FieldPath) {
				result[idx] = f
			}
		} else {
			seen[k] = len(result)
			result = append(result, f)
		}
	}
	return result
}

// collectStringFields recursively traverses CRD properties and collects all
// string fields at any nesting depth. For each string field found, it records
// the leaf field name (CamelCase) and the full dot-separated path from the
// resource root. This ensures deeply nested fields like
// spec.amazonManagedKafkaEventSourceConfig.schemaRegistryConfig.eventRecordFormat
// are captured and available for matching against Terraform fields.
func collectStringFields(props map[string]crdProperty, resourceName string, pathPrefix string, fields *[]types.GoField) {
	for fieldName, prop := range props {
		camelName := jsonToCamel(fieldName)
		fieldPath := pathPrefix + "." + camelName

		switch prop.Type {
		case "string":
			*fields = append(*fields, types.GoField{
				StructName: resourceName,
				FieldName:  camelName,
				FieldPath:  fieldPath,
				GoType:     "*string",
				JSONTag:    fieldName,
			})
		case "object":
			// Recurse into nested object properties
			if prop.Properties != nil {
				collectStringFields(prop.Properties, resourceName, fieldPath, fields)
			}
		}
	}
}

// ParseAllCRDs parses all CRD YAML files in config/crd/bases/ for a given
// repo directory and returns the combined set of heuristic-matching string fields.
func (p *CRDParser) ParseAllCRDs(repoDir string) ([]types.GoField, error) {
	fields, _, err := p.ParseAllCRDsWithResources(repoDir)
	return fields, err
}

// ParseAllCRDsWithResources parses all CRD YAML files in config/crd/bases/ for
// a given repo directory and returns both the string fields and the list of
// resource kinds found across all CRDs.
func (p *CRDParser) ParseAllCRDsWithResources(repoDir string) ([]types.GoField, []string, error) {
	basesDir := filepath.Join(repoDir, "config", "crd", "bases")
	pattern := filepath.Join(basesDir, "*.yaml")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, nil, fmt.Errorf("globbing CRD files in %s: %w", basesDir, err)
	}

	if len(files) == 0 {
		return nil, nil, fmt.Errorf("no CRD YAML files found in %s", basesDir)
	}

	var allFields []types.GoField
	var resourceKinds []string
	for _, file := range files {
		// Skip common CRDs (fieldexports, iamroleselectors) that aren't service resources
		base := filepath.Base(file)
		if strings.HasPrefix(base, "services.k8s.aws_") {
			continue
		}

		fields, kind, err := p.ParseCRDFile(file)
		if err != nil {
			// Log warning but continue
			fmt.Fprintf(os.Stderr, "Warning: failed to parse CRD %s: %v\n", base, err)
			continue
		}
		if kind != "" {
			resourceKinds = append(resourceKinds, kind)
		}
		allFields = append(allFields, fields...)
	}

	sort.Strings(resourceKinds)
	return allFields, resourceKinds, nil
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
