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
		task.Date = now.Format(dateFormat)
	}

	t, err := time.Parse(dateFormat, task.Date)
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
			task.Date = now.Format(dateFormat)
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

	if err := checkDate(&task); err != nil {
		writeJSON(w, map[string]string{
			"error": err.Error(),
		})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{
			"error": err.Error(),
		})
		return
	}
	writeJSON(w, map[string]string{
		"id": strconv.FormatInt(id, 10),
	})

}
