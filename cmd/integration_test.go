package cmd

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/aws-controllers-k8s/ack-scanner/pkg/matcher"
	"github.com/aws-controllers-k8s/ack-scanner/pkg/parser"
	"github.com/aws-controllers-k8s/ack-scanner/pkg/reporter"
	"github.com/aws-controllers-k8s/ack-scanner/pkg/types"
)

// testdataDir returns the absolute path to the cmd/testdata directory.
func testdataDir(t *testing.T) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("could not determine test file location")
	}
	return filepath.Join(filepath.Dir(filename), "testdata")
}

// TestIntegration_ScanPipeline tests the end-to-end ACK scan pipeline using
// fixture data: findTypesFile → ParseTypesFile → ParseGeneratorConfig → classifyFields.
func TestIntegration_ScanPipeline(t *testing.T) {
	td := testdataDir(t)
	repoDir := filepath.Join(td, "iam-controller")

	// Step 1: Find the types file
	typesFile, err := findTypesFile(repoDir)
	if err != nil {
		t.Fatalf("findTypesFile failed: %v", err)
	}
	expectedTypesFile := filepath.Join(repoDir, "apis", "v1alpha1", "types.go")
	if typesFile != expectedTypesFile {
		t.Errorf("findTypesFile returned %q, want %q", typesFile, expectedTypesFile)
	}

	// Step 2: Parse the types file for *string fields matching heuristics
	goParser := &parser.GoASTParser{}
	fields, err := goParser.ParseTypesFile(typesFile)
	if err != nil {
		t.Fatalf("ParseTypesFile failed: %v", err)
	}

	// Expect 5 heuristic-matching fields (all contain "Policy" or "Document"):
	// - RoleSpec.AssumeRolePolicyDocument
	// - RoleSpec.PermissionsBoundaryPolicyARN
	// - RoleSpec.InlinePolicyDocument
	// - PolicySpec.PolicyDocument
	// - PolicySpec.PolicyName
	if len(fields) != 5 {
		t.Errorf("ParseTypesFile found %d fields, want 5. Fields: %+v", len(fields), fields)
	}

	// Verify specific fields are found
	fieldNames := make(map[string]bool)
	for _, f := range fields {
		fieldNames[f.StructName+"."+f.FieldName] = true
	}
	expectedFields := []string{
		"RoleSpec.AssumeRolePolicyDocument",
		"RoleSpec.PermissionsBoundaryPolicyARN",
		"RoleSpec.InlinePolicyDocument",
		"PolicySpec.PolicyDocument",
		"PolicySpec.PolicyName",
	}
	for _, ef := range expectedFields {
		if !fieldNames[ef] {
			t.Errorf("expected field %q not found in parsed results", ef)
		}
	}

	// Step 3: Parse generator.yaml
	genParser := &parser.GeneratorParser{}
	generatorPath := filepath.Join(repoDir, "generator.yaml")
	genConfig, err := genParser.ParseGeneratorConfig(generatorPath)
	if err != nil {
		t.Fatalf("ParseGeneratorConfig failed: %v", err)
	}

	// Verify generator config has expected annotations
	if _, ok := genConfig.Resources["Role"]; !ok {
		t.Error("generator config missing 'Role' resource")
	}
	if _, ok := genConfig.Resources["Policy"]; !ok {
		t.Error("generator config missing 'Policy' resource")
	}

	// Step 4: Classify fields using classifyFields
	ctrl := types.ControllerRepo{
		RepoName:    "iam-controller",
		ServiceName: "iam",
	}
	results := classifyFields(ctrl, fields, genConfig)

	if len(results) != 5 {
		t.Fatalf("classifyFields returned %d results, want 5", len(results))
	}

	// Verify classification:
	// - AssumeRolePolicyDocument → annotated (is_iam_policy) - matched via Role resource
	// - PolicyDocument → annotated (is_document) - matched via Policy resource
	// - InlinePolicyDocument → unannotated (gap) - no annotation in generator.yaml
	// - PermissionsBoundaryPolicyARN → unannotated (gap) - no annotation
	// - PolicyName → unannotated (gap) - no annotation
	annotationMap := make(map[string]types.AnnotationType)
	for _, r := range results {
		annotationMap[r.FieldName] = r.AnnotationType
	}

	if annotationMap["AssumeRolePolicyDocument"] != types.AnnotationIAMPolicy {
		t.Errorf("AssumeRolePolicyDocument annotation = %q, want %q",
			annotationMap["AssumeRolePolicyDocument"], types.AnnotationIAMPolicy)
	}
	if annotationMap["PolicyDocument"] != types.AnnotationDocument {
		t.Errorf("PolicyDocument annotation = %q, want %q",
			annotationMap["PolicyDocument"], types.AnnotationDocument)
	}
	if annotationMap["InlinePolicyDocument"] != types.AnnotationNone {
		t.Errorf("InlinePolicyDocument annotation = %q, want %q",
			annotationMap["InlinePolicyDocument"], types.AnnotationNone)
	}
	if annotationMap["PermissionsBoundaryPolicyARN"] != types.AnnotationNone {
		t.Errorf("PermissionsBoundaryPolicyARN annotation = %q, want %q",
			annotationMap["PermissionsBoundaryPolicyARN"], types.AnnotationNone)
	}
	if annotationMap["PolicyName"] != types.AnnotationNone {
		t.Errorf("PolicyName annotation = %q, want %q",
			annotationMap["PolicyName"], types.AnnotationNone)
	}
}

