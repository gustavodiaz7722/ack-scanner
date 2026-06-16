package parser

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/aws-controllers-k8s/ack-scanner/pkg/types"
)

// GeneratorParser parses generator.yaml to extract annotation information.
type GeneratorParser struct{}

// ParseGeneratorConfig parses a generator.yaml and returns annotated fields.
func (p *GeneratorParser) ParseGeneratorConfig(filePath string) (*types.GeneratorConfig, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading generator config %s: %w", filePath, err)
	}

	// Handle empty file: return empty config
	if len(data) == 0 {
		return &types.GeneratorConfig{
			Resources: make(map[string]types.ResourceConfig),
		}, nil
	}

	var config types.GeneratorConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parsing generator config %s: %w", filePath, err)
	}

	// Ensure Resources map is initialized even if YAML had no resources key
	if config.Resources == nil {
		config.Resources = make(map[string]types.ResourceConfig)
	}

	return &config, nil
}

// GetAnnotatedFields extracts all fields with either is_document or is_iam_policy
// set to true, returning them as AnnotatedField structs.
func (p *GeneratorParser) GetAnnotatedFields(config *types.GeneratorConfig) []types.AnnotatedField {
	if config == nil {
		return nil
	}

	var annotated []types.AnnotatedField
	for resourceName, resource := range config.Resources {
		for fieldName, field := range resource.Fields {
			if field.IsDocument {
				annotated = append(annotated, types.AnnotatedField{
					ResourceName:   resourceName,
					FieldName:      fieldName,
					AnnotationType: types.AnnotationDocument,
				})
			} else if field.IsIAMPolicy {
				annotated = append(annotated, types.AnnotatedField{
					ResourceName:   resourceName,
					FieldName:      fieldName,
					AnnotationType: types.AnnotationIAMPolicy,
				})
			}
		}
	}
	return annotated
}
