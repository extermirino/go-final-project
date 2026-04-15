package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go1f/pkg/db"
)

func addTask(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJSON(w, map[string]string{
			"error": "invalid json",
		})
		return
	}

	if task.Title == "" {
		writeJSON(w, map[string]string{
			"error": "title is empty",
		})
		return
	}

	err = checkDate(&task)
	if err != nil {
		writeJSON(w, map[string]string{
			"error": err.Error(),
		})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{
			"error": "db error",
		})
		return
	}

	writeJSON(w, map[string]string{
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
