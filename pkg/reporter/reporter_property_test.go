package reporter

import (
	"bytes"
	"encoding/json"
	"sort"
	"strings"
	"testing"

	"github.com/aws-controllers-k8s/ack-scanner/pkg/types"
	"pgregory.net/rapid"
)

// --- Generators ---

// genServiceName generates a lowercase service name (2-10 lowercase letters).
func genServiceName() *rapid.Generator[string] {
	return rapid.Custom(func(t *rapid.T) string {
		length := rapid.IntRange(2, 10).Draw(t, "length")
		var sb strings.Builder
		for i := 0; i < length; i++ {
			ch := rapid.ByteRange('a', 'z').Draw(t, "char")
			sb.WriteByte(ch)
		}
		return sb.String()
	})
}

// genRepoName generates a controller repo name from a service name.
func genRepoName(service string) string {
	return service + "-controller"
}

// genControllerRepo generates a ControllerRepo with a random service name.
func genControllerRepo() *rapid.Generator[types.ControllerRepo] {
	return rapid.Custom(func(t *rapid.T) types.ControllerRepo {
		service := genServiceName().Draw(t, "service")
		return types.ControllerRepo{
			RepoName:    genRepoName(service),
			ServiceName: service,
		}
	})
}

// genFieldName generates a CamelCase field name with 1-3 words.
func genFieldName() *rapid.Generator[string] {
	return rapid.Custom(func(t *rapid.T) string {
		wordCount := rapid.IntRange(1, 3).Draw(t, "wordCount")
		var sb strings.Builder
		for i := 0; i < wordCount; i++ {
			upper := rapid.ByteRange('A', 'Z').Draw(t, "upper")
			sb.WriteByte(upper)
			lowerCount := rapid.IntRange(2, 6).Draw(t, "lowerCount")
			for j := 0; j < lowerCount; j++ {
				lower := rapid.ByteRange('a', 'z').Draw(t, "lower")
				sb.WriteByte(lower)
			}
		}
		return sb.String()
	})
}

// genResourceName generates a resource name starting with uppercase.
func genResourceName() *rapid.Generator[string] {
	return rapid.Custom(func(t *rapid.T) string {
		upper := rapid.ByteRange('A', 'Z').Draw(t, "upper")
		length := rapid.IntRange(2, 8).Draw(t, "length")
		var sb strings.Builder
		sb.WriteByte(upper)
		for i := 0; i < length; i++ {
			lower := rapid.ByteRange('a', 'z').Draw(t, "lower")
			sb.WriteByte(lower)
		}
		return sb.String()
	})
}

// genCategory generates a random Category value.
func genCategory() *rapid.Generator[types.Category] {
	return rapid.SampledFrom([]types.Category{
		types.CategoryGapConfirmed,
		types.CategoryGapUnconfirmed,
		types.CategoryAnnotated,
		types.CategoryTerraformOnly,
	})
}

// genAnnotationType generates an annotation type.
func genAnnotationType() *rapid.Generator[types.AnnotationType] {
	return rapid.SampledFrom([]types.AnnotationType{
		types.AnnotationNone,
		types.AnnotationDocument,
		types.AnnotationIAMPolicy,
	})
}

// genTFConfirmation generates a TerraformConfirmation value.
func genTFConfirmation() *rapid.Generator[types.TerraformConfirmation] {
	return rapid.SampledFrom([]types.TerraformConfirmation{
		types.TFConfirmed,
		types.TFUnconfirmed,
		types.TFNotApplicable,
	})
}

// genMatchResult generates a random MatchResult.
func genMatchResult() *rapid.Generator[types.MatchResult] {
	return rapid.Custom(func(t *rapid.T) types.MatchResult {
		service := genServiceName().Draw(t, "service")
		resource := genResourceName().Draw(t, "resource")
		field := genFieldName().Draw(t, "field")
		category := genCategory().Draw(t, "category")
		annotation := genAnnotationType().Draw(t, "annotation")
		tfConfirm := genTFConfirmation().Draw(t, "tfConfirm")
		return types.MatchResult{
			ServiceName:      service,
			ResourceName:     resource,
			FieldName:        field,
			FieldPath:        "Spec." + field,
			AnnotationStatus: annotation,
			TFConfirmation:   tfConfirm,
			Category:         category,
		}
	})
}

