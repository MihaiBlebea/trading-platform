package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/MihaiBlebea/trading-platform/account"
	"github.com/MihaiBlebea/trading-platform/activity"
	"github.com/MihaiBlebea/trading-platform/order"
	"github.com/MihaiBlebea/trading-platform/pos"
)

type PlaceOrderRequest struct {
	Symbol    string  `json:"symbol"`
	Amount    float32 `json:"amount"`
	Type      string  `json:"type"`
	Direction string  `json:"direction"`
	Quantity  int     `json:"quantity"`
}

type PlaceOrderResponse struct {
	Success bool         `json:"success"`
	Error   string       `json:"error,omitempty"`
	Order   *order.Order `json:"order,omitempty"`
}

type OrdersResponse struct {
	Success bool          `json:"success"`
	Error   string        `json:"error,omitempty"`
	Orders  []order.Order `json:"orders"`
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
		if header == "" {
			serverError(w, errors.New("could not find authorization header"))
			return
		}
		apiToken := strings.Split(header, "Bearer ")[1]

		accountRepo, err := account.NewAccountRepo()
		if err != nil {
			serverError(w, err)
			return
		}

		orderRepo, err := order.NewOrderRepo()
		if err != nil {
			serverError(w, err)
			return
		}

		positionRepo, err := pos.NewPositionRepo()
		if err != nil {
			serverError(w, err)
			return
		}

		orderPlacer := activity.NewOrderPlacer(accountRepo, orderRepo, positionRepo)
		order, err := orderPlacer.PlaceOrder(
			apiToken,
			req.Type,
			req.Direction,
			req.Symbol,
			req.Amount,
			req.Quantity,
		)
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

func ordersHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			serverError(w, errors.New("could not find authorization header"))
			return
		}
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

		orders, err := repo.WithAccountId(account.ID)
		if err != nil {
			serverError(w, err)
			return
		}

		resp := OrdersResponse{
			Success: true,
			Orders:  orders,
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
