package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/andrelince/github-proxy/config"
	"github.com/andrelince/github-proxy/pkg/ghcli"
	ghcli_mocks "github.com/andrelince/github-proxy/pkg/ghcli/mocks"
	"github.com/andrelince/github-proxy/rest/definitions"
	"github.com/cucumber/godog"
	"github.com/gorilla/mux"
	"go.uber.org/mock/gomock"
)

type restFeature struct {
	resp    *httptest.ResponseRecorder
	server  *http.Server
	baseURL string
	// mocks
	ghcliMock *ghcli_mocks.MockGithubClient
}

func (a *restFeature) init(t *testing.T) {
	r := mux.NewRouter()
	ctrl := gomock.NewController(t)

	a.ghcliMock = ghcli_mocks.NewMockGithubClient(ctrl)

	h := NewHandler(a.ghcliMock)
	a.resp = httptest.NewRecorder()
	a.server = NewRest(r, h, config.Config{Port: "8080"})
	a.baseURL = "http://localhost:8080"

	go func() {
		if err := a.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
	time.Sleep(100 * time.Millisecond)
}

func (a *restFeature) iSendAGetRequestTo(endpoint string) error {
	reqURL, err := url.JoinPath(a.baseURL, endpoint)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	a.resp.Code = resp.StatusCode
	return nil
}

func (a *restFeature) theResponseCodeShouldBe(code int) error {
	if code != a.resp.Code {
		return fmt.Errorf("unexpected response code: %d, expected: %d", code, a.resp.Code)
	}
	return nil
}

func (a *restFeature) iCreateARepository(name, description string) error {
	a.ghcliMock.
		EXPECT().
		CreateRepository(gomock.Any(), ghcli.RepositoryInput{Name: name, Description: description, Private: true}).
		Times(1).
		Return(ghcli.Repository{Name: name, Description: description, Private: true}, nil)

	reqURL, err := url.JoinPath(a.baseURL, "/repository")
	if err != nil {
		return err
	}

	payload, err := json.Marshal(
		definitions.CreateRepositoryRequest{Name: name, Description: description, Private: true})
	if err != nil {
		return fmt.Errorf("failed to marshal JSON payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, reqURL, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	a.resp.Code = resp.StatusCode
	return nil
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		Name:                 "rest api features",
		TestSuiteInitializer: nil,
		ScenarioInitializer: func(ctx *godog.ScenarioContext) {
			api := &restFeature{}

			ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
				api.init(t)
				return ctx, nil
			})
			ctx.Step(`^i send a GET request to "([^"]*)"$`, api.iSendAGetRequestTo)
			ctx.Step(`^the response code should be "([^"]*)"$`, api.theResponseCodeShouldBe)
			ctx.Step(`^i create a repository with name "([^"]*)" and description "([^"]*)"$`, api.iCreateARepository)

			ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
				return ctx, api.server.Shutdown(ctx)
			})
		},
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
