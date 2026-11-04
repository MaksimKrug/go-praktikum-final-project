package api

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

func signinHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var request struct {
			Password string `json:"password"`
		}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&request)
		if err != nil {
			writeJSON(w, map[string]string{
				"error": "Not correct JSON: " + err.Error(),
			})
			return
		}
		var extra any
		if err := decoder.Decode(&extra); err != io.EOF {
			writeJSON(w, map[string]string{
				"error": "Тело запроса должно содержать один JSON-объект",
			})
			return
		}

		password := request.Password

		envPassword := os.Getenv("TODO_PASSWORD")
		if envPassword == "" {
			writeJSON(w, map[string]string{
				"error": "Authentication off",
			})
			return
		}
		if password == envPassword {
			token, err := createToken(password)
			if err != nil {
				writeJSON(w, map[string]string{
					"error": err.Error(),
				})
				return
			}
			writeJSON(w, map[string]string{"token": token})
		} else {
			writeJSON(w, map[string]string{
				"error": "Wrong password",
			})
		}
	default:
		w.Header().Set("Allow", http.MethodPost)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusMethodNotAllowed)
		writeJSON(w, map[string]string{
			"error": "Request method not supported",
		})
	}
}
