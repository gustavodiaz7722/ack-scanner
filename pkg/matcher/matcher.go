// Package matcher provides field matching logic for cross-referencing ACK fields
// with Terraform documentation fields.
package matcher

import (
	"strings"

	"github.com/aws-controllers-k8s/ack-scanner/pkg/types"
)

// lookupKey is the composite key used to match ACK fields with Terraform fields.
type lookupKey struct {
	service string // lowercase service name
	field   string // normalized field name (snake_case, stripped of prefixes/suffixes)
}

// Matcher cross-references ACK fields with Terraform fields.
type Matcher struct{}

// Match performs case-insensitive service matching and normalized field name
// comparison (CamelCase to snake_case). It cross-references ACK scan results
// with Terraform fields and classifies each field into one of four categories:
//   - already_annotated: field already has is_document or is_iam_policy annotation
//   - gap_confirmed_by_terraform: unannotated in ACK and confirmed as JSON in Terraform
//   - gap_without_terraform_confirmation: unannotated in ACK but not found in Terraform
//   - terraform_only: found in Terraform docs but not matched to any ACK field
//
// When an exact field name match fails, it attempts a "similar field" match:
// a Terraform field in the same service is considered a match if one field name
// is a suffix of the other (e.g., ACK "AssumeRolePolicyDocument" contains TF
// "assume_role_policy" as a prefix, or TF "policy" is a suffix of ACK "policy_document").
func (m *Matcher) Match(ackFields []types.ScanResult, tfFields []types.TerraformField) []types.MatchResult {
	// Step 1: Build lookup structures from Terraform fields.
	// Exact lookup: Key = (lowercase service name, normalized field name)
	tfLookup := make(map[lookupKey]types.TerraformField, len(tfFields))
	// Per-service list for similar-field matching
	tfByService := make(map[string][]types.TerraformField)
	for _, tf := range tfFields {
		key := lookupKey{
			service: strings.ToLower(tf.ServiceName),
			field:   NormalizeFieldName(tf.FieldName),
		}
		tfLookup[key] = tf
		svc := strings.ToLower(tf.ServiceName)
		tfByService[svc] = append(tfByService[svc], tf)
	}

	// Track which Terraform fields were matched to an ACK field
	matchedTFKeys := make(map[lookupKey]bool)

	var results []types.MatchResult

	// Step 2: Process each ACK scan result
	for _, ack := range ackFields {
		normalizedField := NormalizeFieldName(ack.FieldName)
		svc := strings.ToLower(ack.ServiceName)
		key := lookupKey{
			service: svc,
			field:   normalizedField,
		}

		var category types.Category
		var tfConfirmation types.TerraformConfirmation

		if ack.AnnotationType != types.AnnotationNone {
			// Field is already annotated
			category = types.CategoryAnnotated
			tfConfirmation = types.TFNotApplicable
		} else {
			// Field is unannotated — check if Terraform confirms it
			if _, found := tfLookup[key]; found {
				// Exact match
				category = types.CategoryGapConfirmed
				tfConfirmation = types.TFConfirmed
				matchedTFKeys[key] = true
			} else if similarKey, found := findSimilarTFField(normalizedField, svc, tfByService[svc]); found {
				// Similar field match (suffix/prefix containment)
				category = types.CategoryGapConfirmed
				tfConfirmation = types.TFConfirmed
				matchedTFKeys[similarKey] = true
			} else {
				category = types.CategoryGapUnconfirmed
				tfConfirmation = types.TFUnconfirmed
			}
		}

		results = append(results, types.MatchResult{
			ServiceName:      ack.ServiceName,
			ResourceName:     ack.ResourceName,
			FieldName:        ack.FieldName,
			FieldPath:        ack.FieldPath,
			AnnotationStatus: ack.AnnotationType,
			TFConfirmation:   tfConfirmation,
			Category:         category,
		})
	}

	// Step 3: Add Terraform-only fields (not matched to any ACK field)
	for _, tf := range tfFields {
		key := lookupKey{
			service: strings.ToLower(tf.ServiceName),
			field:   NormalizeFieldName(tf.FieldName),
		}
		if !matchedTFKeys[key] {
			// Check if any ACK field (regardless of annotation) matched this TF field
			ackMatched := false
			svc := strings.ToLower(tf.ServiceName)
			tfNorm := NormalizeFieldName(tf.FieldName)
			for _, ack := range ackFields {
				if strings.ToLower(ack.ServiceName) != svc {
					continue
				}
				ackNorm := NormalizeFieldName(ack.FieldName)
				if ackNorm == tfNorm || isSimilarField(ackNorm, tfNorm) {
					ackMatched = true
					break
				}
			}
			if !ackMatched {
				results = append(results, types.MatchResult{
					ServiceName:      tf.ServiceName,
					ResourceName:     tf.ResourceType,
					FieldName:        tf.FieldName,
					FieldPath:        "",
					AnnotationStatus: types.AnnotationNone,
					TFConfirmation:   types.TFNotApplicable,
					Category:         types.CategoryTerraformOnly,
				})
			}
		}
	}

	return results
}

