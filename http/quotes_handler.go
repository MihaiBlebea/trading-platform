package http

import (
	"errors"
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
			serverError(w, errors.New("Invalid start"))
			return
		}
		symbol := query.Get("symbol")
		if symbol == "" {
			serverError(w, errors.New("Invalid symbol"))
			return
		}

		qs, err := q.GetQuotes(strings.ToUpper(symbol), startDate, "1m")
		if err != nil {
			serverError(w, err)
			return
		}

		resp := QuotesResponse{
			Success: true,
			Data:    qs,
		}
		sendResponse(w, resp, http.StatusOK)
	})
}
