package rest

import (
	"net/http"
)

type Handler struct{}

func NewHandler() Handler {
	return Handler{}
}

func (h Handler) Health(writer http.ResponseWriter, request *http.Request) {
	if _, err := writer.Write([]byte(`OK`)); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
}
