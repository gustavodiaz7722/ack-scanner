package reporter

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/aws-controllers-k8s/ack-scanner/pkg/types"
)

func TestNewReporter(t *testing.T) {
	r := NewReporter("json")
	if r.format != "json" {
		t.Errorf("expected format 'json', got %q", r.format)
	}
}

func TestGenerateSummary(t *testing.T) {
	results := []types.MatchResult{
		{ServiceName: "s3", Category: types.CategoryGapConfirmed},
		{ServiceName: "s3", Category: types.CategoryGapConfirmed},
		{ServiceName: "iam", Category: types.CategoryGapConfirmed},
		{ServiceName: "s3", Category: types.CategoryGapUnconfirmed},
		{ServiceName: "eks", Category: types.CategoryAnnotated},
		{ServiceName: "sns", Category: types.CategoryTerraformOnly},
	}

	summary := GenerateSummary(results)

	if summary.TotalFields != 6 {
		t.Errorf("expected TotalFields=6, got %d", summary.TotalFields)
	}
	if summary.AnnotatedCount != 1 {
		t.Errorf("expected AnnotatedCount=1, got %d", summary.AnnotatedCount)
	}
	if summary.GapConfirmedCount != 3 {
		t.Errorf("expected GapConfirmedCount=3, got %d", summary.GapConfirmedCount)
	}
	if summary.GapUnconfirmedCount != 1 {
		t.Errorf("expected GapUnconfirmedCount=1, got %d", summary.GapUnconfirmedCount)
	}
	if summary.TerraformOnlyCount != 1 {
		t.Errorf("expected TerraformOnlyCount=1, got %d", summary.TerraformOnlyCount)
	}

	// GapsPerService: s3=3 (2 confirmed + 1 unconfirmed), iam=1 (1 confirmed)
	if summary.GapsPerService["s3"] != 3 {
		t.Errorf("expected GapsPerService[s3]=3, got %d", summary.GapsPerService["s3"])
	}
	if summary.GapsPerService["iam"] != 1 {
		t.Errorf("expected GapsPerService[iam]=1, got %d", summary.GapsPerService["iam"])
	}

	// ServicesByPriority: s3 (2 confirmed) > iam (1 confirmed)
	if len(summary.ServicesByPriority) != 2 {
		t.Fatalf("expected 2 services in priority list, got %d", len(summary.ServicesByPriority))
	}
	if summary.ServicesByPriority[0].ServiceName != "s3" {
		t.Errorf("expected first priority service='s3', got %q", summary.ServicesByPriority[0].ServiceName)
	}
	if summary.ServicesByPriority[0].ConfirmedGapCount != 2 {
		t.Errorf("expected first priority count=2, got %d", summary.ServicesByPriority[0].ConfirmedGapCount)
	}
	if summary.ServicesByPriority[1].ServiceName != "iam" {
		t.Errorf("expected second priority service='iam', got %q", summary.ServicesByPriority[1].ServiceName)
	}
}

func TestSortResults(t *testing.T) {
	results := []types.MatchResult{
		{ServiceName: "eks", FieldName: "addon_config", Category: types.CategoryAnnotated},
		{ServiceName: "iam", FieldName: "policy_doc", Category: types.CategoryGapUnconfirmed},
		{ServiceName: "s3", FieldName: "policy", Category: types.CategoryGapConfirmed},
		{ServiceName: "sns", FieldName: "filter", Category: types.CategoryTerraformOnly},
		{ServiceName: "iam", FieldName: "role_policy", Category: types.CategoryGapConfirmed},
	}

	sorted := sortResults(results)

	// Expected order: gap_confirmed (iam, s3), gap_unconfirmed (iam), annotated (eks), terraform_only (sns)
	expectedOrder := []struct {
		service  string
		field    string
		category types.Category
	}{
		{"iam", "role_policy", types.CategoryGapConfirmed},
		{"s3", "policy", types.CategoryGapConfirmed},
		{"iam", "policy_doc", types.CategoryGapUnconfirmed},
		{"eks", "addon_config", types.CategoryAnnotated},
		{"sns", "filter", types.CategoryTerraformOnly},
	}

	for i, exp := range expectedOrder {
		if sorted[i].ServiceName != exp.service || sorted[i].FieldName != exp.field || sorted[i].Category != exp.category {
			t.Errorf("position %d: expected (%s, %s, %s), got (%s, %s, %s)",
				i, exp.service, exp.field, exp.category,
				sorted[i].ServiceName, sorted[i].FieldName, sorted[i].Category)
		}
	}
}

