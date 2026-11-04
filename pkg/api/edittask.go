package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/MaksimKrug/go-praktikum-final-project/pkg/db"
)

func editTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&task)
	if err != nil {
		writeJSON(w, map[string]string{
			"error": "Not correct JSON: " + err.Error(),
		})
		return
	}

	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		writeJSON(w, map[string]string{
			"error": "Request body contains more than one JSON",
		})
		return
	}

	if task.Title == "" {
		writeJSON(w, map[string]string{
			"error": "Title not specified",
		})
		return
	}

	parsedId, err := strconv.ParseInt(task.ID, 10, 64)
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

	if err := checkDate(&task); err != nil {
		writeJSON(w, map[string]string{
			"error": err.Error(),
		})
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{
			"error": err.Error(),
		})
		return
	}
	writeJSON(w, struct{}{})
}
