package api

import (
	"encoding/json"
	"net/http"
)

var jwtSecret = []byte("super_secret_key")

func Init() {
	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/task", auth(taskHandler))
	http.HandleFunc("/api/task/done", auth(doneTask))
	http.HandleFunc("/api/tasks", auth(tasksHandler))
	http.HandleFunc("/api/signin", signInHandler)
}

func writeJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, err error) {
	writeJson(w, map[string]string{
		"error": err.Error(),
	})
}
