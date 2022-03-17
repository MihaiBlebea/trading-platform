package http

import (
	"net/http"
	"strings"

	"github.com/MihaiBlebea/trading-platform/quotes"
)

type QuotesResponse struct {
	Success bool           `json:"success"`
	Error   string         `json:"error,omitempty"`
	Data    []quotes.Quote `json:"data,omitempty"`
}

func historicDataHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := quotes.New()

		query := r.URL.Query()
		startDate := query.Get("start")
		if startDate == "" {
			resp := QuotesResponse{
				Success: false,
				Error:   "Invalid param start",
			}
			sendResponse(w, resp, http.StatusInternalServerError)
			return
		}
		symbol := query.Get("symbol")
		if symbol == "" {
			resp := QuotesResponse{
				Success: false,
				Error:   "Invalid param symbol",
			}
			sendResponse(w, resp, http.StatusInternalServerError)
			return
		}

		qs, err := q.GetQuotes(strings.ToUpper(symbol), startDate, "1m")
		if err != nil {
			resp := QuotesResponse{
				Success: false,
				Error:   err.Error(),
			}
			sendResponse(w, resp, http.StatusInternalServerError)
			return
		}

		resp := QuotesResponse{
			Success: true,
			Data:    qs,
		}
		sendResponse(w, resp, http.StatusOK)
	})
}
