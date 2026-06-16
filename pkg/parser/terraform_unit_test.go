package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aws-controllers-k8s/ack-scanner/pkg/types"
)

func TestServiceFromFileName(t *testing.T) {
	tests := []struct {
		fileName string
		expected string
	}{
		{"sns_topic_subscription.html.markdown", "sns"},
		{"iam_role.html.markdown", "iam"},
		{"eks_addon.html.markdown", "eks"},
		{"s3_bucket_policy.html.markdown", "s3"},
		{"lambda_function.html.markdown", "lambda"},
		{"singleword.html.markdown", "singleword"},
	}

	for _, tc := range tests {
		t.Run(tc.fileName, func(t *testing.T) {
			got := ServiceFromFileName(tc.fileName)
			if got != tc.expected {
				t.Errorf("ServiceFromFileName(%q) = %q, want %q", tc.fileName, got, tc.expected)
			}
		})
	}
}

func TestResourceFromFileName(t *testing.T) {
	tests := []struct {
		fileName string
		expected string
	}{
		{"sns_topic_subscription.html.markdown", "topic_subscription"},
		{"iam_role.html.markdown", "role"},
		{"eks_addon.html.markdown", "addon"},
		{"s3_bucket_policy.html.markdown", "bucket_policy"},
		{"singleword.html.markdown", ""},
	}

	for _, tc := range tests {
		t.Run(tc.fileName, func(t *testing.T) {
			got := ResourceFromFileName(tc.fileName)
			if got != tc.expected {
				t.Errorf("ResourceFromFileName(%q) = %q, want %q", tc.fileName, got, tc.expected)
			}
		})
	}
}

func TestExtractArguments(t *testing.T) {
	content := `---
subcategory: "SNS"
---

# Resource: aws_sns_topic_subscription

## Example Usage

Some example here.

## Argument Reference

The following arguments are required:

* ` + "`topic_arn`" + ` - (Required) ARN of the SNS topic to subscribe to.
* ` + "`protocol`" + ` - (Required) Protocol to use.
* ` + "`filter_policy`" + ` - (Optional) JSON String with the filter policy that will be used in the subscription.
* ` + "`endpoint`" + ` - (Required) Endpoint to send data to.

## Attribute Reference

This resource exports the following attributes:
`

	p := &TerraformParser{}
	args := p.ExtractArguments(content)

	if len(args) != 4 {
		t.Fatalf("expected 4 arguments, got %d", len(args))
	}

	// Verify first entry
	if args[0].FieldName != "topic_arn" {
		t.Errorf("args[0].FieldName = %q, want %q", args[0].FieldName, "topic_arn")
	}
	if !args[0].Required {
		t.Error("args[0].Required should be true")
	}

	// Verify optional field
	if args[2].FieldName != "filter_policy" {
		t.Errorf("args[2].FieldName = %q, want %q", args[2].FieldName, "filter_policy")
	}
	if args[2].Required {
		t.Error("args[2].Required should be false")
	}
	if args[2].Description != "JSON String with the filter policy that will be used in the subscription." {
		t.Errorf("args[2].Description = %q", args[2].Description)
	}
}

func TestExtractJsonEncodeFields(t *testing.T) {
	content := `## Example Usage

` + "```hcl" + `
resource "aws_iam_role" "example" {
  name = "example"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
    }]
  })

  inline_policy = data.aws_iam_policy_document.example.json
}
` + "```" + `

## Argument Reference
`

	p := &TerraformParser{}
	// When targeting aws_iam_role, both fields should be detected
	fields := p.ExtractJsonEncodeFields(content, "aws_iam_role")

	fieldSet := make(map[string]struct{})
	for _, f := range fields {
		fieldSet[f] = struct{}{}
	}

	if _, ok := fieldSet["assume_role_policy"]; !ok {
		t.Error("expected assume_role_policy to be detected from jsonencode")
	}
	if _, ok := fieldSet["inline_policy"]; !ok {
		t.Error("expected inline_policy to be detected from data.aws_iam_policy_document")
	}
	if len(fields) != 2 {
		t.Errorf("expected 2 fields, got %d: %v", len(fields), fields)
	}
}

