package http

import (
	"errors"
	"net/http"

	"github.com/MihaiBlebea/trading-platform/di"
	"github.com/MihaiBlebea/trading-platform/symbols"
)

type SymbolResponse struct {
	Success bool            `json:"success"`
	Error   string          `json:"error,omitempty"`
	Data    *symbols.Symbol `json:"symbol,omitempty"`
}

func symbolHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		di := di.NewContainer()

		query := r.URL.Query()
		symbol := query.Get("symbol")
		if symbol == "" {
			serverError(w, errors.New("invalid symbol"))
			return
		}

		symbolRepo, err := di.GetSymbolRepo()
		if err != nil {
			serverError(w, err)
			return
		}

		s, err := symbolRepo.WithSymbol(symbol)
		if err != nil {
			serverError(w, err)
			return
		}

		resp := SymbolResponse{
			Success: true,
			Data:    s,
		}
		sendResponse(w, resp, http.StatusOK)
	})
}