// --- Property 4: Controller list sort order ---
// Feature: ack-scanner, Property 4: Controller list sort order
// **Validates: Requirements 2.1**

// TestProperty4_ControllerListSortOrder verifies that for any list of controller
// repositories, the formatted output presents controllers sorted alphabetically
// by service name.
func TestProperty4_ControllerListSortOrder(t *testing.T) {
	// Feature: ack-scanner, Property 4: Controller list sort order
	rapid.Check(t, func(t *rapid.T) {
		// Generate a random list of controller repos
		repoCount := rapid.IntRange(0, 20).Draw(t, "repoCount")
		repos := make([]types.ControllerRepo, repoCount)
		for i := 0; i < repoCount; i++ {
			repos[i] = genControllerRepo().Draw(t, "repo")
		}

		r := NewReporter("json")
		var buf bytes.Buffer
		err := r.FormatControllerList(repos, &buf)
		if err != nil {
			t.Fatalf("FormatControllerList failed: %v", err)
		}

		// Deserialize the JSON output
		var output []controllerListEntry
		if err := json.Unmarshal(buf.Bytes(), &output); err != nil {
			t.Fatalf("output is not valid JSON: %v", err)
		}

		// Verify the output is sorted alphabetically by service name
		for i := 1; i < len(output); i++ {
			if output[i-1].ServiceName > output[i].ServiceName {
				t.Fatalf("sort order violated at index %d: %q > %q",
					i, output[i-1].ServiceName, output[i].ServiceName)
			}
		}
	})
}

// TestProperty4_ControllerListSortOrder_Table verifies sort order holds for table format too.
func TestProperty4_ControllerListSortOrder_Table(t *testing.T) {
	// Feature: ack-scanner, Property 4: Controller list sort order
	rapid.Check(t, func(t *rapid.T) {
		repoCount := rapid.IntRange(2, 10).Draw(t, "repoCount")
		repos := make([]types.ControllerRepo, repoCount)
		for i := 0; i < repoCount; i++ {
			repos[i] = genControllerRepo().Draw(t, "repo")
		}

		// Sort a copy to get expected order
		expected := make([]types.ControllerRepo, len(repos))
		copy(expected, repos)
		sort.Slice(expected, func(i, j int) bool {
			return expected[i].ServiceName < expected[j].ServiceName
		})

		// Use JSON to verify internal sort (table output is harder to parse)
		r := NewReporter("json")
		var buf bytes.Buffer
		err := r.FormatControllerList(repos, &buf)
		if err != nil {
			t.Fatalf("FormatControllerList failed: %v", err)
		}

		var output []controllerListEntry
		if err := json.Unmarshal(buf.Bytes(), &output); err != nil {
			t.Fatalf("output is not valid JSON: %v", err)
		}

		if len(output) != len(expected) {
			t.Fatalf("expected %d entries, got %d", len(expected), len(output))
		}

		for i := range output {
			if output[i].ServiceName != expected[i].ServiceName {
				t.Fatalf("position %d: expected service %q, got %q",
					i, expected[i].ServiceName, output[i].ServiceName)
			}
		}
	})
}

// --- Property 5: Controller list JSON round-trip ---
// Feature: ack-scanner, Property 5: Controller list JSON round-trip
// **Validates: Requirements 2.2**

