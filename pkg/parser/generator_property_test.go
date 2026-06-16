// Feature: ack-scanner, Property 7: Generator.yaml annotation extraction
package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aws-controllers-k8s/ack-scanner/pkg/types"
	"gopkg.in/yaml.v3"
	"pgregory.net/rapid"
)

// genGeneratorResourceName generates simple uppercase identifiers like "Role", "Bucket", "Addon".
func genGeneratorResourceName() *rapid.Generator[string] {
	names := []string{
		"Role", "Bucket", "Addon", "Policy", "Cluster", "Function",
		"Table", "Queue", "Topic", "Stream", "Instance", "Volume",
		"Subnet", "Gateway", "Certificate", "Domain", "Repository",
		"Pipeline", "Stack", "Secret", "Key", "Alarm", "Rule",
	}
	return rapid.SampledFrom(names)
}

// genGeneratorFieldName generates CamelCase field names.
func genGeneratorFieldName() *rapid.Generator[string] {
	names := []string{
		"PolicyDocument", "AssumeRolePolicyDocument", "ConfigurationValues",
		"FilterPolicy", "TemplateBody", "SchemaDefinition",
		"InlinePolicy", "AccessPolicy", "EncryptionConfiguration",
		"Tags", "Name", "Description", "ClusterName", "Status",
		"ARN", "ID", "CreatedAt", "Region", "VpcId", "SubnetIds",
		"ServiceAccountRoleArn", "RoleArn", "KmsKeyId", "BucketName",
	}
	return rapid.SampledFrom(names)
}

// annotationChoice represents the annotation state of a field.
type annotationChoice int

const (
	noAnnotation annotationChoice = iota
	documentAnnotation
	iamPolicyAnnotation
)

// testField holds generated field data for verification.
type testField struct {
	Resource   string
	Field      string
	Annotation annotationChoice
}

// TestProperty7_GeneratorYAMLAnnotationExtraction verifies that for any valid
// generator.yaml with known annotations, the parser extracts exactly those
// annotated fields with correct annotation type and resource association.
//
// **Validates: Requirements 3.2**
func TestProperty7_GeneratorYAMLAnnotationExtraction(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		// Generate between 1 and 5 resources
		numResources := rapid.IntRange(1, 5).Draw(rt, "numResources")

		// Track which resources we pick (use unique names)
		resourceNames := make([]string, 0, numResources)
		usedResources := make(map[string]bool)
		for len(resourceNames) < numResources {
			name := genGeneratorResourceName().Draw(rt, "resourceName")
			if !usedResources[name] {
				usedResources[name] = true
				resourceNames = append(resourceNames, name)
			}
		}

		// For each resource, generate 1-4 fields with random annotations
		var allFields []testField
		configResources := make(map[string]types.ResourceConfig)

		for _, resName := range resourceNames {
			numFields := rapid.IntRange(1, 4).Draw(rt, "numFields_"+resName)
			usedFields := make(map[string]bool)
			fieldConfigs := make(map[string]types.FieldConfig)

			for len(fieldConfigs) < numFields {
				fieldName := genGeneratorFieldName().Draw(rt, "fieldName_"+resName)
				if usedFields[fieldName] {
					continue
				}
				usedFields[fieldName] = true

				// Randomly assign annotation: 0=none, 1=document, 2=iam_policy
				choice := annotationChoice(rapid.IntRange(0, 2).Draw(rt, "annotation_"+resName+"_"+fieldName))

				fc := types.FieldConfig{}
				switch choice {
				case documentAnnotation:
					fc.IsDocument = true
				case iamPolicyAnnotation:
					fc.IsIAMPolicy = true
				}

				fieldConfigs[fieldName] = fc
				allFields = append(allFields, testField{
					Resource:   resName,
					Field:      fieldName,
					Annotation: choice,
				})
			}

			configResources[resName] = types.ResourceConfig{Fields: fieldConfigs}
		}

		// Build the GeneratorConfig and serialize to YAML
		config := types.GeneratorConfig{Resources: configResources}
		yamlData, err := yaml.Marshal(&config)
		if err != nil {
			t.Fatalf("failed to marshal YAML: %v", err)
		}

		// Write to a temp file
		dir, err := os.MkdirTemp("", "generator-property-test-*")
		if err != nil {
			t.Fatalf("failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(dir)

		path := filepath.Join(dir, "generator.yaml")
		if err := os.WriteFile(path, yamlData, 0644); err != nil {
			t.Fatalf("failed to write temp file: %v", err)
		}

		// Parse using ParseGeneratorConfig
		p := &GeneratorParser{}
		parsed, err := p.ParseGeneratorConfig(path)
		if err != nil {
			t.Fatalf("ParseGeneratorConfig failed: %v", err)
		}

		// Get annotated fields
		annotated := p.GetAnnotatedFields(parsed)

		// Build a lookup of annotated fields for verification
		type fieldKey struct {
			Resource string
			Field    string
		}
		annotatedMap := make(map[fieldKey]types.AnnotationType)
		for _, af := range annotated {
			annotatedMap[fieldKey{af.ResourceName, af.FieldName}] = af.AnnotationType
		}

		// Count expected annotated fields
		expectedCount := 0
		for _, tf := range allFields {
			if tf.Annotation != noAnnotation {
				expectedCount++
			}
		}

		// Verify 1: every field with is_document=true or is_iam_policy=true is
		// present in the result with the correct annotation type
		for _, tf := range allFields {
			key := fieldKey{tf.Resource, tf.Field}
			switch tf.Annotation {
			case documentAnnotation:
				at, found := annotatedMap[key]
				if !found {
					t.Errorf("field %s.%s has is_document=true but was not in annotated results", tf.Resource, tf.Field)
				} else if at != types.AnnotationDocument {
					t.Errorf("field %s.%s expected annotation type 'document', got '%s'", tf.Resource, tf.Field, at)
				}
			case iamPolicyAnnotation:
				at, found := annotatedMap[key]
				if !found {
					t.Errorf("field %s.%s has is_iam_policy=true but was not in annotated results", tf.Resource, tf.Field)
				} else if at != types.AnnotationIAMPolicy {
					t.Errorf("field %s.%s expected annotation type 'iam_policy', got '%s'", tf.Resource, tf.Field, at)
				}
			case noAnnotation:
				// Verify 2: no field without an annotation appears in the result
				_, found := annotatedMap[key]
				if found {
					t.Errorf("field %s.%s has no annotation but appeared in annotated results", tf.Resource, tf.Field)
				}
			}
		}

		// Verify 3: the total count matches (no extra entries)
		if len(annotated) != expectedCount {
			t.Errorf("expected %d annotated fields, got %d", expectedCount, len(annotated))
		}

		// Verify 4: resource associations are correct - every returned field
		// must have its resource name matching what we generated
		for _, af := range annotated {
			key := fieldKey{af.ResourceName, af.FieldName}
			found := false
			for _, tf := range allFields {
				if tf.Resource == key.Resource && tf.Field == key.Field {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("annotated field %s.%s not found in generated fields", af.ResourceName, af.FieldName)
			}
		}
	})
}
