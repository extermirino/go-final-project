package api

import (
	"errors"
	"net/http"
	"time"

	"go1f/pkg/db"
)

func doneTask(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
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

	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			writeError(w, err)
			return
		}

		writeJson(w, map[string]string{})
		return
	}

	date, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		writeError(w, err)
		return
	}

	err = db.UpdateDate(date, task.ID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJson(w, map[string]string{})
}
