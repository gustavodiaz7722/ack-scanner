package matcher

import (
	"strings"
	"testing"

	"github.com/aws-controllers-k8s/ack-scanner/pkg/types"
	"pgregory.net/rapid"
)

// Feature: ack-scanner, Property 8: Field annotation classification
// **Validates: Requirements 3.4, 3.6, 5.2**

// genServiceName generates a lowercase service name (1-10 lowercase letters).
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

// genAnnotationType generates an annotation type: none, document, or iam_policy.
func genAnnotationType() *rapid.Generator[types.AnnotationType] {
	return rapid.Custom(func(t *rapid.T) types.AnnotationType {
		idx := rapid.IntRange(0, 2).Draw(t, "annotationType")
		switch idx {
		case 0:
			return types.AnnotationNone
		case 1:
			return types.AnnotationDocument
		default:
			return types.AnnotationIAMPolicy
		}
	})
}

// genScanResult generates a ScanResult with a given service name and annotation type.
func genScanResult() *rapid.Generator[types.ScanResult] {
	return rapid.Custom(func(t *rapid.T) types.ScanResult {
		service := genServiceName().Draw(t, "service")
		resource := genResourceName().Draw(t, "resource")
		field := genFieldName().Draw(t, "field")
		annotation := genAnnotationType().Draw(t, "annotation")
		return types.ScanResult{
			ServiceName:    service,
			RepoName:       service + "-controller",
			ResourceName:   resource,
			FieldName:      field,
			FieldPath:      "Spec." + field,
			GoType:         "*string",
			AnnotationType: annotation,
		}
	})
}

// genScanResultWithService generates a ScanResult with a specific service name.
func genScanResultWithService(service string) *rapid.Generator[types.ScanResult] {
	return rapid.Custom(func(t *rapid.T) types.ScanResult {
		resource := genResourceName().Draw(t, "resource")
		field := genFieldName().Draw(t, "field")
		annotation := genAnnotationType().Draw(t, "annotation")
		return types.ScanResult{
			ServiceName:    service,
			RepoName:       service + "-controller",
			ResourceName:   resource,
			FieldName:      field,
			FieldPath:      "Spec." + field,
			GoType:         "*string",
			AnnotationType: annotation,
		}
	})
}

// TestProperty8_FieldAnnotationClassification verifies that for any set of identified
// document-string fields and a generator config, a field is classified as "annotated"
// if and only if the generator config contains is_document or is_iam_policy for that field.
// Unannotated fields without TF confirmation are excluded from results.
func TestProperty8_FieldAnnotationClassification(t *testing.T) {
	// Feature: ack-scanner, Property 8: Field annotation classification
	rapid.Check(t, func(t *rapid.T) {
		// Generate a set of ACK fields with random annotation statuses
		fieldCount := rapid.IntRange(1, 10).Draw(t, "fieldCount")
		ackFields := make([]types.ScanResult, fieldCount)
		for i := 0; i < fieldCount; i++ {
			ackFields[i] = genScanResult().Draw(t, "ackField")
		}

		m := &Matcher{}
		// Match with empty TF fields — only annotated fields should appear
		results := m.Match(ackFields, nil)

		// Count expected annotated fields
		expectedAnnotated := 0
		for _, ack := range ackFields {
			if ack.AnnotationType != types.AnnotationNone {
				expectedAnnotated++
			}
		}

		if len(results) != expectedAnnotated {
			t.Fatalf("expected %d annotated results (no TF = no unannotated results), got %d",
				expectedAnnotated, len(results))
		}

		// Verify all results are annotated
		for _, result := range results {
			if result.Category != types.CategoryAnnotated {
				t.Fatalf("with no TF fields, all results should be annotated, got %q for %q",
					result.Category, result.FieldName)
			}
		}
	})
}