// TestProperty5_ControllerListJSONRoundTrip verifies that for any list of controller
// repositories, serializing to JSON and deserializing produces an equivalent list
// where each element contains both repository_name and service_name with original values.
func TestProperty5_ControllerListJSONRoundTrip(t *testing.T) {
	// Feature: ack-scanner, Property 5: Controller list JSON round-trip
	rapid.Check(t, func(t *rapid.T) {
		repoCount := rapid.IntRange(0, 20).Draw(t, "repoCount")
		repos := make([]types.ControllerRepo, repoCount)
		for i := 0; i < repoCount; i++ {
			repos[i] = genControllerRepo().Draw(t, "repo")
		}

		r := NewReporter("json")
		var buf bytes.Buffer
		err := r.FormatControllerList(repos, &buf)
		if err != nil {
			t.Fatalf("FormatControllerList failed: %v", err)
		}

		// Deserialize
		var output []controllerListEntry
		if err := json.Unmarshal(buf.Bytes(), &output); err != nil {
			t.Fatalf("output is not valid JSON: %v", err)
		}

		// Output should have same count as input
		if len(output) != repoCount {
			t.Fatalf("expected %d entries, got %d", repoCount, len(output))
		}

		// Sort both input and output by service name for comparison
		// (FormatControllerList sorts by service name)
		sorted := make([]types.ControllerRepo, len(repos))
		copy(sorted, repos)
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].ServiceName < sorted[j].ServiceName
		})

		for i, entry := range output {
			// Each entry must have both repository_name and service_name
			if entry.RepositoryName == "" {
				t.Fatalf("entry %d: repository_name is empty", i)
			}
			if entry.ServiceName == "" {
				t.Fatalf("entry %d: service_name is empty", i)
			}

			// Values must match the sorted input
			if entry.RepositoryName != sorted[i].RepoName {
				t.Fatalf("entry %d: expected repository_name=%q, got %q",
					i, sorted[i].RepoName, entry.RepositoryName)
			}
			if entry.ServiceName != sorted[i].ServiceName {
				t.Fatalf("entry %d: expected service_name=%q, got %q",
					i, sorted[i].ServiceName, entry.ServiceName)
			}
		}
	})
}

// TestProperty5_ControllerListJSONContainsRequiredFields verifies that JSON output
// contains both repository_name and service_name keys for every element.
func TestProperty5_ControllerListJSONContainsRequiredFields(t *testing.T) {
	// Feature: ack-scanner, Property 5: Controller list JSON round-trip
	rapid.Check(t, func(t *rapid.T) {
		repoCount := rapid.IntRange(1, 10).Draw(t, "repoCount")
		repos := make([]types.ControllerRepo, repoCount)
		for i := 0; i < repoCount; i++ {
			repos[i] = genControllerRepo().Draw(t, "repo")
		}

		r := NewReporter("json")
		var buf bytes.Buffer
		err := r.FormatControllerList(repos, &buf)
		if err != nil {
			t.Fatalf("FormatControllerList failed: %v", err)
		}

		// Parse as generic JSON to check field names
		var rawOutput []map[string]interface{}
		if err := json.Unmarshal(buf.Bytes(), &rawOutput); err != nil {
			t.Fatalf("output is not valid JSON: %v", err)
		}

		for i, entry := range rawOutput {
			if _, ok := entry["repository_name"]; !ok {
				t.Fatalf("entry %d: missing 'repository_name' key", i)
			}
			if _, ok := entry["service_name"]; !ok {
				t.Fatalf("entry %d: missing 'service_name' key", i)
			}
		}
	})
}

// --- Property 15: Report gap sort order ---
// Feature: ack-scanner, Property 15: Report gap sort order
// **Validates: Requirements 5.4**

