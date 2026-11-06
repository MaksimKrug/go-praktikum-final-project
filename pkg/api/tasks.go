package api

import (
	"net/http"

	"github.com/MaksimKrug/go-praktikum-final-project/pkg/db"
)

const tasksLimit = 10

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tasks, err := db.Tasks(tasksLimit, r.URL.Query().Get("search"))
		if err != nil {
			writeDBError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, TasksResp{Tasks: tasks})
	default:
		w.Header().Set("Allow", http.MethodGet)
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "request method not supported",
		})
	}

}
