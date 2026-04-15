package api

import (
	"errors"
	"net/http"

	"go1f/pkg/db"
)

func deleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	if id == "" {
		writeError(w, errors.New("id is empty"))
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJson(w, map[string]string{})
}
