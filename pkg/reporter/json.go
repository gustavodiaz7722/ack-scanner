package reporter

import (
	"encoding/json"
	"io"

	"github.com/aws-controllers-k8s/ack-scanner/pkg/types"
)

// reportJSON is the top-level JSON output structure for the gap analysis report.
type reportJSON struct {
	Fields  []types.MatchResult `json:"fields"`
	Summary types.ReportSummary `json:"summary"`
}

// formatReportJSON outputs the report as a structured JSON document with
// a "fields" array and a "summary" object. Uses indented formatting for
// human readability.
func formatReportJSON(results []types.MatchResult, summary types.ReportSummary, w io.Writer) error {
	report := reportJSON{
		Fields:  results,
		Summary: summary,
	}

	// Ensure empty slices serialize as [] not null
	if report.Fields == nil {
		report.Fields = []types.MatchResult{}
	}
	if report.Summary.ServicesByPriority == nil {
		report.Summary.ServicesByPriority = []types.ServicePriority{}
	}
	if report.Summary.GapsPerService == nil {
		report.Summary.GapsPerService = make(map[string]int)
	}

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte("\n"))
	return err
}

// controllerListEntry is the JSON representation of a single controller
// in the list output.
type controllerListEntry struct {
	RepositoryName string `json:"repository_name"`
	ServiceName    string `json:"service_name"`
}

// formatControllerListJSON outputs the controller list as a JSON array of
// objects containing repository_name and service_name fields.
func formatControllerListJSON(repos []types.ControllerRepo, w io.Writer) error {
	entries := make([]controllerListEntry, len(repos))
	for i, repo := range repos {
		entries[i] = controllerListEntry{
			RepositoryName: repo.RepoName,
			ServiceName:    repo.ServiceName,
		}
	}

	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte("\n"))
	return err
}
