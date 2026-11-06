package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/MaksimKrug/go-praktikum-final-project/pkg/db"
)

func checkDate(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(db.DateFormat)
	}

	t, err := time.Parse(db.DateFormat, task.Date)
	if err != nil {
		return fmt.Errorf("task date not parsed: %w", err)
	}

	var next string
	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	if afterNow(now, t) {
		if task.Repeat == "" {
			task.Date = now.Format(db.DateFormat)
		} else {
			task.Date = next
		}
	}

	return nil
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&task)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid JSON: " + err.Error(),
		})
		return
	}

	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "request body must contain one JSON object",
		})
		return
	}

	if task.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "title not specified",
		})
		return
	}

	if err := checkDate(&task); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeDBError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{
		"id": strconv.FormatInt(id, 10),
	})

}
