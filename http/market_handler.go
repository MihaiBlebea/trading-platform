package http

import (
	"net/http"

	"github.com/MihaiBlebea/trading-platform/di"
)

type MarketStatusResponse struct {
	Success      bool   `json:"success"`
	Error        string `json:"error,omitempty"`
	MarketOpenUS bool   `json:"market_open_us"`
}

func MarketStatusHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		di := di.NewContainer()

		status, err := di.GetMarketStatus()
		if err != nil {
			serverError(w, err)
			return
		}

		isOpen := status.IsOpen()

		resp := MarketStatusResponse{
			Success:      true,
			MarketOpenUS: isOpen,
		}
		sendResponse(w, resp, http.StatusOK)
	})
}
