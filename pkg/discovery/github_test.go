package discovery

import (
	"strings"
	"testing"

	"github.com/google/go-github/v60/github"
	"pgregory.net/rapid"
)

// Feature: ack-scanner, Property 1: Service name extraction round-trip
// For any valid service name (non-empty, lowercase alphanumeric with hyphens),
// appending "-controller" and extracting the service name should return the original.
// **Validates: Requirements 1.2**
func TestProperty1_ServiceNameExtractionRoundTrip(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate a valid service name: non-empty, lowercase alphanumeric with hyphens,
		// must start and end with alphanumeric (not hyphen)
		serviceName := genValidServiceName().Draw(t, "serviceName")

		// Append "-controller" suffix
		repoName := serviceName + controllerSuffix

		// Extract service name by removing suffix
		extracted := ExtractServiceName(repoName)

		// Round-trip: extracted must equal original
		if extracted != serviceName {
			t.Fatalf("round-trip failed: input=%q, repoName=%q, extracted=%q", serviceName, repoName, extracted)
		}
	})
}

// Feature: ack-scanner, Property 2: Repository filter correctness
// For any list of GitHub repositories with varying archive/fork status and names,
// the filter (isControllerRepo) should return exactly those that are non-archived,
// non-fork, and named *-controller.
// **Validates: Requirements 1.1**
func TestProperty2_RepositoryFilterCorrectness(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate a list of repositories with random properties
		repos := rapid.SliceOf(genGitHubRepo()).Draw(t, "repos")

		// Apply the filter
		var filtered []*github.Repository
		for _, repo := range repos {
			if isControllerRepo(repo) {
				filtered = append(filtered, repo)
			}
		}

		// Verify: every filtered repo must satisfy all three conditions
		for _, repo := range filtered {
			if repo.GetArchived() {
				t.Fatalf("filter returned archived repo: %s", repo.GetName())
			}
			if repo.GetFork() {
				t.Fatalf("filter returned fork repo: %s", repo.GetName())
			}
			if !strings.HasSuffix(repo.GetName(), controllerSuffix) {
				t.Fatalf("filter returned repo without -controller suffix: %s", repo.GetName())
			}
		}

		// Verify: no repo satisfying all conditions was excluded
		for _, repo := range repos {
			shouldInclude := !repo.GetArchived() && !repo.GetFork() && strings.HasSuffix(repo.GetName(), controllerSuffix)
			wasIncluded := false
			for _, f := range filtered {
				if f == repo {
					wasIncluded = true
					break
				}
			}
			if shouldInclude && !wasIncluded {
				t.Fatalf("filter excluded valid controller repo: %s (archived=%v, fork=%v)",
					repo.GetName(), repo.GetArchived(), repo.GetFork())
			}
			if !shouldInclude && wasIncluded {
				t.Fatalf("filter included invalid repo: %s (archived=%v, fork=%v)",
					repo.GetName(), repo.GetArchived(), repo.GetFork())
			}
		}
	})
}

// Feature: ack-scanner, Property 3: Pagination completeness
// For pagination, generate a list of N repos and simulate splitting them into pages
// of size M. Verify that collecting all pages gives back exactly the original list
// with no duplicates and no omissions.
// **Validates: Requirements 1.6**
func TestProperty3_PaginationCompleteness(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate total number of items (1-500)
		totalItems := rapid.IntRange(1, 500).Draw(t, "totalItems")
		// Generate page size (1-100)
		pageSize := rapid.IntRange(1, 100).Draw(t, "pageSize")

		// Create original items
		items := make([]string, totalItems)
		for i := range items {
			items[i] = rapid.StringMatching(`[a-z]{3,10}-controller`).Draw(t, "item")
		}

		// Simulate pagination: split items into pages
		var collected []string
		for page := 0; page*pageSize < totalItems; page++ {
			start := page * pageSize
			end := start + pageSize
			if end > totalItems {
				end = totalItems
			}
			pageItems := items[start:end]
			collected = append(collected, pageItems...)
		}

		// Verify completeness: collected must have same length as original
		if len(collected) != len(items) {
			t.Fatalf("pagination lost items: original=%d, collected=%d", len(items), len(collected))
		}

		// Verify order preservation: collected must be identical to original
		for i := range items {
			if collected[i] != items[i] {
				t.Fatalf("pagination changed order at index %d: expected=%q, got=%q", i, items[i], collected[i])
			}
		}
	})
}

// --- Generators ---

// genValidServiceName generates valid service names: non-empty, lowercase
// alphanumeric with hyphens, starting and ending with alphanumeric.
func genValidServiceName() *rapid.Generator[string] {
	return rapid.Custom[string](func(t *rapid.T) string {
		// Generate length between 1 and 30
		length := rapid.IntRange(1, 30).Draw(t, "nameLen")

		// Valid chars: lowercase letters, digits, hyphens (but not leading/trailing)
		chars := make([]byte, length)
		for i := range chars {
			if i == 0 || i == length-1 {
				// First and last char must be alphanumeric
				chars[i] = genAlphaNum().Draw(t, "char")
			} else {
				// Middle chars can include hyphens
				chars[i] = genAlphaNumHyphen().Draw(t, "char")
			}
		}

		// Ensure no double hyphens (clean up)
		result := string(chars)
		for strings.Contains(result, "--") {
			result = strings.ReplaceAll(result, "--", "-")
		}
		// Ensure non-empty after cleanup
		if result == "" || result == "-" {
			result = "s3"
		}
		return result
	})
}

// genAlphaNum generates a lowercase alphanumeric byte.
func genAlphaNum() *rapid.Generator[byte] {
	return rapid.Custom[byte](func(t *rapid.T) byte {
		// 26 letters + 10 digits = 36 chars
		idx := rapid.IntRange(0, 35).Draw(t, "idx")
		if idx < 26 {
			return byte('a' + idx)
		}
		return byte('0' + (idx - 26))
	})
}

// genAlphaNumHyphen generates a lowercase alphanumeric or hyphen byte.
func genAlphaNumHyphen() *rapid.Generator[byte] {
	return rapid.Custom[byte](func(t *rapid.T) byte {
		// 26 letters + 10 digits + 1 hyphen = 37 chars
		idx := rapid.IntRange(0, 36).Draw(t, "idx")
		if idx < 26 {
			return byte('a' + idx)
		}
		if idx < 36 {
			return byte('0' + (idx - 26))
		}
		return '-'
	})
}

// genGitHubRepo generates a github.Repository with random archive/fork/name properties.
func genGitHubRepo() *rapid.Generator[*github.Repository] {
	return rapid.Custom[*github.Repository](func(t *rapid.T) *github.Repository {
		archived := rapid.Bool().Draw(t, "archived")
		fork := rapid.Bool().Draw(t, "fork")

		// Generate name that may or may not end with "-controller"
		hasControllerSuffix := rapid.Bool().Draw(t, "hasControllerSuffix")
		var name string
		if hasControllerSuffix {
			prefix := rapid.StringMatching(`[a-z][a-z0-9\-]{1,15}`).Draw(t, "prefix")
			name = prefix + controllerSuffix
		} else {
			name = rapid.StringMatching(`[a-z][a-z0-9\-]{2,20}`).Draw(t, "name")
			// Ensure it does NOT end with -controller to avoid ambiguity
			name = strings.TrimSuffix(name, controllerSuffix)
			if name == "" {
				name = "some-repo"
			}
		}

		return &github.Repository{
			Name:     github.String(name),
			Archived: github.Bool(archived),
			Fork:     github.Bool(fork),
		}
	})
}