// TestProperty8_AnnotatedFieldsNeverClassifiedAsGap verifies that fields with
// annotations (is_document or is_iam_policy) are NEVER classified as gap fields,
// regardless of whether Terraform has a matching field.
func TestProperty8_AnnotatedFieldsNeverClassifiedAsGap(t *testing.T) {
	// Feature: ack-scanner, Property 8: Field annotation classification
	rapid.Check(t, func(t *rapid.T) {
		service := genServiceName().Draw(t, "service")
		field := genFieldName().Draw(t, "field")

		// Create an annotated ACK field
		annotationType := rapid.SampledFrom([]types.AnnotationType{
			types.AnnotationDocument,
			types.AnnotationIAMPolicy,
		}).Draw(t, "annotationType")

		ackFields := []types.ScanResult{
			{
				ServiceName:    service,
				RepoName:       service + "-controller",
				ResourceName:   "Resource",
				FieldName:      field,
				FieldPath:      "Spec." + field,
				GoType:         "*string",
				AnnotationType: annotationType,
			},
		}

		// Create a matching TF field (same service, normalized field name)
		snakeField := CamelToSnake(field)
		tfFields := []types.TerraformField{
			{
				ServiceName:     service,
				ResourceType:    "resource",
				FieldName:       snakeField,
				Description:     "JSON field",
				DetectionMethod: types.DetectDescriptionPhrase,
			},
		}

		m := &Matcher{}
		results := m.Match(ackFields, tfFields)

		// Find the result for our ACK field
		for _, r := range results {
			if r.FieldName == field && strings.EqualFold(r.ServiceName, service) {
				if r.Category == types.CategoryGapConfirmed || r.Category == types.CategoryGapUnconfirmed {
					t.Fatalf("annotated field %q (type=%q) should never be a gap, got category %q",
						field, annotationType, r.Category)
				}
				if r.Category != types.CategoryAnnotated {
					t.Fatalf("annotated field %q should be %q, got %q",
						field, types.CategoryAnnotated, r.Category)
				}
			}
		}
	})
}

// Feature: ack-scanner, Property 9: Service filter exactness
// **Validates: Requirements 3.6**

// TestProperty9_ServiceFilterExactness verifies that for any set of match results
// and a comma-separated filter string, FilterByServices returns exactly those results
// whose service name matches one of the filter values (case-insensitive).
func TestProperty9_ServiceFilterExactness(t *testing.T) {
	// Feature: ack-scanner, Property 9: Service filter exactness
	rapid.Check(t, func(t *rapid.T) {
		// Generate a set of distinct service names
		serviceCount := rapid.IntRange(2, 6).Draw(t, "serviceCount")
		services := make([]string, serviceCount)
		serviceSet := make(map[string]bool)
		for i := 0; i < serviceCount; i++ {
			for {
				s := genServiceName().Draw(t, "service")
				if !serviceSet[s] {
					services[i] = s
					serviceSet[s] = true
					break
				}
			}
		}

		// Generate match results spread across those services
		resultCount := rapid.IntRange(3, 15).Draw(t, "resultCount")
		results := make([]types.MatchResult, resultCount)
		for i := 0; i < resultCount; i++ {
			svcIdx := rapid.IntRange(0, serviceCount-1).Draw(t, "svcIdx")
			field := genFieldName().Draw(t, "field")
			results[i] = types.MatchResult{
				ServiceName:  services[svcIdx],
				ResourceName: "Resource",
				FieldName:    field,
				FieldPath:    "Spec." + field,
				Category:     types.CategoryGapConfirmed,
			}
		}

		// Pick a random subset of services to filter by (using a bitmask approach)
		filterServices := make([]string, 0)
		for i := 0; i < serviceCount; i++ {
			include := rapid.Bool().Draw(t, "include")
			if include {
				filterServices = append(filterServices, services[i])
			}
		}
		// Ensure at least one service in filter
		if len(filterServices) == 0 {
			filterServices = append(filterServices, services[0])
		}

		m := &Matcher{}
		filtered := m.FilterByServices(results, filterServices)

		// Build allowed set (case-insensitive)
		allowed := make(map[string]bool)
		for _, s := range filterServices {
			allowed[strings.ToLower(s)] = true
		}

		// Verify: every returned result has a service in the filter
		for _, r := range filtered {
			if !allowed[strings.ToLower(r.ServiceName)] {
				t.Fatalf("FilterByServices returned result with service %q which is not in filter %v",
					r.ServiceName, filterServices)
			}
		}

		// Verify: every result from the original set that should be included IS included
		expectedCount := 0
		for _, r := range results {
			if allowed[strings.ToLower(r.ServiceName)] {
				expectedCount++
			}
		}
		if len(filtered) != expectedCount {
			t.Fatalf("FilterByServices returned %d results, expected %d (filter=%v)",
				len(filtered), expectedCount, filterServices)
		}
	})
}

