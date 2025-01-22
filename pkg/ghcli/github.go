package ghcli

import (
	"context"
	"fmt"
	"github.com/google/go-github/v55/github"
	"net/http"
)

//go:generate mockgen -destination=./mocks/mock_github_client.go -package=ghcli_mocks github.com/andrelince/github-proxy/pkg/ghcli GithubClient
type GithubClient interface {
	CreateRepository(ctx context.Context, in RepositoryInput) (Repository, error)
	ListRepositories(ctx context.Context) ([]Repository, error)
	DeleteRepository(ctx context.Context, name string) error
	ListOpenPRs(ctx context.Context, owner, repo string, num int) ([]PullRequest, error)
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
		fmt.Println(err)
		return Repository{}, err
	} else if resp.StatusCode != http.StatusCreated {
		fmt.Println(resp)
		return Repository{}, fmt.Errorf("failed to create repository: %s", resp.Status)
	} else if out == nil {
		return Repository{}, fmt.Errorf("repository is empty")
	}
	return Repository{
		ID:          out.GetID(),
		Name:        out.GetName(),
		Description: out.GetDescription(),
		Private:     out.GetPrivate(),
	}, nil
}

func (c *GHClient) ListRepositories(ctx context.Context) ([]Repository, error) {
	out, resp, err := c.client.Repositories.List(ctx, "andrelince", &github.RepositoryListOptions{
		Visibility: "public",
	})
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list repositories: %s", resp.Status)
	}

	res := make([]Repository, len(out))
	for i := range out {
		res[i] = Repository{
			ID:          out[i].GetID(),
			Name:        out[i].GetName(),
			Description: out[i].GetDescription(),
			Private:     out[i].GetPrivate(),
		}
	}

	return res, nil
}

func (c *GHClient) DeleteRepository(ctx context.Context, name string) error {
	resp, err := c.client.Repositories.Delete(ctx, "andrelince", name)
	if err != nil {
		return err
	} else if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete repository: %s", resp.Status)
	}
	return nil
}

func (c *GHClient) ListOpenPRs(ctx context.Context, owner, repo string, num int) ([]PullRequest, error) {
	if num > 10 { // enforce max number
		num = 10
	}
	out, resp, err := c.client.PullRequests.List(ctx, owner, repo, &github.PullRequestListOptions{
		State: "open",
		ListOptions: github.ListOptions{
			PerPage: num,
			Page:    1,
		},
	})
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list prs: %s", resp.Status)
	}

	res := make([]PullRequest, len(out))
	for i := range out {
		res[i] = PullRequest{
			ID:          out[i].GetID(),
			Title:       out[i].GetTitle(),
			Body:        out[i].GetBody(),
			Contributor: out[i].GetUser().GetName(),
		}
	}

	return res, nil
}
