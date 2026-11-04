package api

import (
	"net/http"
	"strconv"

	"github.com/MaksimKrug/go-praktikum-final-project/pkg/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	parsedId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		writeJSON(w, map[string]string{
			"error": err.Error(),
		})
		return
	}
	if parsedId <= 0 {
		writeJSON(w, map[string]string{
			"error": "Id should be greater than zero",
		})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{
			"error": err.Error(),
		})
		return
	}
	writeJSON(w, task)
}