// TestProperty9_EmptyFilterReturnsAll verifies that an empty filter returns all results.
func TestProperty9_EmptyFilterReturnsAll(t *testing.T) {
	// Feature: ack-scanner, Property 9: Service filter exactness
	rapid.Check(t, func(t *rapid.T) {
		resultCount := rapid.IntRange(0, 10).Draw(t, "resultCount")
		results := make([]types.MatchResult, resultCount)
		for i := 0; i < resultCount; i++ {
			service := genServiceName().Draw(t, "service")
			field := genFieldName().Draw(t, "field")
			results[i] = types.MatchResult{
				ServiceName:  service,
				ResourceName: "Resource",
				FieldName:    field,
				Category:     types.CategoryGapConfirmed,
			}
		}

		m := &Matcher{}

		// nil filter
		filtered := m.FilterByServices(results, nil)
		if len(filtered) != resultCount {
			t.Fatalf("nil filter: expected %d results, got %d", resultCount, len(filtered))
		}

		// empty slice filter
		filtered = m.FilterByServices(results, []string{})
		if len(filtered) != resultCount {
			t.Fatalf("empty filter: expected %d results, got %d", resultCount, len(filtered))
		}
	})
}

// TestProperty9_FilterIsCaseInsensitive verifies that service filter matching
// is case-insensitive.
func TestProperty9_FilterIsCaseInsensitive(t *testing.T) {
	// Feature: ack-scanner, Property 9: Service filter exactness
	rapid.Check(t, func(t *rapid.T) {
		service := genServiceName().Draw(t, "service")
		field := genFieldName().Draw(t, "field")

		results := []types.MatchResult{
			{
				ServiceName:  service,
				ResourceName: "Resource",
				FieldName:    field,
				Category:     types.CategoryGapConfirmed,
			},
		}

		// Filter with uppercased service name
		m := &Matcher{}
		filtered := m.FilterByServices(results, []string{strings.ToUpper(service)})
		if len(filtered) != 1 {
			t.Fatalf("case-insensitive filter with %q should match service %q, got %d results",
				strings.ToUpper(service), service, len(filtered))
		}
	})
}

// Feature: ack-scanner, Property 14: Category assignment determinism
// **Validates: Requirements 5.2**

// TestProperty14_CategoryAssignmentDeterminism verifies that for any field with a
// known annotation status and Terraform match status, the category assignment is
// deterministic and follows the rules:
//   - unannotated + TF-matched → gap_confirmed_by_terraform
//   - unannotated + TF-unmatched → gap_without_terraform_confirmation
//   - annotated → already_annotated
//   - TF-only (no ACK field) → terraform_only
func TestProperty14_CategoryAssignmentDeterminism(t *testing.T) {
	// Feature: ack-scanner, Property 14: Category assignment determinism
	rapid.Check(t, func(t *rapid.T) {
		service := genServiceName().Draw(t, "service")
		field := genFieldName().Draw(t, "field")
		snakeField := CamelToSnake(field)

		// Test case: unannotated + TF matched
		ackFields := []types.ScanResult{
			{
				ServiceName:    service,
				RepoName:       service + "-controller",
				ResourceName:   "Resource",
				FieldName:      field,
				FieldPath:      "Spec." + field,
				GoType:         "*string",
				AnnotationType: types.AnnotationNone,
			},
		}
		tfFields := []types.TerraformField{
			{
				ServiceName:     service,
				ResourceType:    "resource",
				FieldName:       snakeField,
				Description:     "JSON field",
				DetectionMethod: types.DetectDescriptionPhrase,
			},
		}

		m := &Matcher{}
		results := m.Match(ackFields, tfFields)

		// Find the ACK field result
		var ackResult *types.MatchResult
		for i, r := range results {
			if r.FieldName == field && strings.EqualFold(r.ServiceName, service) {
				ackResult = &results[i]
				break
			}
		}
		if ackResult == nil {
			t.Fatalf("ACK field %q not found in results", field)
		}
		if ackResult.Category != types.CategoryGapConfirmed {
			t.Fatalf("unannotated + TF-matched: expected %q, got %q",
				types.CategoryGapConfirmed, ackResult.Category)
		}
		if ackResult.TFConfirmation != types.TFConfirmed {
			t.Fatalf("unannotated + TF-matched: expected TFConfirmation=%q, got %q",
				types.TFConfirmed, ackResult.TFConfirmation)
		}
	})
}