// TestIntegration_TerraformScanPipeline tests the Terraform documentation
// parsing pipeline using fixture data.
func TestIntegration_TerraformScanPipeline(t *testing.T) {
	td := testdataDir(t)
	docsDir := filepath.Join(td, "terraform")

	tfParser := &parser.TerraformParser{}
	fields, err := tfParser.ParseAllDocs(docsDir)
	if err != nil {
		t.Fatalf("ParseAllDocs failed: %v", err)
	}

	// The iam_role.html.markdown fixture should yield:
	// - assume_role_policy: detected via both jsonencode() in example AND "JSON string" in description
	// - policy: detected via data.aws_iam_policy_document reference in example
	if len(fields) < 1 {
		t.Fatalf("ParseAllDocs returned %d fields, want at least 1", len(fields))
	}

	// Check that assume_role_policy was detected
	fieldMap := make(map[string]types.TerraformField)
	for _, f := range fields {
		fieldMap[f.FieldName] = f
	}

	if tf, ok := fieldMap["assume_role_policy"]; ok {
		if tf.ServiceName != "iam" {
			t.Errorf("assume_role_policy service = %q, want %q", tf.ServiceName, "iam")
		}
		if tf.ResourceType != "role" {
			t.Errorf("assume_role_policy resource type = %q, want %q", tf.ResourceType, "role")
		}
		// Should be detected by both description phrase and jsonencode
		if tf.DetectionMethod != types.DetectBoth {
			t.Errorf("assume_role_policy detection method = %q, want %q", tf.DetectionMethod, types.DetectBoth)
		}
	} else {
		t.Error("assume_role_policy not found in parsed Terraform fields")
	}

	// Check that "policy" field was detected via data.aws_iam_policy_document reference
	if _, ok := fieldMap["policy"]; ok {
		// The "policy" field is detected via data.aws_iam_policy_document.inline_policy.json
		if fieldMap["policy"].DetectionMethod != types.DetectJsonEncodeExample {
			t.Errorf("policy detection method = %q, want %q",
				fieldMap["policy"].DetectionMethod, types.DetectJsonEncodeExample)
		}
	} else {
		t.Error("policy not found in parsed Terraform fields")
	}
}

// TestIntegration_FullReportPipeline tests the full report pipeline: ACK scan
// results + Terraform results → Matcher → Reporter (JSON output).
func TestIntegration_FullReportPipeline(t *testing.T) {
	// Simulate ACK scan results (as if from scanning iam-controller)
	ackResults := []types.ScanResult{
		{
			ServiceName:    "iam",
			RepoName:       "iam-controller",
			ResourceName:   "RoleSpec",
			FieldName:      "AssumeRolePolicyDocument",
			FieldPath:      "RoleSpec.AssumeRolePolicyDocument",
			GoType:         "*string",
			AnnotationType: types.AnnotationIAMPolicy,
		},
		{
			ServiceName:    "iam",
			RepoName:       "iam-controller",
			ResourceName:   "RoleSpec",
			FieldName:      "InlinePolicyDocument",
			FieldPath:      "RoleSpec.InlinePolicyDocument",
			GoType:         "*string",
			AnnotationType: types.AnnotationNone,
		},
		{
			ServiceName:    "iam",
			RepoName:       "iam-controller",
			ResourceName:   "PolicySpec",
			FieldName:      "PolicyDocument",
			FieldPath:      "PolicySpec.PolicyDocument",
			GoType:         "*string",
			AnnotationType: types.AnnotationDocument,
		},
	}

	// Simulate Terraform scan results
	tfFields := []types.TerraformField{
		{
			ServiceName:     "iam",
			ResourceType:    "role",
			FieldName:       "assume_role_policy",
			Description:     "Policy document associated with the role as a JSON string.",
			DetectionMethod: types.DetectBoth,
		},
		{
			ServiceName:     "iam",
			ResourceType:    "role",
			FieldName:       "inline_policy_document",
			Description:     "Inline policy JSON document.",
			DetectionMethod: types.DetectDescriptionPhrase,
		},
		{
			ServiceName:     "iam",
			ResourceType:    "policy",
			FieldName:       "policy",
			Description:     "The policy document as a JSON formatted string.",
			DetectionMethod: types.DetectDescriptionPhrase,
		},
	}

	// Run matcher
	m := &matcher.Matcher{}
	matchResults := m.Match(ackResults, tfFields)

	// Verify we get results
	if len(matchResults) == 0 {
		t.Fatal("Match returned zero results")
	}

	// Generate JSON report
	rep := reporter.NewReporter("json")
	var buf bytes.Buffer
	err := rep.GenerateReport(matchResults, &buf)
	if err != nil {
		t.Fatalf("GenerateReport failed: %v", err)
	}

	output := buf.String()
	if output == "" {
		t.Fatal("GenerateReport produced empty output")
	}

	// Verify the JSON can be parsed
	var report struct {
		Fields  []json.RawMessage `json:"fields"`
		Summary json.RawMessage   `json:"summary"`
	}
	if err := json.Unmarshal(buf.Bytes(), &report); err != nil {
		t.Fatalf("failed to unmarshal report JSON: %v", err)
	}

	// Verify there are field entries
	if len(report.Fields) == 0 {
		t.Error("report JSON has empty 'fields' array")
	}

	// Verify summary is present
	if report.Summary == nil {
		t.Error("report JSON missing 'summary' object")
	}

	// Verify summary counts are consistent
	var summary types.ReportSummary
	if err := json.Unmarshal(report.Summary, &summary); err != nil {
		t.Fatalf("failed to unmarshal summary: %v", err)
	}
	totalCategories := summary.AnnotatedCount + summary.GapConfirmedCount +
		summary.GapUnconfirmedCount + summary.TerraformOnlyCount
	if totalCategories != summary.TotalFields {
		t.Errorf("summary category sum (%d) != total fields (%d)",
			totalCategories, summary.TotalFields)
	}
}

