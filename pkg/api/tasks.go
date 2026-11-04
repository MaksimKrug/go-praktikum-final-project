package api

import (
	"net/http"

	"github.com/MaksimKrug/go-praktikum-final-project/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tasks, err := db.Tasks(10, r.URL.Query().Get("search"))
		if err != nil {
			writeJSON(w, map[string]string{
				"error": err.Error(),
			})
			return
		}

		writeJSON(w, TasksResp{Tasks: tasks})
	default:
		w.Header().Set("Allow", http.MethodGet)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusMethodNotAllowed)
		writeJSON(w, map[string]string{
			"error": "Request method not supported",
		})
	}

}