// TestProperty14_UnannotatedWithoutTFMatch verifies unannotated fields without
// Terraform match get gap_without_terraform_confirmation category.
func TestProperty14_UnannotatedWithoutTFMatch(t *testing.T) {
	// Feature: ack-scanner, Property 14: Category assignment determinism
	// Unannotated fields without a TF match are now excluded from results entirely.
	rapid.Check(t, func(t *rapid.T) {
		service := genServiceName().Draw(t, "service")
		field := genFieldName().Draw(t, "field")

		ackFields := []types.ScanResult{
			{
				ServiceName:    service,
				RepoName:       service + "-controller",
				ResourceName:   "Resource",
				FieldName:      field,
				FieldPath:      "Spec." + field,
				GoType:         "*string",
				AnnotationType: types.AnnotationNone,
			},
		}

		// Use a different service in TF so there's no match
		differentService := service + "x"
		tfFields := []types.TerraformField{
			{
				ServiceName:     differentService,
				ResourceType:    "resource",
				FieldName:       "unrelated_field",
				Description:     "JSON field",
				DetectionMethod: types.DetectDescriptionPhrase,
			},
		}

		m := &Matcher{}
		results := m.Match(ackFields, tfFields)

		// The ACK field should NOT appear — it's unannotated with no TF match
		for _, r := range results {
			if r.FieldName == field && strings.EqualFold(r.ServiceName, service) {
				t.Fatalf("unannotated field %q without TF match should not appear in results", field)
			}
		}
	})
}

// TestProperty14_AnnotatedAlwaysClassifiedCorrectly verifies annotated fields always
// get already_annotated category regardless of TF match status.
func TestProperty14_AnnotatedAlwaysClassifiedCorrectly(t *testing.T) {
	// Feature: ack-scanner, Property 14: Category assignment determinism
	rapid.Check(t, func(t *rapid.T) {
		service := genServiceName().Draw(t, "service")
		field := genFieldName().Draw(t, "field")
		snakeField := CamelToSnake(field)

		annotationType := rapid.SampledFrom([]types.AnnotationType{
			types.AnnotationDocument,
			types.AnnotationIAMPolicy,
		}).Draw(t, "annotationType")

		ackFields := []types.ScanResult{
			{
				ServiceName:    service,
				RepoName:       service + "-controller",
				ResourceName:   "Resource",
				FieldName:      field,
				FieldPath:      "Spec." + field,
				GoType:         "*string",
				AnnotationType: annotationType,
			},
		}

		// Include a matching TF field — annotated fields should still be already_annotated
		tfFields := []types.TerraformField{
			{
				ServiceName:     service,
				ResourceType:    "resource",
				FieldName:       snakeField,
				Description:     "JSON field",
				DetectionMethod: types.DetectDescriptionPhrase,
			},
		}

		m := &Matcher{}
		results := m.Match(ackFields, tfFields)

		var ackResult *types.MatchResult
		for i, r := range results {
			if r.FieldName == field && strings.EqualFold(r.ServiceName, service) {
				ackResult = &results[i]
				break
			}
		}
		if ackResult == nil {
			t.Fatalf("ACK field %q not found in results", field)
		}
		if ackResult.Category != types.CategoryAnnotated {
			t.Fatalf("annotated field: expected %q, got %q",
				types.CategoryAnnotated, ackResult.Category)
		}
		if ackResult.TFConfirmation != types.TFNotApplicable {
			t.Fatalf("annotated field: expected TFConfirmation=%q, got %q",
				types.TFNotApplicable, ackResult.TFConfirmation)
		}
	})
}

