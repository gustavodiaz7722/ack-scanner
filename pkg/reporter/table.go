package reporter

import (
	"fmt"
	"io"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"

	"github.com/aws-controllers-k8s/ack-scanner/pkg/types"
)

// formatReportTable outputs the report as a column-aligned table with a summary section.
func formatReportTable(results []types.MatchResult, summary types.ReportSummary, w io.Writer) error {
	table := tablewriter.NewTable(w,
		tablewriter.WithHeaderAutoFormat(tw.Off),
	)
	table.Header([]string{"SERVICE", "RESOURCE", "FIELD", "ANNOTATION", "TF CONFIRMED", "CATEGORY"})

	for _, r := range results {
		table.Append([]string{
			r.ServiceName,
			r.ResourceName,
			r.FieldName,
			string(r.AnnotationStatus),
			string(r.TFConfirmation),
			string(r.Category),
		})
	}

	if err := table.Render(); err != nil {
		return err
	}

	// Print summary section
	fmt.Fprintf(w, "\nSummary:\n")
	fmt.Fprintf(w, "  Total Fields:         %d\n", summary.TotalFields)
	fmt.Fprintf(w, "  Annotated:            %d\n", summary.AnnotatedCount)
	fmt.Fprintf(w, "  Gap (confirmed):      %d\n", summary.GapConfirmedCount)
	fmt.Fprintf(w, "  Gap (unconfirmed):    %d\n", summary.GapUnconfirmedCount)
	fmt.Fprintf(w, "  Terraform Only:       %d\n", summary.TerraformOnlyCount)

	if len(summary.ServicesByPriority) > 0 {
		fmt.Fprintf(w, "\nServices by Priority (confirmed gaps):\n")
		for _, sp := range summary.ServicesByPriority {
			fmt.Fprintf(w, "  %s: %d\n", sp.ServiceName, sp.ConfirmedGapCount)
		}
	}

	return nil
}

// formatControllerListTable outputs the controller list as a column-aligned table
// with a header row containing SERVICE and REPOSITORY columns.
func formatControllerListTable(repos []types.ControllerRepo, w io.Writer) error {
	table := tablewriter.NewTable(w,
		tablewriter.WithHeaderAutoFormat(tw.Off),
	)
	table.Header([]string{"SERVICE", "REPOSITORY"})

	for _, repo := range repos {
		table.Append([]string{repo.ServiceName, repo.RepoName})
	}

	return table.Render()
}

// formatReportTableWithResources outputs the report as a table, including all
// controllers. For now, delegates to formatReportTable (table format doesn't
// easily show empty controllers inline).
func formatReportTableWithResources(results []types.MatchResult, summary types.ReportSummary, controllerResources map[string][]string, w io.Writer) error {
	// For table format, just output results as before. The table format
	// doesn't lend itself well to showing empty controllers.
	return formatReportTable(results, summary, w)
}
