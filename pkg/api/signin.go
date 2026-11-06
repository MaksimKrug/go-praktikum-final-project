package api

import (
	"encoding/json"
	"io"
	"net/http"
)

func (a *apiHandlers) signinHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var request struct {
			Password string `json:"password"`
		}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&request)
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

		password := request.Password

		envPassword := a.config.Password
		if envPassword == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "authentication off",
			})
			return
		}
		if password == envPassword {
			token, err := createToken(password, a.config.JWTSecret)
			if err != nil {
				writeJSON(w, http.StatusInternalServerError, map[string]string{
					"error": err.Error(),
				})
				return
			}
			writeJSON(w, http.StatusOK, map[string]string{"token": token})
		} else {
			writeJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "wrong password",
			})
		}
	default:
		w.Header().Set("Allow", http.MethodPost)
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "request method not supported",
		})
	}
}
