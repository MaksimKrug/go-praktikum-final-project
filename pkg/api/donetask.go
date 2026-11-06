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
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
			return
		}
		if parsedId <= 0 {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "id should be greater than zero",
			})
			return
		}
		task, err := db.GetTask(id)
		if err != nil {
			writeDBError(w, err)
			return
		}
		if task.Repeat == "" {
			err := db.DeleteTask(id)
			if err != nil {
				writeDBError(w, err)
				return
			}
		} else {
			next, err := NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				writeJSON(w, http.StatusInternalServerError, map[string]string{
					"error": err.Error(),
				})
				return
			}
			err = db.UpdateDate(next, id)
			if err != nil {
				writeDBError(w, err)
				return
			}
		}

		writeJSON(w, http.StatusOK, struct{}{})
	default:
		w.Header().Set("Allow", http.MethodPost)
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "request method not supported",
		})
	}
}