// TestIntegration_JSONOutputFormat verifies that the JSON reporter output
// conforms to the expected structure with "fields" array and "summary" object.
func TestIntegration_JSONOutputFormat(t *testing.T) {
	// Create representative match results
	matchResults := []types.MatchResult{
		{
			ServiceName:      "iam",
			ResourceName:     "Role",
			FieldName:        "AssumeRolePolicyDocument",
			FieldPath:        "RoleSpec.AssumeRolePolicyDocument",
			AnnotationStatus: types.AnnotationIAMPolicy,
			TFConfirmation:   types.TFNotApplicable,
			Category:         types.CategoryAnnotated,
		},
		{
			ServiceName:      "iam",
			ResourceName:     "Role",
			FieldName:        "InlinePolicyDocument",
			FieldPath:        "RoleSpec.InlinePolicyDocument",
			AnnotationStatus: types.AnnotationNone,
			TFConfirmation:   types.TFConfirmed,
			Category:         types.CategoryGapConfirmed,
		},
		{
			ServiceName:      "s3",
			ResourceName:     "BucketSpec",
			FieldName:        "BucketPolicyDocument",
			FieldPath:        "BucketSpec.BucketPolicyDocument",
			AnnotationStatus: types.AnnotationNone,
			TFConfirmation:   types.TFUnconfirmed,
			Category:         types.CategoryGapUnconfirmed,
		},
	}

	rep := reporter.NewReporter("json")
	var buf bytes.Buffer
	err := rep.GenerateReport(matchResults, &buf)
	if err != nil {
		t.Fatalf("GenerateReport (json) failed: %v", err)
	}

	// Parse and validate structure
	var report struct {
		Fields  []map[string]interface{} `json:"fields"`
		Summary map[string]interface{}   `json:"summary"`
	}
	if err := json.Unmarshal(buf.Bytes(), &report); err != nil {
		t.Fatalf("failed to unmarshal JSON output: %v", err)
	}

	// Verify "fields" array exists and has entries
	if report.Fields == nil {
		t.Fatal("JSON output missing 'fields' key")
	}
	if len(report.Fields) != 3 {
		t.Errorf("expected 3 fields, got %d", len(report.Fields))
	}

	// Verify each field entry has required keys
	requiredFieldKeys := []string{
		"service_name", "resource_name", "field_name", "field_path",
		"annotation_status", "terraform_confirmation", "category",
	}
	for i, field := range report.Fields {
		for _, key := range requiredFieldKeys {
			if _, ok := field[key]; !ok {
				t.Errorf("field[%d] missing required key %q", i, key)
			}
		}
	}

	// Verify "summary" object exists and has required keys
	if report.Summary == nil {
		t.Fatal("JSON output missing 'summary' key")
	}
	requiredSummaryKeys := []string{
		"total_fields", "annotated_count", "gap_confirmed_count",
		"gap_unconfirmed_count", "terraform_only_count",
		"gaps_per_service", "services_by_priority",
	}
	for _, key := range requiredSummaryKeys {
		if _, ok := report.Summary[key]; !ok {
			t.Errorf("summary missing required key %q", key)
		}
	}

	// Verify numeric summary values
	totalFields, _ := report.Summary["total_fields"].(float64)
	annotated, _ := report.Summary["annotated_count"].(float64)
	gapConfirmed, _ := report.Summary["gap_confirmed_count"].(float64)
	gapUnconfirmed, _ := report.Summary["gap_unconfirmed_count"].(float64)
	tfOnly, _ := report.Summary["terraform_only_count"].(float64)

	if int(totalFields) != 3 {
		t.Errorf("total_fields = %v, want 3", totalFields)
	}
	if int(annotated) != 1 {
		t.Errorf("annotated_count = %v, want 1", annotated)
	}
	if int(gapConfirmed) != 1 {
		t.Errorf("gap_confirmed_count = %v, want 1", gapConfirmed)
	}
	if int(gapUnconfirmed) != 1 {
		t.Errorf("gap_unconfirmed_count = %v, want 1", gapUnconfirmed)
	}
	if int(tfOnly) != 0 {
		t.Errorf("terraform_only_count = %v, want 0", tfOnly)
	}
}

