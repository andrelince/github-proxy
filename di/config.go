package di

import (
	"github.com/andrelince/github-proxy/config"
	"github.com/andrelince/github-proxy/pkg/env"
	"github.com/andrelince/github-proxy/pkg/ghcli"
	"github.com/andrelince/github-proxy/rest"
	"github.com/gorilla/mux"
	"go.uber.org/dig"
)

func buildConfig(c *dig.Container) error {
	// inject base config
	if err := c.Provide(func() (config.Config, error) {
		return env.New(config.Config{})
	}); err != nil {
		return err
	}

	// provide mux router
	if err := c.Provide(func() *mux.Router {
		return mux.NewRouter()
	}); err != nil {
		return err
	}

	// provide github client
	if err := c.Provide(func(conf config.Config) ghcli.GithubClient {
		return ghcli.NewGitHubClient(conf.GitHubAuthToken)
	}); err != nil {
		return err
	}

	// provide rest api handler
	if err := c.Provide(rest.NewHandler); err != nil {
		return err
	}

	// provide rest api server
	if err := c.Provide(rest.NewRest); err != nil {
		return err
	}

	return nil
}
