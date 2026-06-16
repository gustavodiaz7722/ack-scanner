// Package discovery handles GitHub repository discovery for ACK controllers.
package discovery

import (
	"context"
	"sort"
	"strings"

	"github.com/google/go-github/v60/github"

	"github.com/aws-controllers-k8s/ack-scanner/pkg/types"
)

const controllerSuffix = "-controller"

// GitHubDiscoverer discovers ACK controller repositories.
type GitHubDiscoverer struct {
	client *github.Client
	org    string
}

// NewGitHubDiscoverer creates a new GitHubDiscoverer for the given organization.
func NewGitHubDiscoverer(client *github.Client, org string) *GitHubDiscoverer {
	return &GitHubDiscoverer{
		client: client,
		org:    org,
	}
}

// DiscoverControllers returns all non-archived, non-fork *-controller repos
// in the configured organization, sorted by service name.
func (d *GitHubDiscoverer) DiscoverControllers(ctx context.Context) ([]types.ControllerRepo, error) {
	var controllers []types.ControllerRepo

	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	for {
		repos, resp, err := d.client.Repositories.ListByOrg(ctx, d.org, opt)
		if err != nil {
			return nil, err
		}

		for _, repo := range repos {
			if !isControllerRepo(repo) {
				continue
			}
			controllers = append(controllers, types.ControllerRepo{
				RepoName:    repo.GetName(),
				ServiceName: ExtractServiceName(repo.GetName()),
			})
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	sort.Slice(controllers, func(i, j int) bool {
		return controllers[i].ServiceName < controllers[j].ServiceName
	})

	return controllers, nil
}

// isControllerRepo returns true if the repository is a non-archived, non-fork
// repo whose name ends with "-controller".
func isControllerRepo(repo *github.Repository) bool {
	if repo.GetArchived() {
		return false
	}
	if repo.GetFork() {
		return false
	}
	return strings.HasSuffix(repo.GetName(), controllerSuffix)
}

// ExtractServiceName extracts the service name from a controller repository
// name by removing the "-controller" suffix.
func ExtractServiceName(repoName string) string {
	return strings.TrimSuffix(repoName, controllerSuffix)
}
