package matcher

import (
	"testing"

	"github.com/aws-controllers-k8s/ack-scanner/pkg/types"
)

func TestMatch_CategoryAssignment(t *testing.T) {
	m := &Matcher{}

	ackFields := []types.ScanResult{
		{
			ServiceName:    "iam",
			RepoName:       "iam-controller",
			ResourceName:   "Role",
			FieldName:      "AssumeRolePolicyDocument",
			FieldPath:      "Spec.AssumeRolePolicyDocument",
			GoType:         "*string",
			AnnotationType: types.AnnotationDocument,
		},
		{
			ServiceName:    "sns",
			RepoName:       "sns-controller",
			ResourceName:   "Subscription",
			FieldName:      "FilterPolicy",
			FieldPath:      "Spec.FilterPolicy",
			GoType:         "*string",
			AnnotationType: types.AnnotationNone,
		},
		{
			ServiceName:    "eks",
			RepoName:       "eks-controller",
			ResourceName:   "Addon",
			FieldName:      "ConfigurationValues",
			FieldPath:      "Spec.ConfigurationValues",
			GoType:         "*string",
			AnnotationType: types.AnnotationNone,
		},
	}

	tfFields := []types.TerraformField{
		{
			ServiceName:     "sns",
			ResourceType:    "topic_subscription",
			FieldName:       "filter_policy",
			Description:     "JSON String with the filter policy",
			DetectionMethod: types.DetectDescriptionPhrase,
		},
		{
			ServiceName:     "lambda",
			ResourceType:    "function",
			FieldName:       "policy",
			Description:     "JSON policy document",
			DetectionMethod: types.DetectDescriptionPhrase,
		},
	}

	results := m.Match(ackFields, tfFields)

	// Expected: 3 results (annotated IAM, confirmed SNS, terraform-only lambda)
	// EKS ConfigurationValues is excluded — no TF match and not annotated
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d: %+v", len(results), results)
	}

	// Verify category assignments
	categoryMap := make(map[string]types.Category)
	for _, r := range results {
		categoryMap[r.ServiceName+"/"+r.FieldName] = r.Category
	}

	// IAM AssumeRolePolicyDocument is already annotated
	if cat := categoryMap["iam/AssumeRolePolicyDocument"]; cat != types.CategoryAnnotated {
		t.Errorf("iam/AssumeRolePolicyDocument: expected %s, got %s", types.CategoryAnnotated, cat)
	}

	// SNS FilterPolicy is unannotated and found in Terraform → gap_confirmed
	if cat := categoryMap["sns/FilterPolicy"]; cat != types.CategoryGapConfirmed {
		t.Errorf("sns/FilterPolicy: expected %s, got %s", types.CategoryGapConfirmed, cat)
	}

	// Lambda policy is in Terraform but not in ACK → terraform_only
	if cat := categoryMap["lambda/policy"]; cat != types.CategoryTerraformOnly {
		t.Errorf("lambda/policy: expected %s, got %s", types.CategoryTerraformOnly, cat)
	}

	// EKS ConfigurationValues should NOT appear (no TF match, not annotated)
	if _, exists := categoryMap["eks/ConfigurationValues"]; exists {
		t.Error("eks/ConfigurationValues should not appear — no TF match and not annotated")
	}
}

func TestMatch_TFConfirmation(t *testing.T) {
	m := &Matcher{}

	ackFields := []types.ScanResult{
		{
			ServiceName:    "iam",
			RepoName:       "iam-controller",
			ResourceName:   "Role",
			FieldName:      "AssumeRolePolicyDocument",
			FieldPath:      "Spec.AssumeRolePolicyDocument",
			GoType:         "*string",
			AnnotationType: types.AnnotationDocument,
		},
		{
			ServiceName:    "sns",
			RepoName:       "sns-controller",
			ResourceName:   "Subscription",
			FieldName:      "FilterPolicy",
			FieldPath:      "Spec.FilterPolicy",
			GoType:         "*string",
			AnnotationType: types.AnnotationNone,
		},
		{
			ServiceName:    "eks",
			RepoName:       "eks-controller",
			ResourceName:   "Addon",
			FieldName:      "ConfigurationValues",
			FieldPath:      "Spec.ConfigurationValues",
			GoType:         "*string",
			AnnotationType: types.AnnotationNone,
		},
	}

	tfFields := []types.TerraformField{
		{
			ServiceName:     "sns",
			ResourceType:    "topic_subscription",
			FieldName:       "filter_policy",
			Description:     "JSON String with the filter policy",
			DetectionMethod: types.DetectDescriptionPhrase,
		},
	}

	results := m.Match(ackFields, tfFields)

	confirmationMap := make(map[string]types.TerraformConfirmation)
	for _, r := range results {
		confirmationMap[r.ServiceName+"/"+r.FieldName] = r.TFConfirmation
	}

	// Already annotated → not-applicable
	if tf := confirmationMap["iam/AssumeRolePolicyDocument"]; tf != types.TFNotApplicable {
		t.Errorf("iam/AssumeRolePolicyDocument: expected %s, got %s", types.TFNotApplicable, tf)
	}

	// Unannotated + found in TF → confirmed
	if tf := confirmationMap["sns/FilterPolicy"]; tf != types.TFConfirmed {
		t.Errorf("sns/FilterPolicy: expected %s, got %s", types.TFConfirmed, tf)
	}

	// Unannotated + not found in TF → excluded from results
	if _, exists := confirmationMap["eks/ConfigurationValues"]; exists {
		t.Error("eks/ConfigurationValues should not appear — no TF match and not annotated")
	}
}

