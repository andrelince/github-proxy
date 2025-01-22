package rest

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/andrelince/github-proxy/pkg/ghcli"
	"github.com/andrelince/github-proxy/rest/definitions"
)

type Handler struct {
	githubClient ghcli.GithubClient
}

func NewHandler(githubClient ghcli.GithubClient) Handler {
	return Handler{
		githubClient: githubClient,
	}
}

func (h Handler) Health(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	if _, err := writer.Write([]byte(`OK`)); err != nil {
		http.Error(writer, "internal Server Error", http.StatusInternalServerError)
	}
}

func (h Handler) CreateRepo(writer http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, "failed to read request body", http.StatusBadRequest)
		return
	}
	defer request.Body.Close()

	var repoRequest definitions.CreateRepositoryRequest
	if err := json.Unmarshal(body, &repoRequest); err != nil {
		http.Error(writer, "invalid JSON format", http.StatusBadRequest)
		return
	}

	out, err := h.githubClient.CreateRepository(request.Context(), ghcli.RepositoryInput{
		Name:        repoRequest.Name,
		Description: repoRequest.Description,
		Private:     repoRequest.Private,
	})

	if err != nil {
		http.Error(writer, "failed to create repository", http.StatusInternalServerError)
		return
	}

	resp := definitions.RepositoryResponse{
		ID:          out.ID,
		Name:        out.Name,
		Description: out.Description,
		Private:     out.Private,
	}

	bytes, err := json.Marshal(resp)
	if err != nil {
		http.Error(writer, "failed construct response", http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusCreated)
	writer.Write(bytes)
}

func (h Handler) ListRepos(writer http.ResponseWriter, request *http.Request) {
	out, err := h.githubClient.ListRepositories(request.Context())
	if err != nil {
		http.Error(writer, "failed to get repositories", http.StatusInternalServerError)
		return
	}

	resp := make([]definitions.RepositoryResponse, len(out))

	for i := range out {
		resp[i] = definitions.RepositoryResponse{
			ID:          out[i].ID,
			Name:        out[i].Name,
			Description: out[i].Description,
			Private:     out[i].Private,
		}
	}

	bytes, err := json.Marshal(resp)
	if err != nil {
		http.Error(writer, "failed construct response", http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write(bytes)
}

func (h Handler) DeleteRepo(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	name := vars["name"]

	if name == "" && strings.HasPrefix(name, "testonly") {
		// extra safety check
		http.Error(writer, "invalid repository name", http.StatusBadRequest)
		return
	}

	if err := h.githubClient.DeleteRepository(request.Context(), name); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}

func (h Handler) ListOpenPRs(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	owner := vars["owner"]
	repository := vars["repository"]

	if owner == "" || repository == "" {
		http.Error(writer, "invalid owner/repository name", http.StatusBadRequest)
		return
	}

	num, err := strconv.Atoi(request.URL.Query().Get("num"))
	if num == 0 {
		num = 10
	}

	out, err := h.githubClient.ListOpenPRs(request.Context(), owner, repository, num)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := make([]definitions.PRResponse, len(out))

	for i := range out {
		resp[i] = definitions.PRResponse{
			ID:          out[i].ID,
			Title:       out[i].Title,
			Body:        out[i].Body,
			Contributor: out[i].Contributor,
		}
	}

	bytes, err := json.Marshal(resp)
	if err != nil {
		http.Error(writer, "failed construct response", http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write(bytes)
}