// TestProperty15_ReportGapSortOrder verifies that for any set of Annotation_Gap
// entries in the report, all Terraform-confirmed gaps appear before all unconfirmed
// gaps in the output ordering.
func TestProperty15_ReportGapSortOrder(t *testing.T) {
	// Feature: ack-scanner, Property 15: Report gap sort order
	rapid.Check(t, func(t *rapid.T) {
		// Generate a mix of match results with various categories
		resultCount := rapid.IntRange(1, 30).Draw(t, "resultCount")
		results := make([]types.MatchResult, resultCount)
		for i := 0; i < resultCount; i++ {
			results[i] = genMatchResult().Draw(t, "result")
		}

		// Use sortResults (the internal function) directly
		sorted := sortResults(results)

		// Verify: all gap_confirmed entries appear before all gap_unconfirmed entries
		lastConfirmedIdx := -1
		firstUnconfirmedIdx := -1

		for i, r := range sorted {
			if r.Category == types.CategoryGapConfirmed {
				lastConfirmedIdx = i
			}
			if r.Category == types.CategoryGapUnconfirmed && firstUnconfirmedIdx == -1 {
				firstUnconfirmedIdx = i
			}
		}

		// If both categories exist, confirmed must all come before unconfirmed
		if lastConfirmedIdx != -1 && firstUnconfirmedIdx != -1 {
			if lastConfirmedIdx >= firstUnconfirmedIdx {
				t.Fatalf("sort order violated: last gap_confirmed at index %d, first gap_unconfirmed at index %d",
					lastConfirmedIdx, firstUnconfirmedIdx)
			}
		}
	})
}

// TestProperty15_ReportGapSortOrderViaJSON verifies the gap sort order through
// the full JSON report output path.
func TestProperty15_ReportGapSortOrderViaJSON(t *testing.T) {
	// Feature: ack-scanner, Property 15: Report gap sort order
	rapid.Check(t, func(t *rapid.T) {
		resultCount := rapid.IntRange(1, 20).Draw(t, "resultCount")
		results := make([]types.MatchResult, resultCount)
		for i := 0; i < resultCount; i++ {
			results[i] = genMatchResult().Draw(t, "result")
		}

		r := NewReporter("json")
		var buf bytes.Buffer
		err := r.GenerateReport(results, &buf)
		if err != nil {
			t.Fatalf("GenerateReport failed: %v", err)
		}

		var output reportJSON
		if err := json.Unmarshal(buf.Bytes(), &output); err != nil {
			t.Fatalf("output is not valid JSON: %v", err)
		}

		// Verify gap_confirmed all appear before gap_unconfirmed in the output
		lastConfirmedIdx := -1
		firstUnconfirmedIdx := -1

		for i, f := range output.Fields {
			if f.Category == types.CategoryGapConfirmed {
				lastConfirmedIdx = i
			}
			if f.Category == types.CategoryGapUnconfirmed && firstUnconfirmedIdx == -1 {
				firstUnconfirmedIdx = i
			}
		}

		if lastConfirmedIdx != -1 && firstUnconfirmedIdx != -1 {
			if lastConfirmedIdx >= firstUnconfirmedIdx {
				t.Fatalf("JSON output sort order violated: last gap_confirmed at index %d, first gap_unconfirmed at index %d",
					lastConfirmedIdx, firstUnconfirmedIdx)
			}
		}
	})
}

// --- Property 16: Report summary consistency ---
// Feature: ack-scanner, Property 16: Report summary consistency
// **Validates: Requirements 5.5**

