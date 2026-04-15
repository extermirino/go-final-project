package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"go1f/pkg/db"
)

func updateTask(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeError(w, errors.New("invalid json"))
		return
	}

	if task.ID == "" {
		writeError(w, errors.New("id is empty"))
		return
	}

	if task.Title == "" {
		writeError(w, errors.New("title is empty"))
		return
	}

	err = checkDate(&task)
	if err != nil {
		writeError(w, err)
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJson(w, map[string]string{})
}
