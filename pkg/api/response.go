package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/MaksimKrug/go-praktikum-final-project/pkg/db"
)

func writeJSON(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("write JSON response: %v", err)
	}
}

func writeDBError(w http.ResponseWriter, err error) {
	if errors.Is(err, db.ErrNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": db.ErrNotFound.Error()})
		return
	}
	log.Printf("database operation: %v", err)
	writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
}
