package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/MihaiBlebea/trading-platform/account"
	"github.com/MihaiBlebea/trading-platform/activity"
	"github.com/MihaiBlebea/trading-platform/order"
	"github.com/MihaiBlebea/trading-platform/symbols"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
)

type Order struct {
	order.Order
	Title string `json:"title"`
}

type PlaceOrderRequest struct {
	Symbol     string  `json:"symbol"`
	Amount     float64 `json:"amount"`
	Type       string  `json:"type"`
	Direction  string  `json:"direction"`
	Quantity   int     `json:"quantity"`
	StopLoss   float64 `json:"stop-loss"`
	TakeProfit float64 `json:"take-profit"`
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
	Success bool    `json:"success"`
	Error   string  `json:"error,omitempty"`
	Orders  []Order `json:"orders"`
}

func PlaceOrderHandler(cont *dig.Container) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req PlaceOrderRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			serverError(w, cont, err)
			return
		}

		header := r.Header.Get("Authorization")
		if header == "" {
			serverError(w, cont, errors.New("could not find authorization header"))
			return
		}
		apiToken := strings.Split(header, "Bearer ")[1]

		err = cont.Invoke(func(
			symbolService *symbols.Service,
			orderPlacer *activity.OrderPlacer,
			logger *logrus.Logger) {

			if !symbolService.Exists(strings.ToUpper(req.Symbol)) {
				serverError(w, cont, errors.New("symbol not found"))
				return
			}

			order, err := orderPlacer.PlaceOrder(
				apiToken,
				req.Type,
				req.Direction,
				req.Symbol,
				req.Amount,
				float64(req.Quantity),
				req.StopLoss,
				req.TakeProfit,
			)
			if err != nil {
				serverError(w, cont, err)
				return
			}

			resp := OrderResponse{
				Success: true,
				Order:   order,
			}
			sendResponse(w, logger, resp, http.StatusOK)
		})
		if err != nil {
			serverError(w, cont, err)
			return
		}
	})
}

func CancelOrderHandler(cont *dig.Container) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req CancelOrderRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			serverError(w, cont, err)
			return
		}

		header := r.Header.Get("Authorization")
		if header == "" {
			serverError(w, cont, errors.New("could not find authorization header"))
			return
		}
		apiToken := strings.Split(header, "Bearer ")[1]

		err = cont.Invoke(func(orderCanceller *activity.OrderCanceller, logger *logrus.Logger) {
			order, err := orderCanceller.CancelOrder(
				apiToken,
				req.OrderID,
			)
			if err != nil {
				serverError(w, cont, err)
				return
			}

			resp := OrderResponse{
				Success: true,
				Order:   order,
			}
			sendResponse(w, logger, resp, http.StatusOK)
		})
		if err != nil {
			serverError(w, cont, err)
			return
		}
	})
}

func OrdersHandler(cont *dig.Container) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			serverError(w, cont, errors.New("could not find authorization header"))
			return
		}
		apiToken := strings.Split(header, "Bearer ")[1]

		err := cont.Invoke(func(
			accountRepo *account.AccountRepo,
			orderRepo *order.OrderRepo,
			symbolRepo *symbols.SymbolRepo,
			logger *logrus.Logger) {

			account, err := accountRepo.WithToken(apiToken)
			if err != nil {
				serverError(w, cont, err)
				return
			}

			orders, err := orderRepo.WithAccountId(account.ID)
			if err != nil {
				serverError(w, cont, err)
				return
			}

			resp := OrdersResponse{
				Success: true,
				Orders:  []Order{},
			}

			for _, o := range orders {
				var title string
				if s, err := symbolRepo.WithSymbol(o.Symbol); err == nil {
					title = s.Title
				} else {
					logger.Error(err)
				}

				resp.Orders = append(resp.Orders, Order{o, title})
			}

			sendResponse(w, logger, resp, http.StatusOK)
		})
		if err != nil {
			serverError(w, cont, err)
			return
		}
	})
}
