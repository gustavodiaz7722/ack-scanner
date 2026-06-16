package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aws-controllers-k8s/ack-scanner/pkg/types"
)

// TestInvalidOutputFlag_ErrorContainsValidOptions verifies that an invalid --output
// flag value produces an error message listing all valid options (json, table, markdown).
// Validates: Requirements 6.10
func TestInvalidOutputFlag_ErrorContainsValidOptions(t *testing.T) {
	invalidValues := []string{"yaml", "csv", "xml", "HTML", "plaintext", ""}

	for _, val := range invalidValues {
		t.Run("value="+val, func(t *testing.T) {
			githubToken = "token"
			cacheDir = t.TempDir()
			outputFormat = val

			err := persistentPreRun(rootCmd, nil)
			if err == nil {
				t.Fatalf("expected error for invalid output format %q, got nil", val)
			}

			errMsg := err.Error()
			for _, valid := range []string{"json", "table", "markdown"} {
				if !strings.Contains(errMsg, valid) {
					t.Errorf("error message %q does not list valid option %q", errMsg, valid)
				}
			}
		})
	}
}

// TestFindTypesFile_MissingApisDir verifies that findTypesFile returns an error
// when the apis/ directory does not exist.
// Validates: Requirements 1.3, 1.4
func TestFindTypesFile_MissingApisDir(t *testing.T) {
	repoDir := t.TempDir()
	// Do NOT create an apis/ directory.

	_, err := findTypesFile(repoDir)
	if err == nil {
		t.Fatal("expected error for missing apis/ directory, got nil")
	}

	if !strings.Contains(err.Error(), "apis/") {
		t.Errorf("error message %q does not mention apis/", err.Error())
	}
}

// TestFindTypesFile_EmptyApisDir verifies that findTypesFile returns an error
// when the apis/ directory exists but contains no version subdirectories.
// Validates: Requirements 1.3, 1.4
func TestFindTypesFile_EmptyApisDir(t *testing.T) {
	repoDir := t.TempDir()
	apisDir := filepath.Join(repoDir, "apis")
	if err := os.MkdirAll(apisDir, 0o755); err != nil {
		t.Fatalf("failed to create apis/ dir: %v", err)
	}

	_, err := findTypesFile(repoDir)
	if err == nil {
		t.Fatal("expected error for empty apis/ directory, got nil")
	}

	if !strings.Contains(err.Error(), "version") {
		t.Errorf("error message %q does not mention 'version'", err.Error())
	}
}

// TestFindTypesFile_MissingTypesGo verifies that findTypesFile returns an error
// when the apis/{version}/ directory exists but types.go is absent.
// Validates: Requirements 1.3, 1.4
func TestFindTypesFile_MissingTypesGo(t *testing.T) {
	repoDir := t.TempDir()
	versionDir := filepath.Join(repoDir, "apis", "v1alpha1")
	if err := os.MkdirAll(versionDir, 0o755); err != nil {
		t.Fatalf("failed to create version dir: %v", err)
	}
	// Do NOT create types.go in the version directory.

	_, err := findTypesFile(repoDir)
	if err == nil {
		t.Fatal("expected error for missing types.go, got nil")
	}

	if !strings.Contains(err.Error(), "types.go") {
		t.Errorf("error message %q does not mention 'types.go'", err.Error())
	}
}

