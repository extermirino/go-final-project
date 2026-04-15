package api

import (
	"errors"
	"net/http"

	"go1f/pkg/db"
)

func getTask(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeError(w, errors.New("id is empty"))
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJson(w, task)
}
