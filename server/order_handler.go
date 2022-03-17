package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/MihaiBlebea/trading-platform/account"
	"github.com/MihaiBlebea/trading-platform/order"
)

type PlaceOrderRequest struct {
	Symbol    string  `json:"symbol"`
	Amount    float32 `json:"amount"`
	Type      string  `json:"type"`
	Direction string  `json:"direction"`
}

type PlaceOrderResponse struct {
	Success bool         `json:"success"`
	Error   string       `json:"error,omitempty"`
	Order   *order.Order `json:"order,omitempty"`
}

func placeOrderHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req PlaceOrderRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			serverError(w, err)
			return
		}

		header := r.Header.Get("Authorization")
		apiToken := strings.Split(header, "Bearer ")[1]

		accountRepo, err := account.NewAccountRepo()
		if err != nil {
			serverError(w, err)
			return
		}

		account, err := accountRepo.WithToken(apiToken)
		if err != nil {
			serverError(w, err)
			return
		}

		repo, err := order.NewOrderRepo()
		if err != nil {
			serverError(w, err)
			return
		}

		order := order.NewOrder(account.ID, req.Type, req.Direction, req.Amount, req.Symbol)
		order, err = repo.Save(order)
		if err != nil {
			serverError(w, err)
			return
		}

		resp := PlaceOrderResponse{
			Success: true,
			Order:   order,
		}
		sendResponse(w, resp, http.StatusOK)
	})
}

func serverError(w http.ResponseWriter, err error) {
	if err != nil {
		resp := PlaceOrderResponse{
			Success: false,
			Error:   err.Error(),
		}
		sendResponse(w, resp, http.StatusInternalServerError)
	}
}