func TestExtractJsonEncodeFields_IgnoresOtherResources(t *testing.T) {
	content := `## Example Usage

` + "```hcl" + `
resource "aws_iam_role" "example" {
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
  })
}

resource "aws_eks_addon" "example" {
  addon_name   = "vpc-cni"
  cluster_name = aws_eks_cluster.example.name
}
` + "```" + `

## Argument Reference
`

	p := &TerraformParser{}
	// When targeting aws_eks_addon, assume_role_policy should NOT be detected
	// because it's in the aws_iam_role block, not aws_eks_addon
	fields := p.ExtractJsonEncodeFields(content, "aws_eks_addon")

	for _, f := range fields {
		if f == "assume_role_policy" {
			t.Error("assume_role_policy should not be detected when targeting aws_eks_addon")
		}
	}
}

func TestExtractJsonEncodeFields_NoExampleSection(t *testing.T) {
	content := `## Argument Reference

* ` + "`policy`" + ` - (Required) The policy document.

## Attribute Reference
`

	p := &TerraformParser{}
	fields := p.ExtractJsonEncodeFields(content, "aws_test_resource")

	if len(fields) != 0 {
		t.Errorf("expected 0 fields for content without Example Usage, got %d", len(fields))
	}
}

func TestContainsJSONPhrase(t *testing.T) {
	tests := []struct {
		desc     string
		expected bool
	}{
		{"JSON String with the filter policy", true},
		{"A json-encoded representation of the policy", true},
		{"JSON formatted string of the configuration", true},
		{"The policy document as a JSON string", true},
		{"A JSON schema for validation", true},
		{"The IAM policy document attached to the role", true},
		{"ARN of the SNS topic", false},
		{"Name of the resource", false},
		{"The configuration values as YAML", false},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got := containsJSONPhrase(tc.desc)
			if got != tc.expected {
				t.Errorf("containsJSONPhrase(%q) = %v, want %v", tc.desc, got, tc.expected)
			}
		})
	}
}

func TestParseResourceDoc(t *testing.T) {
	// Create a temp file with terraform doc content
	content := `---
subcategory: "SNS"
---

# Resource: aws_sns_topic_subscription

## Example Usage

### With filter policy

` + "```hcl" + `
resource "aws_sns_topic_subscription" "example" {
  topic_arn     = aws_sns_topic.example.arn
  protocol      = "lambda"
  endpoint      = aws_lambda_function.example.arn
  filter_policy = jsonencode({
    store = ["example_corp"]
  })
}
` + "```" + `

## Argument Reference

The following arguments are required:

* ` + "`topic_arn`" + ` - (Required) ARN of the SNS topic to subscribe to.
* ` + "`protocol`" + ` - (Required) Protocol to use.
* ` + "`filter_policy`" + ` - (Optional) JSON String with the filter policy that will be used in the subscription to filter messages.
* ` + "`redrive_policy`" + ` - (Optional) JSON String with the redrive policy document.

## Attribute Reference

This resource exports the following attributes:
`

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "sns_topic_subscription.html.markdown")
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	p := &TerraformParser{}
	fields, err := p.ParseResourceDoc(filePath)
	if err != nil {
		t.Fatal(err)
	}

	// filter_policy should be "both" (description + jsonencode)
	// redrive_policy should be "description_phrase" only
	fieldMap := make(map[string]types.TerraformField)
	for _, f := range fields {
		fieldMap[f.FieldName] = f
	}

	if len(fields) != 2 {
		t.Fatalf("expected 2 JSON fields, got %d: %+v", len(fields), fields)
	}

	fp, ok := fieldMap["filter_policy"]
	if !ok {
		t.Fatal("expected filter_policy in results")
	}
	if fp.DetectionMethod != types.DetectBoth {
		t.Errorf("filter_policy DetectionMethod = %q, want %q", fp.DetectionMethod, types.DetectBoth)
	}
	if fp.ServiceName != "sns" {
		t.Errorf("filter_policy ServiceName = %q, want %q", fp.ServiceName, "sns")
	}
	if fp.ResourceType != "topic_subscription" {
		t.Errorf("filter_policy ResourceType = %q, want %q", fp.ResourceType, "topic_subscription")
	}

	rp, ok := fieldMap["redrive_policy"]
	if !ok {
		t.Fatal("expected redrive_policy in results")
	}
	if rp.DetectionMethod != types.DetectDescriptionPhrase {
		t.Errorf("redrive_policy DetectionMethod = %q, want %q", rp.DetectionMethod, types.DetectDescriptionPhrase)
	}
}

