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
func (m *Matcher) Match(ackFields []types.ScanResult, tfFields []types.TerraformField) []types.MatchResult {
	// Step 1: Build a lookup map from Terraform fields.
	// Key = (lowercase service name, normalized field name)
	tfLookup := make(map[lookupKey]types.TerraformField, len(tfFields))
	for _, tf := range tfFields {
		key := lookupKey{
			service: strings.ToLower(tf.ServiceName),
			field:   NormalizeFieldName(tf.FieldName),
		}
		tfLookup[key] = tf
	}

	// Track which Terraform fields were matched to an ACK field
	matchedTFKeys := make(map[lookupKey]bool)

	var results []types.MatchResult

	// Step 2: Process each ACK scan result
	for _, ack := range ackFields {
		normalizedField := NormalizeFieldName(ack.FieldName)
		key := lookupKey{
			service: strings.ToLower(ack.ServiceName),
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
				category = types.CategoryGapConfirmed
				tfConfirmation = types.TFConfirmed
				matchedTFKeys[key] = true
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
			// We need to verify it wasn't matched by an annotated ACK field either
			ackMatched := false
			for _, ack := range ackFields {
				ackKey := lookupKey{
					service: strings.ToLower(ack.ServiceName),
					field:   NormalizeFieldName(ack.FieldName),
				}
				if ackKey == key {
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