// TestFindTypesFile_SelectsLatestVersion verifies that findTypesFile picks the
// latest (alphabetically last) version directory when multiple exist.
func TestFindTypesFile_SelectsLatestVersion(t *testing.T) {
	repoDir := t.TempDir()

	// Create multiple version directories.
	versions := []string{"v1alpha1", "v1alpha2", "v1beta1"}
	for _, v := range versions {
		vDir := filepath.Join(repoDir, "apis", v)
		if err := os.MkdirAll(vDir, 0o755); err != nil {
			t.Fatalf("failed to create dir %s: %v", v, err)
		}
	}

	// Only create types.go in the latest version.
	latestTypesFile := filepath.Join(repoDir, "apis", "v1beta1", "types.go")
	if err := os.WriteFile(latestTypesFile, []byte("package v1beta1\n"), 0o644); err != nil {
		t.Fatalf("failed to write types.go: %v", err)
	}

	result, err := findTypesFile(repoDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != latestTypesFile {
		t.Errorf("expected %q, got %q", latestTypesFile, result)
	}
}

// TestLookupAnnotation_NilConfig verifies that lookupAnnotation returns
// AnnotationNone when the config is nil.
// Validates: Requirements 6.8
func TestLookupAnnotation_NilConfig(t *testing.T) {
	field := types.GoField{
		StructName: "Bucket",
		FieldName:  "PolicyDocument",
		FieldPath:  "Spec.PolicyDocument",
		GoType:     "*string",
	}

	result := lookupAnnotation(field, nil)
	if result != types.AnnotationNone {
		t.Errorf("expected AnnotationNone for nil config, got %q", result)
	}
}

// TestLookupAnnotation_SpecSuffixStripping verifies that lookupAnnotation
// correctly strips the "Spec" suffix from struct names to match resource names
// in the generator config.
// Validates: Requirements 6.8
func TestLookupAnnotation_SpecSuffixStripping(t *testing.T) {
	config := &types.GeneratorConfig{
		Resources: map[string]types.ResourceConfig{
			"Bucket": {
				Fields: map[string]types.FieldConfig{
					"PolicyDocument": {IsDocument: true},
				},
			},
		},
	}

	// Field with "Spec" suffix in struct name should match "Bucket" resource.
	field := types.GoField{
		StructName: "BucketSpec",
		FieldName:  "PolicyDocument",
		FieldPath:  "Spec.PolicyDocument",
		GoType:     "*string",
	}

	result := lookupAnnotation(field, config)
	if result != types.AnnotationDocument {
		t.Errorf("expected AnnotationDocument after Spec suffix stripping, got %q", result)
	}
}

// TestLookupAnnotation_SDKSuffixStripping verifies that lookupAnnotation
// correctly strips the "_SDK" suffix from struct names.
func TestLookupAnnotation_SDKSuffixStripping(t *testing.T) {
	config := &types.GeneratorConfig{
		Resources: map[string]types.ResourceConfig{
			"Role": {
				Fields: map[string]types.FieldConfig{
					"AssumeRolePolicyDocument": {IsIAMPolicy: true},
				},
			},
		},
	}

	field := types.GoField{
		StructName: "Role_SDK",
		FieldName:  "AssumeRolePolicyDocument",
		FieldPath:  "Spec.AssumeRolePolicyDocument",
		GoType:     "*string",
	}

	result := lookupAnnotation(field, config)
	if result != types.AnnotationIAMPolicy {
		t.Errorf("expected AnnotationIAMPolicy after _SDK suffix stripping, got %q", result)
	}
}

// TestLookupAnnotation_FieldPathFallback verifies that lookupAnnotation falls
// back to FieldPath lookup when direct FieldName match fails.
func TestLookupAnnotation_FieldPathFallback(t *testing.T) {
	config := &types.GeneratorConfig{
		Resources: map[string]types.ResourceConfig{
			"Topic": {
				Fields: map[string]types.FieldConfig{
					"Spec.FilterPolicy": {IsDocument: true},
				},
			},
		},
	}

	field := types.GoField{
		StructName: "Topic",
		FieldName:  "FilterPolicy",
		FieldPath:  "Spec.FilterPolicy",
		GoType:     "*string",
	}

	result := lookupAnnotation(field, config)
	if result != types.AnnotationDocument {
		t.Errorf("expected AnnotationDocument via FieldPath fallback, got %q", result)
	}
}

// TestLookupAnnotation_NoMatch verifies that lookupAnnotation returns
// AnnotationNone when neither the resource nor the field is found.
func TestLookupAnnotation_NoMatch(t *testing.T) {
	config := &types.GeneratorConfig{
		Resources: map[string]types.ResourceConfig{
			"Bucket": {
				Fields: map[string]types.FieldConfig{
					"PolicyDocument": {IsDocument: true},
				},
			},
		},
	}

	// Field doesn't match any resource in config.
	field := types.GoField{
		StructName: "UnknownResource",
		FieldName:  "UnknownField",
		FieldPath:  "Spec.UnknownField",
		GoType:     "*string",
	}

	result := lookupAnnotation(field, config)
	if result != types.AnnotationNone {
		t.Errorf("expected AnnotationNone for unmatched field, got %q", result)
	}
}

// TestFilterControllers_EmptyFilter verifies that filterControllers with an
// empty filter string returns an empty list.
// Validates: Requirements 4.4
func TestFilterControllers_EmptyFilter(t *testing.T) {
	controllers := []types.ControllerRepo{
		{RepoName: "s3-controller", ServiceName: "s3"},
		{RepoName: "ec2-controller", ServiceName: "ec2"},
	}

	result := filterControllers(controllers, "")
	if len(result) != 0 {
		t.Errorf("expected empty result for empty filter, got %d items", len(result))
	}
}

// TestFilterControllers_CorrectFiltering verifies that filterControllers
// correctly filters by exact service name match.
// Validates: Requirements 4.4
func TestFilterControllers_CorrectFiltering(t *testing.T) {
	controllers := []types.ControllerRepo{
		{RepoName: "s3-controller", ServiceName: "s3"},
		{RepoName: "ec2-controller", ServiceName: "ec2"},
		{RepoName: "iam-controller", ServiceName: "iam"},
		{RepoName: "eks-controller", ServiceName: "eks"},
	}

	result := filterControllers(controllers, "s3,iam")
	if len(result) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result))
	}

	// Check that the correct services were kept.
	services := make(map[string]bool)
	for _, r := range result {
		services[r.ServiceName] = true
	}
	if !services["s3"] {
		t.Error("expected s3 in filtered results")
	}
	if !services["iam"] {
		t.Error("expected iam in filtered results")
	}
	if services["ec2"] {
		t.Error("ec2 should not be in filtered results")
	}
	if services["eks"] {
		t.Error("eks should not be in filtered results")
	}
}

