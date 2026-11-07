package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/MaksimKrug/go-praktikum-final-project/pkg/config"
	"github.com/MaksimKrug/go-praktikum-final-project/pkg/db"
)

type apiHandlers struct {
	config config.Config
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:

		now := time.Now()
		nowParam := r.URL.Query().Get("now")
		if nowParam != "" {
			parsedNow, err := time.Parse(db.DateFormat, nowParam)
			if err != nil {
				http.Error(w, "now parameter not correct: "+err.Error(), http.StatusBadRequest)
				return
			}
			now = parsedNow
		}

		dateParam := r.URL.Query().Get("date")
		repeatParam := r.URL.Query().Get("repeat")

		newDate, err := NextDate(now, dateParam, repeatParam)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if _, err := fmt.Fprint(w, newDate); err != nil {
			log.Printf("write response: %v", err)
		}
	default:
		w.Header().Set("Allow", http.MethodGet)
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "request method not supported",
		})
	}

}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		editTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		w.Header().Set("Allow", "GET, POST, PUT, DELETE")
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "request method not supported",
		})
	}
}

func Init(cfg config.Config) {
	a := &apiHandlers{config: cfg}
	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/task", a.auth(taskHandler))
	http.HandleFunc("/api/tasks", a.auth(tasksHandler))
	http.HandleFunc("/api/task/done", a.auth(doneTaskHandler))
	http.HandleFunc("/api/signin", a.signinHandler)
}
