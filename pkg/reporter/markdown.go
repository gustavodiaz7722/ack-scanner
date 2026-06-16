package reporter

import (
	"fmt"
	"io"
	"sort"

	"github.com/aws-controllers-k8s/ack-scanner/pkg/types"
)

// formatReportMarkdown outputs the report as markdown with tables grouped by service.
// Gap fields (confirmed or unconfirmed) are highlighted with a ⚠️ prefix on the field name.
// Within each service group, confirmed gaps sort before unconfirmed gaps.
func formatReportMarkdown(results []types.MatchResult, summary types.ReportSummary, w io.Writer) error {
	// Title
	if _, err := fmt.Fprintf(w, "# ACK Scanner Gap Analysis Report\n\n"); err != nil {
		return err
	}

	// Summary section
	if _, err := fmt.Fprintf(w, "## Summary\n\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "| Metric | Count |\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "| --- | --- |\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "| Total Fields | %d |\n", summary.TotalFields); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "| Already Annotated | %d |\n", summary.AnnotatedCount); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "| Gaps (Terraform Confirmed) | %d |\n", summary.GapConfirmedCount); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "| Gaps (Unconfirmed) | %d |\n", summary.GapUnconfirmedCount); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "| Terraform Only | %d |\n", summary.TerraformOnlyCount); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "\n"); err != nil {
		return err
	}

	// Priority Services section
	if _, err := fmt.Fprintf(w, "## Priority Services\n\n"); err != nil {
		return err
	}
	if len(summary.ServicesByPriority) == 0 {
		if _, err := fmt.Fprintf(w, "No priority services identified.\n\n"); err != nil {
			return err
		}
	} else {
		if _, err := fmt.Fprintf(w, "| Service | Confirmed Gaps |\n"); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, "| --- | --- |\n"); err != nil {
			return err
		}
		for _, sp := range summary.ServicesByPriority {
			if _, err := fmt.Fprintf(w, "| %s | %d |\n", sp.ServiceName, sp.ConfirmedGapCount); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintf(w, "\n"); err != nil {
			return err
		}
	}

	// Group results by service
	serviceOrder, grouped := groupByService(results)

	for _, svc := range serviceOrder {
		svcResults := grouped[svc]
		if _, err := fmt.Fprintf(w, "## %s\n\n", svc); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, "| Resource | Field | Annotation | TF Confirmed | Category |\n"); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, "| --- | --- | --- | --- | --- |\n"); err != nil {
			return err
		}

		// Sort within service: confirmed gaps first, then unconfirmed, then others
		sort.SliceStable(svcResults, func(i, j int) bool {
			pi := categoryPriority(svcResults[i].Category)
			pj := categoryPriority(svcResults[j].Category)
			if pi != pj {
				return pi < pj
			}
			return svcResults[i].FieldName < svcResults[j].FieldName
		})

		for _, r := range svcResults {
			fieldDisplay := r.FieldName
			if r.Category == types.CategoryGapConfirmed || r.Category == types.CategoryGapUnconfirmed {
				fieldDisplay = "⚠️ " + r.FieldName
			}

			annotation := formatAnnotation(r.AnnotationStatus)
			tfConfirmed := formatTFConfirmation(r.TFConfirmation)
			category := formatCategory(r.Category)

			if _, err := fmt.Fprintf(w, "| %s | %s | %s | %s | %s |\n",
				r.ResourceName, fieldDisplay, annotation, tfConfirmed, category); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintf(w, "\n"); err != nil {
			return err
		}
	}

	return nil
}

// formatControllerListMarkdown outputs the controller list as a markdown table.
func formatControllerListMarkdown(repos []types.ControllerRepo, w io.Writer) error {
	if _, err := fmt.Fprintf(w, "| Service | Repository |\n"); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "| --- | --- |\n"); err != nil {
		return err
	}
	for _, repo := range repos {
		if _, err := fmt.Fprintf(w, "| %s | %s |\n", repo.ServiceName, repo.RepoName); err != nil {
			return err
		}
	}
	return nil
}

// groupByService groups results by service name and returns the services
// in sorted order along with the grouped map.
func groupByService(results []types.MatchResult) ([]string, map[string][]types.MatchResult) {
	grouped := make(map[string][]types.MatchResult)
	for _, r := range results {
		grouped[r.ServiceName] = append(grouped[r.ServiceName], r)
	}

	services := make([]string, 0, len(grouped))
	for svc := range grouped {
		services = append(services, svc)
	}
	sort.Strings(services)

	return services, grouped
}

// formatAnnotation returns a human-readable annotation status string.
func formatAnnotation(a types.AnnotationType) string {
	switch a {
	case types.AnnotationDocument:
		return "is_document"
	case types.AnnotationIAMPolicy:
		return "is_iam_policy"
	case types.AnnotationNone:
		return "none"
	default:
		return string(a)
	}
}

// formatTFConfirmation returns a human-readable Terraform confirmation status.
func formatTFConfirmation(tc types.TerraformConfirmation) string {
	switch tc {
	case types.TFConfirmed:
		return "Yes"
	case types.TFUnconfirmed:
		return "No"
	case types.TFNotApplicable:
		return "N/A"
	default:
		return string(tc)
	}
}

// formatCategory returns a human-readable category string.
func formatCategory(c types.Category) string {
	switch c {
	case types.CategoryGapConfirmed:
		return "Gap (Confirmed)"
	case types.CategoryGapUnconfirmed:
		return "Gap (Unconfirmed)"
	case types.CategoryAnnotated:
		return "Annotated"
	case types.CategoryTerraformOnly:
		return "Terraform Only"
	default:
		return string(c)
	}
}
