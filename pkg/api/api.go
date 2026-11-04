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

func Init() {
	http.HandleFunc("/api/nextdate", nextDayHandler)
}
