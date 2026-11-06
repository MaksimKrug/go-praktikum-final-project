package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MaksimKrug/go-praktikum-final-project/pkg/db"
)

func TestDatabaseErrorStatus(t *testing.T) {
	for _, tc := range []struct {
		name string
		err  error
		code int
	}{
		{"missing", db.ErrNotFound, http.StatusNotFound},
		{"wrapped missing", fmt.Errorf("delete task: %w", db.ErrNotFound), http.StatusNotFound},
		{"database failure", errors.New("private database details"), http.StatusInternalServerError},
	} {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			writeDBError(w, tc.err)
			if w.Code != tc.code {
				t.Fatalf("status = %d, want %d", w.Code, tc.code)
			}
			var body map[string]string
			if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil || body["error"] == "" {
				t.Fatalf("invalid error response: %s", w.Body.String())
			}
			if tc.code == http.StatusInternalServerError && body["error"] != "internal server error" {
				t.Fatal("internal database details exposed")
			}
		})
	}
}