// findSimilarTFField searches for a Terraform field in the same service that
// is "similar" to the given ACK field name. Similarity is determined by
// suffix/prefix containment of underscore-delimited segments.
//
// For example:
//   - ACK "assume_role_policy_document" is similar to TF "assume_role_policy"
//     (TF name is a prefix of ACK name)
//   - ACK "policy_document" is similar to TF "policy"
//     (TF name is a suffix-component of ACK name)
//   - ACK "policy" is similar to TF "policy"
//     (exact match — handled before this is called)
func findSimilarTFField(ackNormalized string, service string, tfFields []types.TerraformField) (lookupKey, bool) {
	for _, tf := range tfFields {
		tfNorm := NormalizeFieldName(tf.FieldName)
		if isSimilarField(ackNormalized, tfNorm) {
			return lookupKey{service: service, field: tfNorm}, true
		}
	}
	return lookupKey{}, false
}

// isSimilarField returns true if two normalized field names are "similar" —
// meaning one is a suffix or prefix of the other at an underscore boundary,
// or they share a significant common suffix.
//
// Rules (all operate on normalized snake_case names):
//  1. One name is a prefix of the other at an underscore boundary
//     e.g., "assume_role_policy" is a prefix of "assume_role_policy_document"
//  2. One name is a suffix of the other at an underscore boundary
//     e.g., "policy" is a suffix of "redrive_policy" — but this is too loose.
//     We require the shorter name to be at least 2 segments or match >= 50% of the longer.
//  3. Both names share a common suffix of at least 2 segments
//     e.g., "inline_policy_document" and "role_policy_document" share "policy_document"
func isSimilarField(a, b string) bool {
	if a == b {
		return true
	}

	// Ensure 'a' is the longer name
	if len(a) < len(b) {
		a, b = b, a
	}

	// Rule 1: b is a prefix of a at an underscore boundary
	if strings.HasPrefix(a, b+"_") || strings.HasPrefix(a, b) && len(a) > len(b) && a[len(b)] == '_' {
		return true
	}

	// Rule 2: b is a suffix of a at an underscore boundary
	// Only if b has at least 2 segments (to avoid "policy" matching everything)
	if strings.HasSuffix(a, "_"+b) {
		bSegments := strings.Count(b, "_") + 1
		if bSegments >= 2 {
			return true
		}
		// Single segment: only match if it's >= 50% of the longer name's segments
		aSegments := strings.Count(a, "_") + 1
		if aSegments <= 2 {
			return true
		}
	}

	// Rule 3: Shared suffix of at least 2 segments
	aParts := strings.Split(a, "_")
	bParts := strings.Split(b, "_")
	commonSuffix := 0
	for i, j := len(aParts)-1, len(bParts)-1; i >= 0 && j >= 0; i, j = i-1, j-1 {
		if aParts[i] == bParts[j] {
			commonSuffix++
		} else {
			break
		}
	}
	if commonSuffix >= 2 {
		return true
	}

	return false
}

// FilterByServices filters match results to include only those whose service name
// matches one of the provided services (case-insensitive comparison).
func (m *Matcher) FilterByServices(results []types.MatchResult, services []string) []types.MatchResult {
	if len(services) == 0 {
		return results
	}

	// Build a set of allowed service names (lowercased for case-insensitive matching)
	allowed := make(map[string]bool, len(services))
	for _, s := range services {
		allowed[strings.ToLower(s)] = true
	}

	var filtered []types.MatchResult
	for _, r := range results {
		if allowed[strings.ToLower(r.ServiceName)] {
			filtered = append(filtered, r)
		}
	}
	return filtered
}