func TestMatch_CaseInsensitiveServiceMatching(t *testing.T) {
	m := &Matcher{}

	ackFields := []types.ScanResult{
		{
			ServiceName:    "IAM",
			RepoName:       "iam-controller",
			ResourceName:   "Role",
			FieldName:      "PolicyDocument",
			FieldPath:      "Spec.PolicyDocument",
			GoType:         "*string",
			AnnotationType: types.AnnotationNone,
		},
	}

	tfFields := []types.TerraformField{
		{
			ServiceName:     "iam",
			ResourceType:    "role",
			FieldName:       "policy_document",
			Description:     "JSON policy document",
			DetectionMethod: types.DetectDescriptionPhrase,
		},
	}

	results := m.Match(ackFields, tfFields)

	// Should match despite different casing on service name
	found := false
	for _, r := range results {
		if r.FieldName == "PolicyDocument" && r.Category == types.CategoryGapConfirmed {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected case-insensitive service matching to produce gap_confirmed_by_terraform")
	}

	// The TF field should NOT appear as terraform_only since it was matched
	for _, r := range results {
		if r.Category == types.CategoryTerraformOnly {
			t.Errorf("unexpected terraform_only result: %+v", r)
		}
	}
}

func TestMatch_CamelToSnakeFieldMatching(t *testing.T) {
	m := &Matcher{}

	ackFields := []types.ScanResult{
		{
			ServiceName:    "iam",
			RepoName:       "iam-controller",
			ResourceName:   "Role",
			FieldName:      "AssumeRolePolicyDocument",
			FieldPath:      "Spec.AssumeRolePolicyDocument",
			GoType:         "*string",
			AnnotationType: types.AnnotationNone,
		},
	}

	tfFields := []types.TerraformField{
		{
			ServiceName:     "iam",
			ResourceType:    "role",
			FieldName:       "assume_role_policy_document",
			Description:     "JSON policy document",
			DetectionMethod: types.DetectDescriptionPhrase,
		},
	}

	results := m.Match(ackFields, tfFields)

	// CamelCase ACK field should match snake_case TF field
	found := false
	for _, r := range results {
		if r.FieldName == "AssumeRolePolicyDocument" && r.Category == types.CategoryGapConfirmed {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected CamelCase to snake_case matching to produce gap_confirmed_by_terraform")
	}
}

func TestMatch_EmptyInputs(t *testing.T) {
	m := &Matcher{}

	// Both empty
	results := m.Match(nil, nil)
	if len(results) != 0 {
		t.Errorf("expected 0 results for nil inputs, got %d", len(results))
	}

	// Empty ACK, some TF
	tfFields := []types.TerraformField{
		{
			ServiceName:     "iam",
			ResourceType:    "role",
			FieldName:       "policy",
			Description:     "JSON policy",
			DetectionMethod: types.DetectDescriptionPhrase,
		},
	}
	results = m.Match(nil, tfFields)
	if len(results) != 1 {
		t.Fatalf("expected 1 terraform_only result, got %d", len(results))
	}
	if results[0].Category != types.CategoryTerraformOnly {
		t.Errorf("expected terraform_only, got %s", results[0].Category)
	}

	// Some ACK, empty TF — unannotated fields with no TF match are excluded
	ackFields := []types.ScanResult{
		{
			ServiceName:    "iam",
			RepoName:       "iam-controller",
			ResourceName:   "Role",
			FieldName:      "PolicyDocument",
			FieldPath:      "Spec.PolicyDocument",
			GoType:         "*string",
			AnnotationType: types.AnnotationNone,
		},
	}
	results = m.Match(ackFields, nil)
	if len(results) != 0 {
		t.Fatalf("expected 0 results (no TF to confirm, not annotated), got %d", len(results))
	}

	// Annotated ACK field with empty TF should still appear
	ackFieldsAnnotated := []types.ScanResult{
		{
			ServiceName:    "iam",
			RepoName:       "iam-controller",
			ResourceName:   "Role",
			FieldName:      "PolicyDocument",
			FieldPath:      "Spec.PolicyDocument",
			GoType:         "*string",
			AnnotationType: types.AnnotationDocument,
		},
	}
	results = m.Match(ackFieldsAnnotated, nil)
	if len(results) != 1 {
		t.Fatalf("expected 1 annotated result, got %d", len(results))
	}
	if results[0].Category != types.CategoryAnnotated {
		t.Errorf("expected already_annotated, got %s", results[0].Category)
	}
}

func TestFilterByServices(t *testing.T) {
	m := &Matcher{}

	results := []types.MatchResult{
		{ServiceName: "iam", FieldName: "PolicyDocument", Category: types.CategoryGapConfirmed},
		{ServiceName: "sns", FieldName: "FilterPolicy", Category: types.CategoryGapConfirmed},
		{ServiceName: "eks", FieldName: "ConfigurationValues", Category: types.CategoryGapUnconfirmed},
		{ServiceName: "s3", FieldName: "BucketPolicy", Category: types.CategoryAnnotated},
	}

	// Filter for iam and eks
	filtered := m.FilterByServices(results, []string{"iam", "eks"})
	if len(filtered) != 2 {
		t.Fatalf("expected 2 filtered results, got %d", len(filtered))
	}
	for _, r := range filtered {
		if r.ServiceName != "iam" && r.ServiceName != "eks" {
			t.Errorf("unexpected service in filtered results: %s", r.ServiceName)
		}
	}

	// Empty filter returns all
	filtered = m.FilterByServices(results, nil)
	if len(filtered) != 4 {
		t.Fatalf("expected all 4 results with nil filter, got %d", len(filtered))
	}

	filtered = m.FilterByServices(results, []string{})
	if len(filtered) != 4 {
		t.Fatalf("expected all 4 results with empty filter, got %d", len(filtered))
	}
}

func TestMatch_AnnotatedFieldWithTFMatch(t *testing.T) {
	m := &Matcher{}

	// An annotated ACK field that also exists in Terraform should:
	// 1. Be categorized as already_annotated (not gap_confirmed)
	// 2. NOT produce a terraform_only entry for the TF field
	ackFields := []types.ScanResult{
		{
			ServiceName:    "iam",
			RepoName:       "iam-controller",
			ResourceName:   "Role",
			FieldName:      "AssumeRolePolicyDocument",
			FieldPath:      "Spec.AssumeRolePolicyDocument",
			GoType:         "*string",
			AnnotationType: types.AnnotationDocument,
		},
	}

	tfFields := []types.TerraformField{
		{
			ServiceName:     "iam",
			ResourceType:    "role",
			FieldName:       "assume_role_policy_document",
			Description:     "JSON policy document",
			DetectionMethod: types.DetectDescriptionPhrase,
		},
	}

	results := m.Match(ackFields, tfFields)

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d: %+v", len(results), results)
	}

	if results[0].Category != types.CategoryAnnotated {
		t.Errorf("expected already_annotated, got %s", results[0].Category)
	}
	if results[0].TFConfirmation != types.TFNotApplicable {
		t.Errorf("expected not-applicable, got %s", results[0].TFConfirmation)
	}
}

func TestFilterByServices_CaseInsensitive(t *testing.T) {
	m := &Matcher{}

	results := []types.MatchResult{
		{ServiceName: "IAM", FieldName: "PolicyDocument", Category: types.CategoryGapConfirmed},
		{ServiceName: "sns", FieldName: "FilterPolicy", Category: types.CategoryGapConfirmed},
	}

	// Filter with different casing should still work
	filtered := m.FilterByServices(results, []string{"iam"})
	if len(filtered) != 1 {
		t.Fatalf("expected 1 filtered result for case-insensitive match, got %d", len(filtered))
	}
	if filtered[0].ServiceName != "IAM" {
		t.Errorf("expected IAM, got %s", filtered[0].ServiceName)
	}
}