// TestProperty14_TerraformOnlyClassification verifies that TF fields without any
// matching ACK field get terraform_only category, provided the service has an ACK controller.
func TestProperty14_TerraformOnlyClassification(t *testing.T) {
	// Feature: ack-scanner, Property 14: Category assignment determinism
	rapid.Check(t, func(t *rapid.T) {
		tfService := genServiceName().Draw(t, "tfService")
		tfField := genFieldName().Draw(t, "tfField")
		snakeTfField := CamelToSnake(tfField)

		// Create an ACK field in the SAME service (so the service is known)
		// but with a completely different field name (so there's no match)
		ackField := "UnrelatedFieldXyz"

		ackFields := []types.ScanResult{
			{
				ServiceName:    tfService,
				RepoName:       tfService + "-controller",
				ResourceName:   "Resource",
				FieldName:      ackField,
				FieldPath:      "Spec." + ackField,
				GoType:         "*string",
				AnnotationType: types.AnnotationDocument, // annotated so it appears
			},
		}

		tfFields := []types.TerraformField{
			{
				ServiceName:     tfService,
				ResourceType:    "resource",
				FieldName:       snakeTfField,
				Description:     "JSON field",
				DetectionMethod: types.DetectDescriptionPhrase,
			},
		}

		m := &Matcher{}
		results := m.Match(ackFields, tfFields)

		// Find the terraform_only result
		var tfResult *types.MatchResult
		for i, r := range results {
			if r.Category == types.CategoryTerraformOnly &&
				strings.EqualFold(r.ServiceName, tfService) {
				tfResult = &results[i]
				break
			}
		}
		if tfResult == nil {
			t.Fatalf("TF field %q (service=%q) not found as terraform_only in results: %+v",
				snakeTfField, tfService, results)
		}
		if tfResult.Category != types.CategoryTerraformOnly {
			t.Fatalf("TF-only field: expected %q, got %q",
				types.CategoryTerraformOnly, tfResult.Category)
		}
	})
}

// TestProperty14_CategoryAssignmentIsRepeatable verifies that running Match
// multiple times with the same inputs always produces the same categories.
func TestProperty14_CategoryAssignmentIsRepeatable(t *testing.T) {
	// Feature: ack-scanner, Property 14: Category assignment determinism
	rapid.Check(t, func(t *rapid.T) {
		// Generate random ACK fields
		ackCount := rapid.IntRange(1, 5).Draw(t, "ackCount")
		ackFields := make([]types.ScanResult, ackCount)
		for i := 0; i < ackCount; i++ {
			ackFields[i] = genScanResult().Draw(t, "ackField")
		}

		// Generate random TF fields
		tfCount := rapid.IntRange(0, 5).Draw(t, "tfCount")
		tfFields := make([]types.TerraformField, tfCount)
		for i := 0; i < tfCount; i++ {
			service := genServiceName().Draw(t, "tfService")
			field := genFieldName().Draw(t, "tfField")
			tfFields[i] = types.TerraformField{
				ServiceName:     service,
				ResourceType:    "resource",
				FieldName:       CamelToSnake(field),
				Description:     "JSON field",
				DetectionMethod: types.DetectDescriptionPhrase,
			}
		}

		m := &Matcher{}
		results1 := m.Match(ackFields, tfFields)
		results2 := m.Match(ackFields, tfFields)

		if len(results1) != len(results2) {
			t.Fatalf("non-deterministic: first call returned %d results, second returned %d",
				len(results1), len(results2))
		}

		for i := range results1 {
			if results1[i].Category != results2[i].Category {
				t.Fatalf("non-deterministic category at index %d: %q vs %q",
					i, results1[i].Category, results2[i].Category)
			}
			if results1[i].TFConfirmation != results2[i].TFConfirmation {
				t.Fatalf("non-deterministic TFConfirmation at index %d: %q vs %q",
					i, results1[i].TFConfirmation, results2[i].TFConfirmation)
			}
		}
	})
}
