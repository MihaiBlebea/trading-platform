package http

import (
	"encoding/json"
	"net/http"
)

func healthHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var response struct {
			OK bool `json:"ok"`
		}
		response.OK = true

		sendResponse(w, &response, http.StatusOK)
	})
}

func sendResponse(w http.ResponseWriter, resp interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	b, _ := json.Marshal(resp)

	w.Write(b)
}
