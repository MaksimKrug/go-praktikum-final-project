package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/MaksimKrug/go-praktikum-final-project/pkg/db"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
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
		if task.Repeat == "" {
			err := db.DeleteTask(id)
			if err != nil {
				writeJSON(w, map[string]string{
					"error": err.Error(),
				})
				return
			}
		} else {
			next, err := NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				writeJSON(w, map[string]string{
					"error": err.Error(),
				})
				return
			}
			err = db.UpdateDate(next, id)
			if err != nil {
				writeJSON(w, map[string]string{
					"error": err.Error(),
				})
				return
			}
		}

		writeJSON(w, struct{}{})
	default:
		w.Header().Set("Allow", http.MethodPost)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusMethodNotAllowed)
		writeJSON(w, map[string]string{
			"error": "Request method not supported",
		})
	}
}