// TestFilterControllers_WhitespaceInFilter verifies that filterControllers
// correctly trims whitespace from filter values.
func TestFilterControllers_WhitespaceInFilter(t *testing.T) {
	controllers := []types.ControllerRepo{
		{RepoName: "s3-controller", ServiceName: "s3"},
		{RepoName: "ec2-controller", ServiceName: "ec2"},
	}

	result := filterControllers(controllers, " s3 , ec2 ")
	if len(result) != 2 {
		t.Errorf("expected 2 results with whitespace-trimmed filter, got %d", len(result))
	}
}

// TestFilterControllers_NoMatchingServices verifies that filterControllers
// returns empty when no controllers match the filter.
func TestFilterControllers_NoMatchingServices(t *testing.T) {
	controllers := []types.ControllerRepo{
		{RepoName: "s3-controller", ServiceName: "s3"},
		{RepoName: "ec2-controller", ServiceName: "ec2"},
	}

	result := filterControllers(controllers, "rds,lambda")
	if len(result) != 0 {
		t.Errorf("expected empty result for non-matching filter, got %d items", len(result))
	}
}

// TestScanResultsToMatchResults_Conversion verifies that scanResultsToMatchResults
// correctly converts ScanResult entries to MatchResult entries with proper
// category and TerraformConfirmation assignments.
// Validates: Requirements 6.8
func TestScanResultsToMatchResults_Conversion(t *testing.T) {
	results := []types.ScanResult{
		{
			ServiceName:    "s3",
			RepoName:       "s3-controller",
			ResourceName:   "Bucket",
			FieldName:      "PolicyDocument",
			FieldPath:      "Spec.PolicyDocument",
			GoType:         "*string",
			AnnotationType: types.AnnotationDocument,
		},
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

	matchResults := scanResultsToMatchResults(results)
	if len(matchResults) != 2 {
		t.Fatalf("expected 2 match results, got %d", len(matchResults))
	}

	// First result: annotated → CategoryAnnotated, TFNotApplicable
	mr0 := matchResults[0]
	if mr0.Category != types.CategoryAnnotated {
		t.Errorf("expected CategoryAnnotated for annotated field, got %q", mr0.Category)
	}
	if mr0.TFConfirmation != types.TFNotApplicable {
		t.Errorf("expected TFNotApplicable for annotated field, got %q", mr0.TFConfirmation)
	}
	if mr0.AnnotationStatus != types.AnnotationDocument {
		t.Errorf("expected AnnotationDocument status, got %q", mr0.AnnotationStatus)
	}
	if mr0.ServiceName != "s3" {
		t.Errorf("expected service s3, got %q", mr0.ServiceName)
	}

	// Second result: unannotated → CategoryGapUnconfirmed, TFUnconfirmed
	mr1 := matchResults[1]
	if mr1.Category != types.CategoryGapUnconfirmed {
		t.Errorf("expected CategoryGapUnconfirmed for unannotated field, got %q", mr1.Category)
	}
	if mr1.TFConfirmation != types.TFUnconfirmed {
		t.Errorf("expected TFUnconfirmed for unannotated field, got %q", mr1.TFConfirmation)
	}
	if mr1.AnnotationStatus != types.AnnotationNone {
		t.Errorf("expected AnnotationNone status, got %q", mr1.AnnotationStatus)
	}
	if mr1.ServiceName != "iam" {
		t.Errorf("expected service iam, got %q", mr1.ServiceName)
	}
}

// TestScanResultsToMatchResults_EmptyInput verifies that scanResultsToMatchResults
// handles an empty input slice gracefully.
func TestScanResultsToMatchResults_EmptyInput(t *testing.T) {
	matchResults := scanResultsToMatchResults(nil)
	if matchResults != nil && len(matchResults) != 0 {
		t.Errorf("expected nil or empty result for nil input, got %d items", len(matchResults))
	}
}

// TestScanResultsToMatchResults_IAMPolicyAnnotation verifies that IAM policy
// annotated fields are correctly classified.
func TestScanResultsToMatchResults_IAMPolicyAnnotation(t *testing.T) {
	results := []types.ScanResult{
		{
			ServiceName:    "iam",
			RepoName:       "iam-controller",
			ResourceName:   "Role",
			FieldName:      "AssumeRolePolicyDocument",
			FieldPath:      "Spec.AssumeRolePolicyDocument",
			GoType:         "*string",
			AnnotationType: types.AnnotationIAMPolicy,
		},
	}

	matchResults := scanResultsToMatchResults(results)
	if len(matchResults) != 1 {
		t.Fatalf("expected 1 match result, got %d", len(matchResults))
	}

	mr := matchResults[0]
	if mr.Category != types.CategoryAnnotated {
		t.Errorf("expected CategoryAnnotated for IAM policy annotated field, got %q", mr.Category)
	}
	if mr.AnnotationStatus != types.AnnotationIAMPolicy {
		t.Errorf("expected AnnotationIAMPolicy status, got %q", mr.AnnotationStatus)
	}
}