// TestIntegration_MarkdownOutputFormat verifies that the markdown reporter
// produces output with expected headings and table markers.
func TestIntegration_MarkdownOutputFormat(t *testing.T) {
	matchResults := []types.MatchResult{
		{
			ServiceName:      "iam",
			ResourceName:     "Role",
			FieldName:        "AssumeRolePolicyDocument",
			FieldPath:        "RoleSpec.AssumeRolePolicyDocument",
			AnnotationStatus: types.AnnotationIAMPolicy,
			TFConfirmation:   types.TFNotApplicable,
			Category:         types.CategoryAnnotated,
		},
		{
			ServiceName:      "iam",
			ResourceName:     "Role",
			FieldName:        "InlinePolicyDocument",
			FieldPath:        "RoleSpec.InlinePolicyDocument",
			AnnotationStatus: types.AnnotationNone,
			TFConfirmation:   types.TFConfirmed,
			Category:         types.CategoryGapConfirmed,
		},
		{
			ServiceName:      "s3",
			ResourceName:     "BucketSpec",
			FieldName:        "BucketPolicyDocument",
			FieldPath:        "BucketSpec.BucketPolicyDocument",
			AnnotationStatus: types.AnnotationNone,
			TFConfirmation:   types.TFUnconfirmed,
			Category:         types.CategoryGapUnconfirmed,
		},
	}

	rep := reporter.NewReporter("markdown")
	var buf bytes.Buffer
	err := rep.GenerateReport(matchResults, &buf)
	if err != nil {
		t.Fatalf("GenerateReport (markdown) failed: %v", err)
	}

	output := buf.String()

	// Verify main heading
	if !strings.Contains(output, "# ACK Scanner Gap Analysis Report") {
		t.Error("markdown output missing '# ACK Scanner Gap Analysis Report' heading")
	}

	// Verify Summary heading
	if !strings.Contains(output, "## Summary") {
		t.Error("markdown output missing '## Summary' heading")
	}

	// Verify Priority Services heading
	if !strings.Contains(output, "## Priority Services") {
		t.Error("markdown output missing '## Priority Services' heading")
	}

	// Verify service-level headings
	if !strings.Contains(output, "## iam") {
		t.Error("markdown output missing '## iam' service heading")
	}
	if !strings.Contains(output, "## s3") {
		t.Error("markdown output missing '## s3' service heading")
	}

	// Verify table markers (header separator)
	if !strings.Contains(output, "| --- |") {
		t.Error("markdown output missing table separator markers")
	}

	// Verify table headers for field table
	if !strings.Contains(output, "| Resource | Field | Annotation | TF Confirmed | Category |") {
		t.Error("markdown output missing field table header row")
	}

	// Verify gap fields are highlighted with warning emoji
	if !strings.Contains(output, "⚠️ InlinePolicyDocument") {
		t.Error("markdown output missing gap highlight for InlinePolicyDocument")
	}
	if !strings.Contains(output, "⚠️ BucketPolicyDocument") {
		t.Error("markdown output missing gap highlight for BucketPolicyDocument")
	}

	// Verify summary table has metrics
	if !strings.Contains(output, "| Total Fields |") {
		t.Error("markdown output missing 'Total Fields' in summary table")
	}
	if !strings.Contains(output, "| Already Annotated |") {
		t.Error("markdown output missing 'Already Annotated' in summary table")
	}
	if !strings.Contains(output, "| Gaps (Terraform Confirmed) |") {
		t.Error("markdown output missing 'Gaps (Terraform Confirmed)' in summary table")
	}
}
