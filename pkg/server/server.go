package server

import (
	"fmt"
	"net/http"
	"os"

	"go1f/pkg/api"
)

func Run() error {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	api.Init()

	http.Handle("/", http.FileServer(http.Dir("web")))
	return http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
