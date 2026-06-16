// Package reporter generates formatted output from scan results.
package reporter

import (
	"fmt"
	"io"
	"sort"

	"github.com/aws-controllers-k8s/ack-scanner/pkg/types"
)

// Reporter generates formatted output from scan results.
type Reporter struct {
	format string // "json", "table", "markdown"
}

// NewReporter creates a new Reporter with the specified output format.
func NewReporter(format string) *Reporter {
	return &Reporter{format: format}
}

// GenerateReport creates a gap analysis report from matched results.
// It sorts the results, computes a summary, and dispatches to the
// appropriate formatter based on the configured format.
func (r *Reporter) GenerateReport(results []types.MatchResult, w io.Writer) error {
	sorted := sortResults(results)
	summary := GenerateSummary(sorted)

	switch r.format {
	case "json":
		return formatReportJSON(sorted, summary, w)
	case "table":
		return formatReportTable(sorted, summary, w)
	case "markdown":
		return formatReportMarkdown(sorted, summary, w)
	default:
		return fmt.Errorf("unsupported output format: %s", r.format)
	}
}

// FormatControllerList outputs the controller list in the specified format.
// Controllers are sorted alphabetically by service name before formatting.
func (r *Reporter) FormatControllerList(repos []types.ControllerRepo, w io.Writer) error {
	sorted := make([]types.ControllerRepo, len(repos))
	copy(sorted, repos)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].ServiceName < sorted[j].ServiceName
	})

	switch r.format {
	case "json":
		return formatControllerListJSON(sorted, w)
	case "table":
		return formatControllerListTable(sorted, w)
	case "markdown":
		return formatControllerListMarkdown(sorted, w)
	default:
		return fmt.Errorf("unsupported output format: %s", r.format)
	}
}

// GenerateSummary computes aggregate statistics from match results.
// It counts fields per category, gaps per service (sum of gap_confirmed
// and gap_unconfirmed), and ranks services by confirmed gap count.
func GenerateSummary(results []types.MatchResult) types.ReportSummary {
	summary := types.ReportSummary{
		TotalFields:    len(results),
		GapsPerService: make(map[string]int),
	}

	// Count per-service confirmed gaps for priority ranking
	confirmedPerService := make(map[string]int)

	for _, r := range results {
		switch r.Category {
		case types.CategoryAnnotated:
			summary.AnnotatedCount++
		case types.CategoryGapConfirmed:
			summary.GapConfirmedCount++
			summary.GapsPerService[r.ServiceName]++
			confirmedPerService[r.ServiceName]++
		case types.CategoryGapUnconfirmed:
			summary.GapUnconfirmedCount++
			summary.GapsPerService[r.ServiceName]++
		case types.CategoryTerraformOnly:
			summary.TerraformOnlyCount++
		}
	}

	// Build services by priority sorted descending by confirmed gap count
	for svc, count := range confirmedPerService {
		summary.ServicesByPriority = append(summary.ServicesByPriority, types.ServicePriority{
			ServiceName:       svc,
			ConfirmedGapCount: count,
		})
	}
	sort.Slice(summary.ServicesByPriority, func(i, j int) bool {
		if summary.ServicesByPriority[i].ConfirmedGapCount != summary.ServicesByPriority[j].ConfirmedGapCount {
			return summary.ServicesByPriority[i].ConfirmedGapCount > summary.ServicesByPriority[j].ConfirmedGapCount
		}
		return summary.ServicesByPriority[i].ServiceName < summary.ServicesByPriority[j].ServiceName
	})

	return summary
}

// sortResults sorts match results by category priority:
// gap_confirmed first, then gap_unconfirmed, then annotated, then terraform_only.
// Within the same category, results are sorted by service name then field name.
func sortResults(results []types.MatchResult) []types.MatchResult {
	sorted := make([]types.MatchResult, len(results))
	copy(sorted, results)

	sort.SliceStable(sorted, func(i, j int) bool {
		pi := categoryPriority(sorted[i].Category)
		pj := categoryPriority(sorted[j].Category)
		if pi != pj {
			return pi < pj
		}
		if sorted[i].ServiceName != sorted[j].ServiceName {
			return sorted[i].ServiceName < sorted[j].ServiceName
		}
		return sorted[i].FieldName < sorted[j].FieldName
	})

	return sorted
}

// categoryPriority returns the sort priority for a category.
// Lower values sort first.
func categoryPriority(c types.Category) int {
	switch c {
	case types.CategoryGapConfirmed:
		return 0
	case types.CategoryGapUnconfirmed:
		return 1
	case types.CategoryAnnotated:
		return 2
	case types.CategoryTerraformOnly:
		return 3
	default:
		return 4
	}
}
