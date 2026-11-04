package api

import (
	"fmt"
	"net/http"
	"time"
)

const dateFormat = "20060102"

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	nowParam := r.URL.Query().Get("now")
	if nowParam != "" {
		parsedNow, err := time.Parse(dateFormat, nowParam)
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

	fmt.Fprint(w, newDate)

}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	default:
		w.Header().Set("Allow", http.MethodPost)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusMethodNotAllowed)
		writeJSON(w, map[string]string{
			"error": "Request method not supported",
		})
	}
}

func Init() {
	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/task", taskHandler)
}
