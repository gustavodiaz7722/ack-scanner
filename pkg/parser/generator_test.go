package parser

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/aws-controllers-k8s/ack-scanner/pkg/types"
)

func TestParseGeneratorConfig_ValidFile(t *testing.T) {
	content := `resources:
  Role:
    fields:
      AssumeRolePolicyDocument:
        is_document: true
      Tags:
        is_primary_key: true
  Policy:
    fields:
      PolicyDocument:
        is_iam_policy: true
      Description:
        is_required: true
`
	dir := t.TempDir()
	path := filepath.Join(dir, "generator.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	p := &GeneratorParser{}
	config, err := p.ParseGeneratorConfig(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(config.Resources) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(config.Resources))
	}

	// Check Role resource
	role, ok := config.Resources["Role"]
	if !ok {
		t.Fatal("expected 'Role' resource")
	}
	if !role.Fields["AssumeRolePolicyDocument"].IsDocument {
		t.Error("expected AssumeRolePolicyDocument.IsDocument to be true")
	}
	if role.Fields["Tags"].IsDocument {
		t.Error("expected Tags.IsDocument to be false")
	}

	// Check Policy resource
	policy, ok := config.Resources["Policy"]
	if !ok {
		t.Fatal("expected 'Policy' resource")
	}
	if !policy.Fields["PolicyDocument"].IsIAMPolicy {
		t.Error("expected PolicyDocument.IsIAMPolicy to be true")
	}
	if policy.Fields["Description"].IsIAMPolicy {
		t.Error("expected Description.IsIAMPolicy to be false")
	}
}

func TestParseGeneratorConfig_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "generator.yaml")
	if err := os.WriteFile(path, []byte(""), 0644); err != nil {
		t.Fatal(err)
	}

	p := &GeneratorParser{}
	config, err := p.ParseGeneratorConfig(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if config.Resources == nil {
		t.Fatal("expected non-nil Resources map")
	}
	if len(config.Resources) != 0 {
		t.Fatalf("expected 0 resources, got %d", len(config.Resources))
	}
}

func TestParseGeneratorConfig_FileNotFound(t *testing.T) {
	p := &GeneratorParser{}
	_, err := p.ParseGeneratorConfig("/nonexistent/path/generator.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestParseGeneratorConfig_MalformedYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "generator.yaml")
	content := `resources:
  Role:
    fields:
      - this is invalid yaml for a map
      AssumeRolePolicyDocument:
        is_document: true
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	p := &GeneratorParser{}
	_, err := p.ParseGeneratorConfig(path)
	if err == nil {
		t.Fatal("expected error for malformed YAML")
	}
}

func TestParseGeneratorConfig_NoResourcesKey(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "generator.yaml")
	content := `ignore:
  resource_names:
    - Object
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	p := &GeneratorParser{}
	config, err := p.ParseGeneratorConfig(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if config.Resources == nil {
		t.Fatal("expected non-nil Resources map")
	}
	if len(config.Resources) != 0 {
		t.Fatalf("expected 0 resources, got %d", len(config.Resources))
	}
}

func TestGetAnnotatedFields(t *testing.T) {
	config := &types.GeneratorConfig{
		Resources: map[string]types.ResourceConfig{
			"Role": {
				Fields: map[string]types.FieldConfig{
					"AssumeRolePolicyDocument": {IsDocument: true},
					"Tags":                     {},
				},
			},
			"Policy": {
				Fields: map[string]types.FieldConfig{
					"PolicyDocument": {IsIAMPolicy: true},
					"Description":    {},
				},
			},
			"Addon": {
				Fields: map[string]types.FieldConfig{
					"ConfigurationValues": {IsDocument: true},
					"ClusterName":         {},
				},
			},
		},
	}

	p := &GeneratorParser{}
	annotated := p.GetAnnotatedFields(config)

	if len(annotated) != 3 {
		t.Fatalf("expected 3 annotated fields, got %d", len(annotated))
	}

	// Sort for deterministic assertions
	sort.Slice(annotated, func(i, j int) bool {
		if annotated[i].ResourceName != annotated[j].ResourceName {
			return annotated[i].ResourceName < annotated[j].ResourceName
		}
		return annotated[i].FieldName < annotated[j].FieldName
	})

	expected := []types.AnnotatedField{
		{ResourceName: "Addon", FieldName: "ConfigurationValues", AnnotationType: types.AnnotationDocument},
		{ResourceName: "Policy", FieldName: "PolicyDocument", AnnotationType: types.AnnotationIAMPolicy},
		{ResourceName: "Role", FieldName: "AssumeRolePolicyDocument", AnnotationType: types.AnnotationDocument},
	}

	for i, exp := range expected {
		if annotated[i] != exp {
			t.Errorf("annotated[%d] = %+v, want %+v", i, annotated[i], exp)
		}
	}
}

func TestGetAnnotatedFields_NilConfig(t *testing.T) {
	p := &GeneratorParser{}
	annotated := p.GetAnnotatedFields(nil)
	if annotated != nil {
		t.Fatalf("expected nil for nil config, got %+v", annotated)
	}
}

func TestGetAnnotatedFields_NoAnnotations(t *testing.T) {
	config := &types.GeneratorConfig{
		Resources: map[string]types.ResourceConfig{
			"Bucket": {
				Fields: map[string]types.FieldConfig{
					"Name":       {},
					"Encryption": {},
				},
			},
		},
	}

	p := &GeneratorParser{}
	annotated := p.GetAnnotatedFields(config)
	if len(annotated) != 0 {
		t.Fatalf("expected 0 annotated fields, got %d", len(annotated))
	}
}

func TestParseGeneratorConfig_RealWorldStructure(t *testing.T) {
	// Simulates the EKS controller generator.yaml structure
	content := `operations:
  AssociateIdentityProviderConfig:
    operation_type:
    - Create
resources:
  Addon:
    hooks:
      delta_pre_compare:
        code: customPreCompare(delta, a, b)
    fields:
      ClusterName:
        references:
          resource: Cluster
          path: Spec.Name
      ConfigurationValues:
        is_document: true
      ServiceAccountRoleArn:
        references:
          service_name: iam
          resource: Role
          path: Status.ACKResourceMetadata.ARN
  PodIdentityAssociation:
    fields:
      ClusterName:
        references:
          resource: Cluster
          path: Spec.Name
        is_primary_key: true
      Policy:
        is_iam_policy: true
`
	dir := t.TempDir()
	path := filepath.Join(dir, "generator.yaml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	p := &GeneratorParser{}
	config, err := p.ParseGeneratorConfig(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check Addon.ConfigurationValues is_document
	addon, ok := config.Resources["Addon"]
	if !ok {
		t.Fatal("expected 'Addon' resource")
	}
	if !addon.Fields["ConfigurationValues"].IsDocument {
		t.Error("expected ConfigurationValues.IsDocument to be true")
	}
	if addon.Fields["ClusterName"].IsDocument || addon.Fields["ClusterName"].IsIAMPolicy {
		t.Error("expected ClusterName to have no annotations")
	}

	// Check PodIdentityAssociation.Policy is_iam_policy
	pia, ok := config.Resources["PodIdentityAssociation"]
	if !ok {
		t.Fatal("expected 'PodIdentityAssociation' resource")
	}
	if !pia.Fields["Policy"].IsIAMPolicy {
		t.Error("expected Policy.IsIAMPolicy to be true")
	}

	// Verify GetAnnotatedFields
	annotated := p.GetAnnotatedFields(config)
	if len(annotated) != 2 {
		t.Fatalf("expected 2 annotated fields, got %d", len(annotated))
	}
}
