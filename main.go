package main

import (
	"net/http"

	"github.com/andrelince/github-proxy/di"
)

func main() {
	c, err := di.NewDI()
	if err != nil {
		panic(err)
	}

	if err := c.Invoke(start); err != nil {
		panic(err)
	}
}

func start(server *http.Server) error {
	if err := server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
