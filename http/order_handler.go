package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/MihaiBlebea/trading-platform/di"
	"github.com/MihaiBlebea/trading-platform/order"
)

type PlaceOrderRequest struct {
	Symbol     string  `json:"symbol"`
	Amount     float32 `json:"amount"`
	Type       string  `json:"type"`
	Direction  string  `json:"direction"`
	Quantity   int     `json:"quantity"`
	StopLoss   float32 `json:"stop-loss"`
	TakeProfit float32 `json:"take-profit"`
}

type CancelOrderRequest struct {
	OrderID int `json:"order_id"`
}

type OrderResponse struct {
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

		di := di.NewContainer()

		symbolService, err := di.GetSymbolService()
		if err != nil {
			serverError(w, err)
			return
		}

		if !symbolService.Exists(strings.ToUpper(req.Symbol)) {
			serverError(w, errors.New("symbol not found"))
			return
		}

		orderPlacer, err := di.GetOrderPlacer()
		if err != nil {
			serverError(w, err)
			return
		}

		order, err := orderPlacer.PlaceOrder(
			apiToken,
			req.Type,
			req.Direction,
			req.Symbol,
			req.Amount,
			req.Quantity,
			req.StopLoss,
			req.TakeProfit,
		)
		if err != nil {
			serverError(w, err)
			return
		}

		resp := OrderResponse{
			Success: true,
			Order:   order,
		}
		sendResponse(w, resp, http.StatusOK)
	})
}

func cancelOrderHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req CancelOrderRequest

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

		di := di.NewContainer()

		orderCanceller, err := di.GetOrderCanceller()
		if err != nil {
			serverError(w, err)
			return
		}
		order, err := orderCanceller.CancelOrder(
			apiToken,
			req.OrderID,
		)
		if err != nil {
			serverError(w, err)
			return
		}

		resp := OrderResponse{
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

		di := di.NewContainer()

		accountRepo, err := di.GetAccountRepo()
		if err != nil {
			serverError(w, err)
			return
		}
		orderRepo, err := di.GetOrderRepo()
		if err != nil {
			serverError(w, err)
			return
		}

		account, err := accountRepo.WithToken(apiToken)
		if err != nil {
			serverError(w, err)
			return
		}

		orders, err := orderRepo.WithAccountId(account.ID)
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
		resp := OrderResponse{
			Success: false,
			Error:   err.Error(),
		}
		sendResponse(w, resp, http.StatusInternalServerError)
	}
}
