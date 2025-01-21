package rest

import (
	"net/http"
)

type Handler struct{}

func NewHandler() Handler {
	return Handler{}
}

func (h Handler) Health(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	if _, err := writer.Write([]byte(`OK`)); err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}
}