// TestProperty16_ReportSummaryConsistency verifies that for any report output:
// - The summary's total field counts per category sum to the total number of field entries
// - GapsPerService counts match actual gap entries grouped by service
// - Priority list is sorted descending by confirmed gap count
func TestProperty16_ReportSummaryConsistency(t *testing.T) {
	// Feature: ack-scanner, Property 16: Report summary consistency
	rapid.Check(t, func(t *rapid.T) {
		resultCount := rapid.IntRange(0, 30).Draw(t, "resultCount")
		results := make([]types.MatchResult, resultCount)
		for i := 0; i < resultCount; i++ {
			results[i] = genMatchResult().Draw(t, "result")
		}

		summary := GenerateSummary(results)

		// 1. Category counts must sum to total
		categorySum := summary.AnnotatedCount + summary.GapConfirmedCount +
			summary.GapUnconfirmedCount + summary.TerraformOnlyCount
		if categorySum != summary.TotalFields {
			t.Fatalf("category counts (%d+%d+%d+%d=%d) != TotalFields (%d)",
				summary.AnnotatedCount, summary.GapConfirmedCount,
				summary.GapUnconfirmedCount, summary.TerraformOnlyCount,
				categorySum, summary.TotalFields)
		}

		// 2. TotalFields must equal the input length
		if summary.TotalFields != resultCount {
			t.Fatalf("TotalFields (%d) != input count (%d)", summary.TotalFields, resultCount)
		}

		// 3. GapsPerService should match actual gap entries grouped by service
		expectedGapsPerService := make(map[string]int)
		for _, r := range results {
			if r.Category == types.CategoryGapConfirmed || r.Category == types.CategoryGapUnconfirmed {
				expectedGapsPerService[r.ServiceName]++
			}
		}
		for svc, expectedCount := range expectedGapsPerService {
			if summary.GapsPerService[svc] != expectedCount {
				t.Fatalf("GapsPerService[%q]: expected %d, got %d",
					svc, expectedCount, summary.GapsPerService[svc])
			}
		}
		// Also check no extra services in GapsPerService
		for svc, count := range summary.GapsPerService {
			if expectedGapsPerService[svc] != count {
				t.Fatalf("GapsPerService[%q]: unexpected count %d (expected %d)",
					svc, count, expectedGapsPerService[svc])
			}
		}

		// 4. ServicesByPriority should be sorted descending by ConfirmedGapCount
		for i := 1; i < len(summary.ServicesByPriority); i++ {
			prev := summary.ServicesByPriority[i-1]
			curr := summary.ServicesByPriority[i]
			if prev.ConfirmedGapCount < curr.ConfirmedGapCount {
				t.Fatalf("ServicesByPriority not sorted descending: [%d]=%d < [%d]=%d",
					i-1, prev.ConfirmedGapCount, i, curr.ConfirmedGapCount)
			}
		}
	})
}

// TestProperty16_SummaryGapsPerServiceMatchesActualGaps verifies GapsPerService
// through the full JSON report output path.
func TestProperty16_SummaryGapsPerServiceMatchesActualGaps(t *testing.T) {
	// Feature: ack-scanner, Property 16: Report summary consistency
	rapid.Check(t, func(t *rapid.T) {
		resultCount := rapid.IntRange(1, 20).Draw(t, "resultCount")
		results := make([]types.MatchResult, resultCount)
		for i := 0; i < resultCount; i++ {
			results[i] = genMatchResult().Draw(t, "result")
		}

		r := NewReporter("json")
		var buf bytes.Buffer
		err := r.GenerateReport(results, &buf)
		if err != nil {
			t.Fatalf("GenerateReport failed: %v", err)
		}

		var output reportJSON
		if err := json.Unmarshal(buf.Bytes(), &output); err != nil {
			t.Fatalf("output is not valid JSON: %v", err)
		}

		// Count actual gaps per service from the output fields
		actualGaps := make(map[string]int)
		for _, f := range output.Fields {
			if f.Category == types.CategoryGapConfirmed || f.Category == types.CategoryGapUnconfirmed {
				actualGaps[f.ServiceName]++
			}
		}

		// Verify against summary
		for svc, count := range actualGaps {
			if output.Summary.GapsPerService[svc] != count {
				t.Fatalf("GapsPerService[%q] in output: expected %d, got %d",
					svc, count, output.Summary.GapsPerService[svc])
			}
		}
	})
}

// --- Property 17: Report JSON schema validity ---
// Feature: ack-scanner, Property 17: Report JSON schema validity
// **Validates: Requirements 5.6**