func TestGenerateReportJSON(t *testing.T) {
	r := NewReporter("json")
	results := []types.MatchResult{
		{
			ServiceName:      "iam",
			ResourceName:     "Role",
			FieldName:        "AssumeRolePolicyDocument",
			FieldPath:        "Spec.AssumeRolePolicyDocument",
			AnnotationStatus: types.AnnotationNone,
			TFConfirmation:   types.TFConfirmed,
			Category:         types.CategoryGapConfirmed,
		},
		{
			ServiceName:      "eks",
			ResourceName:     "Addon",
			FieldName:        "ConfigurationValues",
			FieldPath:        "Spec.ConfigurationValues",
			AnnotationStatus: types.AnnotationDocument,
			TFConfirmation:   types.TFNotApplicable,
			Category:         types.CategoryAnnotated,
		},
	}

	var buf bytes.Buffer
	err := r.GenerateReport(results, &buf)
	if err != nil {
		t.Fatalf("GenerateReport failed: %v", err)
	}

	// Verify valid JSON
	var output reportJSON
	if err := json.Unmarshal(buf.Bytes(), &output); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	// Verify fields are present
	if len(output.Fields) != 2 {
		t.Errorf("expected 2 fields, got %d", len(output.Fields))
	}

	// Verify sort order: gap_confirmed comes before annotated
	if output.Fields[0].Category != types.CategoryGapConfirmed {
		t.Errorf("expected first field category=gap_confirmed, got %s", output.Fields[0].Category)
	}
	if output.Fields[1].Category != types.CategoryAnnotated {
		t.Errorf("expected second field category=annotated, got %s", output.Fields[1].Category)
	}

	// Verify summary
	if output.Summary.TotalFields != 2 {
		t.Errorf("expected summary total=2, got %d", output.Summary.TotalFields)
	}
	if output.Summary.GapConfirmedCount != 1 {
		t.Errorf("expected summary gap_confirmed=1, got %d", output.Summary.GapConfirmedCount)
	}
	if output.Summary.AnnotatedCount != 1 {
		t.Errorf("expected summary annotated=1, got %d", output.Summary.AnnotatedCount)
	}
}

func TestFormatControllerListJSON(t *testing.T) {
	r := NewReporter("json")
	repos := []types.ControllerRepo{
		{RepoName: "s3-controller", ServiceName: "s3"},
		{RepoName: "iam-controller", ServiceName: "iam"},
		{RepoName: "eks-controller", ServiceName: "eks"},
	}

	var buf bytes.Buffer
	err := r.FormatControllerList(repos, &buf)
	if err != nil {
		t.Fatalf("FormatControllerList failed: %v", err)
	}

	// Verify valid JSON
	var output []controllerListEntry
	if err := json.Unmarshal(buf.Bytes(), &output); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	// Verify sorted by service name (eks, iam, s3)
	if len(output) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(output))
	}
	if output[0].ServiceName != "eks" {
		t.Errorf("expected first service='eks', got %q", output[0].ServiceName)
	}
	if output[1].ServiceName != "iam" {
		t.Errorf("expected second service='iam', got %q", output[1].ServiceName)
	}
	if output[2].ServiceName != "s3" {
		t.Errorf("expected third service='s3', got %q", output[2].ServiceName)
	}

	// Verify repository_name is present
	if output[0].RepositoryName != "eks-controller" {
		t.Errorf("expected first repo='eks-controller', got %q", output[0].RepositoryName)
	}
}

func TestFormatControllerListJSON_Empty(t *testing.T) {
	r := NewReporter("json")
	var buf bytes.Buffer
	err := r.FormatControllerList([]types.ControllerRepo{}, &buf)
	if err != nil {
		t.Fatalf("FormatControllerList failed: %v", err)
	}

	// Verify empty array, not null
	output := buf.String()
	if output != "[]\n" {
		t.Errorf("expected empty JSON array, got %q", output)
	}
}

func TestGenerateReportJSON_Empty(t *testing.T) {
	r := NewReporter("json")
	var buf bytes.Buffer
	err := r.GenerateReport([]types.MatchResult{}, &buf)
	if err != nil {
		t.Fatalf("GenerateReport failed: %v", err)
	}

	var output reportJSON
	if err := json.Unmarshal(buf.Bytes(), &output); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	if len(output.Fields) != 0 {
		t.Errorf("expected 0 fields, got %d", len(output.Fields))
	}
	if output.Summary.TotalFields != 0 {
		t.Errorf("expected summary total=0, got %d", output.Summary.TotalFields)
	}
}

func TestGenerateReport_UnsupportedFormat(t *testing.T) {
	r := NewReporter("xml")
	var buf bytes.Buffer
	err := r.GenerateReport([]types.MatchResult{}, &buf)
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

func TestFormatControllerList_UnsupportedFormat(t *testing.T) {
	r := NewReporter("xml")
	var buf bytes.Buffer
	err := r.FormatControllerList([]types.ControllerRepo{}, &buf)
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}
