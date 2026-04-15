package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go1f/pkg/db"
)

func addTask(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeError(w, errors.New("invalid json"))
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

	id, err := db.AddTask(&task)
	if err != nil {
		writeError(w, errors.New("db error"))
		return
	}

	writeJson(w, map[string]string{
		"id": fmt.Sprint(id),
	})
}

func checkDate(task *db.Task) error {
	now := time.Now().Format("20060102")

	if task.Date == "" {
		task.Date = now
	}

	_, err := time.Parse("20060102", task.Date)
	if err != nil {
		return err
	}

	if task.Date < now {
		if task.Repeat == "" {
			task.Date = now
		} else {
			next, err := NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				return err
			}
			task.Date = next
		}
	}

	return nil
}