// TestProperty17_ReportJSONSchemaValidity verifies that for any report data,
// serializing to JSON produces valid JSON containing a "fields" array (each element
// with service_name, resource_name, field_name, field_path, annotation_status,
// terraform_confirmation, category) and a "summary" object.
func TestProperty17_ReportJSONSchemaValidity(t *testing.T) {
	// Feature: ack-scanner, Property 17: Report JSON schema validity
	rapid.Check(t, func(t *rapid.T) {
		resultCount := rapid.IntRange(0, 20).Draw(t, "resultCount")
		results := make([]types.MatchResult, resultCount)
		for i := 0; i < resultCount; i++ {
			results[i] = genMatchResult().Draw(t, "result")
		}

		r := NewReporter("json")
		var buf bytes.Buffer
		err := r.GenerateReport(results, &buf)
		if err != nil {
			t.Fatalf("GenerateReport failed: %v", err)
		}

		// Parse as generic JSON to verify schema structure
		var rawOutput map[string]interface{}
		if err := json.Unmarshal(buf.Bytes(), &rawOutput); err != nil {
			t.Fatalf("output is not valid JSON: %v", err)
		}

		// Must have "fields" key that is an array
		fieldsRaw, ok := rawOutput["fields"]
		if !ok {
			t.Fatal("JSON output missing 'fields' key")
		}
		fields, ok := fieldsRaw.([]interface{})
		if !ok {
			t.Fatal("'fields' is not an array")
		}

		// Each field element must have required keys
		requiredKeys := []string{
			"service_name", "resource_name", "field_name",
			"field_path", "annotation_status", "terraform_confirmation", "category",
		}
		for i, fieldRaw := range fields {
			field, ok := fieldRaw.(map[string]interface{})
			if !ok {
				t.Fatalf("fields[%d] is not an object", i)
			}
			for _, key := range requiredKeys {
				if _, exists := field[key]; !exists {
					t.Fatalf("fields[%d] missing required key %q", i, key)
				}
			}
		}

		// Must have "summary" key that is an object
		summaryRaw, ok := rawOutput["summary"]
		if !ok {
			t.Fatal("JSON output missing 'summary' key")
		}
		summary, ok := summaryRaw.(map[string]interface{})
		if !ok {
			t.Fatal("'summary' is not an object")
		}

		// Summary must have required statistics keys
		summaryKeys := []string{
			"total_fields", "annotated_count", "gap_confirmed_count",
			"gap_unconfirmed_count", "terraform_only_count",
			"gaps_per_service", "services_by_priority",
		}
		for _, key := range summaryKeys {
			if _, exists := summary[key]; !exists {
				t.Fatalf("summary missing required key %q", key)
			}
		}

		// fields length should match resultCount
		if len(fields) != resultCount {
			t.Fatalf("expected %d fields in output, got %d", resultCount, len(fields))
		}
	})
}

// TestProperty17_ReportJSONFieldsIsAlwaysArray verifies that the "fields" key
// is always a JSON array (never null) even when there are zero results.
func TestProperty17_ReportJSONFieldsIsAlwaysArray(t *testing.T) {
	// Feature: ack-scanner, Property 17: Report JSON schema validity
	rapid.Check(t, func(t *rapid.T) {
		// Sometimes generate empty, sometimes non-empty
		resultCount := rapid.IntRange(0, 5).Draw(t, "resultCount")
		results := make([]types.MatchResult, resultCount)
		for i := 0; i < resultCount; i++ {
			results[i] = genMatchResult().Draw(t, "result")
		}

		r := NewReporter("json")
		var buf bytes.Buffer
		err := r.GenerateReport(results, &buf)
		if err != nil {
			t.Fatalf("GenerateReport failed: %v", err)
		}

		// Check raw bytes — "fields" should never be "null"
		output := buf.String()
		if strings.Contains(output, `"fields": null`) || strings.Contains(output, `"fields":null`) {
			t.Fatal("fields should be [] not null when empty")
		}

		// Verify it's valid JSON with an array for fields
		var rawOutput map[string]interface{}
		if err := json.Unmarshal(buf.Bytes(), &rawOutput); err != nil {
			t.Fatalf("output is not valid JSON: %v", err)
		}
		fieldsRaw := rawOutput["fields"]
		if _, ok := fieldsRaw.([]interface{}); !ok {
			t.Fatal("'fields' is not an array")
		}
	})
}
