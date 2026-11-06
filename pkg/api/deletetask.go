package api

import (
	"net/http"
	"strconv"

	"github.com/MaksimKrug/go-praktikum-final-project/pkg/db"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	err = db.DeleteTask(id)
	if err != nil {
		writeDBError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, struct{}{})
}