func TestParseAllDocs(t *testing.T) {
	tmpDir := t.TempDir()

	// Create two doc files
	snsDoc := `## Example Usage

` + "```hcl" + `
resource "aws_sns_topic_subscription" "example" {
  filter_policy = jsonencode({})
}
` + "```" + `

## Argument Reference

* ` + "`filter_policy`" + ` - (Optional) JSON String with filter policy.
* ` + "`topic_arn`" + ` - (Required) ARN of the topic.

## Attribute Reference
`

	iamDoc := `## Example Usage

` + "```hcl" + `
resource "aws_iam_role" "example" {
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
  })
}
` + "```" + `

## Argument Reference

* ` + "`assume_role_policy`" + ` - (Required) Policy document for assuming the role.
* ` + "`name`" + ` - (Required) Name of the role.

## Attribute Reference
`

	os.WriteFile(filepath.Join(tmpDir, "sns_topic_subscription.html.markdown"), []byte(snsDoc), 0644)
	os.WriteFile(filepath.Join(tmpDir, "iam_role.html.markdown"), []byte(iamDoc), 0644)

	p := &TerraformParser{}
	fields, err := p.ParseAllDocs(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	// sns: filter_policy (both)
	// iam: assume_role_policy (both - "policy document" in description + jsonencode in example)
	if len(fields) != 2 {
		t.Fatalf("expected 2 fields, got %d: %+v", len(fields), fields)
	}

	fieldMap := make(map[string]types.TerraformField)
	for _, f := range fields {
		fieldMap[f.ServiceName+"_"+f.FieldName] = f
	}

	if f, ok := fieldMap["sns_filter_policy"]; !ok {
		t.Error("expected sns filter_policy in results")
	} else if f.DetectionMethod != types.DetectBoth {
		t.Errorf("sns filter_policy detection = %q, want %q", f.DetectionMethod, types.DetectBoth)
	}

	if f, ok := fieldMap["iam_assume_role_policy"]; !ok {
		t.Error("expected iam assume_role_policy in results")
	} else if f.DetectionMethod != types.DetectBoth {
		t.Errorf("iam assume_role_policy detection = %q, want %q", f.DetectionMethod, types.DetectBoth)
	}
}

func TestParseAllDocs_MalformedFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a directory that looks like a markdown file (will fail to read)
	os.MkdirAll(filepath.Join(tmpDir, "bad_resource.html.markdown"), 0755)

	// Create a valid doc
	validDoc := `## Argument Reference

* ` + "`policy`" + ` - (Required) A JSON formatted policy.

## Attribute Reference
`
	os.WriteFile(filepath.Join(tmpDir, "good_resource.html.markdown"), []byte(validDoc), 0644)

	p := &TerraformParser{}
	fields, err := p.ParseAllDocs(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	// Should still get the valid file's result
	if len(fields) != 1 {
		t.Fatalf("expected 1 field from valid file, got %d", len(fields))
	}
	if fields[0].FieldName != "policy" {
		t.Errorf("expected field 'policy', got %q", fields[0].FieldName)
	}
}
