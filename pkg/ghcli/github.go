package ghcli

import (
	"context"
	"fmt"
	"github.com/andrelince/github-proxy/pkg/ptr"
	"github.com/google/go-github/v55/github"
	"net/http"
)

//go:generate mockgen -destination=./mocks/mock_github_client.go -package=ghcli_mocks github.com/andrelince/github-proxy/pkg/ghcli GithubClient
type GithubClient interface {
	CreateRepository(ctx context.Context, in RepositoryInput) (Repository, error)
	ListRepositories(ctx context.Context) ([]Repository, error)
}

type GHClient struct {
	client *github.Client
}

func NewGitHubClient(personalAccessToken string) *GHClient {
	client := github.NewClient(nil).
		WithAuthToken(personalAccessToken)

	return &GHClient{
		client: client,
	}
}

func (c *GHClient) CreateRepository(ctx context.Context, in RepositoryInput) (Repository, error) {
	out, resp, err := c.client.Repositories.Create(ctx, "", &github.Repository{
		Name:        github.String(in.Name),
		Description: github.String(in.Description),
		Private:     github.Bool(in.Private),
	})
	if err != nil {
		return Repository{}, err
	} else if resp.StatusCode != http.StatusCreated {
		return Repository{}, fmt.Errorf("failed to create repository: %s", resp.Status)
	} else if out == nil {
		return Repository{}, fmt.Errorf("repository is empty")
	}
	return Repository{
		ID:          ptr.Value(out.ID, 0),
		Name:        ptr.Value(out.Name, ""),
		Description: ptr.Value(out.Description, ""),
		Private:     ptr.Value(out.Private, false),
	}, nil
}

func (c *GHClient) ListRepositories(ctx context.Context) ([]Repository, error) {
	out, resp, err := c.client.Repositories.List(ctx, "andrelince", &github.RepositoryListOptions{
		Visibility: "public",
	})
	if err != nil {
		fmt.Println(err)
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list repositories: %s", resp.Status)
	}

	res := make([]Repository, len(out))
	for i := range out {
		res[i] = Repository{
			ID:          ptr.Value(out[i].ID, 0),
			Name:        ptr.Value(out[i].Name, ""),
			Description: ptr.Value(out[i].Description, ""),
			Private:     ptr.Value(out[i].Private, false),
		}
	}

	return res, nil
}
