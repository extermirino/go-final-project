package api

import (
	"go1f/pkg/db"
	"net/http"
)

type TasksResp struct {
	Tasks []db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue("search")

	tasks, err := db.Tasks(50, search)
	if err != nil {
		writeError(w, err)
		return
	}

	if tasks == nil {
		tasks = make([]db.Task, 0)
	}

	writeJson(w, TasksResp{
		Tasks: tasks,
	})
}
