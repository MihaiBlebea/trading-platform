package http_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/MihaiBlebea/trading-platform/account"
	handler "github.com/MihaiBlebea/trading-platform/http"
	"github.com/MihaiBlebea/trading-platform/order"
	"github.com/MihaiBlebea/trading-platform/pos"
	"github.com/MihaiBlebea/trading-platform/symbols"
	"github.com/gorilla/mux"
)

func init() {
	os.Setenv("APP_ENV", "test")
}

func TestPlaceBuyOrder(t *testing.T) {
	cont := setupSuite(t)
	defer tearDown(t, cont)

	// Create an account
	acc, err := account.NewAccount("mihaib", "mihai@gmail.com", "1234")
	if err != nil {
		t.Error(err)
		return
	}

	err = cont.Invoke(func(
		accountRepo *account.AccountRepo,
		symbolRepo *symbols.SymbolRepo) {

		_, err := accountRepo.Save(acc)
		if err != nil {
			t.Error(err)
			return
		}
	})
	if err != nil {
		t.Error(err)
		return
	}

	r := mux.NewRouter()
	r.Handle("/api/v1/order", handler.PlaceOrderHandler(cont)).Methods(http.MethodPost)

	ts := httptest.NewServer(r)
	defer ts.Close()

	// Create request body
	payload := handler.PlaceOrderRequest{
		Symbol:    "AAPL",
		Amount:    1000.00,
		Direction: "buy",
		Type:      "limit",
	}
	b, err := json.Marshal(payload)
	if err != nil {
		t.Error(err)
		return
	}

	client := http.Client{}
	req, err := http.NewRequest("POST", ts.URL+"/api/v1/order", bytes.NewBuffer(b))
	if err != nil {
		t.Error(err)
		return
	}

	req.Header.Set(
		"Authorization",
		fmt.Sprintf("Bearer %s", acc.ApiToken),
	)

	res, err := client.Do(req)
	if err != nil {
		t.Error(err)
		return
	}

	response := handler.OrderResponse{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Error(err)
		return
	}

	if response.Success != true {
		t.Errorf("expected success to be true, got: %v", response.Success)
	}

	if response.Order.Status != "pending" {
		t.Errorf("expected order status to be pending, got: %v", response.Order.Status)
	}

	if response.Order.Direction != "buy" {
		t.Errorf("expected order direction to be buy, got: %v", response.Order.Direction)
	}

	if response.Order.Amount != float64(payload.Amount) {
		t.Errorf(
			"expected order amount to be %v, got: %v",
			payload.Amount,
			response.Order.Amount,
		)
	}

	if response.Order.FillPrice != 0 {
		t.Errorf("expected order fill price to be 0, got: %v", response.Order.FillPrice)
	}

	if response.Order.Symbol != payload.Symbol {
		t.Errorf(
			"expected order symbol to be %v, got: %v",
			payload.Symbol,
			response.Order.Symbol,
		)
	}

	if response.Order.Quantity != 0 {
		t.Errorf("expected order quantity to be 0, got: %v", response.Order.Quantity)
	}
}

func TestPlaceSellOrder(t *testing.T) {
	cont := setupSuite(t)
	defer tearDown(t, cont)

	// Create an account
	acc, err := account.NewAccount("mihaib", "mihai@gmail.com", "1234")
	if err != nil {
		t.Error(err)
		return
	}

	err = cont.Invoke(func(
		accountRepo *account.AccountRepo,
		symbolRepo *symbols.SymbolRepo,
		positionRepo *pos.PositionRepo) {

		acc, err := accountRepo.Save(acc)
		if err != nil {
			t.Error(err)
			return
		}

		_, err = positionRepo.Save(pos.NewPosition(acc.ID, "AAPL", 100, 144.44))
		if err != nil {
			t.Error(err)
			return
		}
	})
	if err != nil {
		t.Error(err)
		return
	}

	r := mux.NewRouter()
	r.Handle("/api/v1/order", handler.PlaceOrderHandler(cont)).Methods(http.MethodPost)

	ts := httptest.NewServer(r)
	defer ts.Close()

	// Create request body
	payload := handler.PlaceOrderRequest{
		Symbol:    "AAPL",
		Quantity:  10,
		Direction: "sell",
		Type:      "limit",
	}
	b, err := json.Marshal(payload)
	if err != nil {
		t.Error(err)
		return
	}

	client := http.Client{}
	req, err := http.NewRequest("POST", ts.URL+"/api/v1/order", bytes.NewBuffer(b))
	if err != nil {
		t.Error(err)
		return
	}

	req.Header.Set(
		"Authorization",
		fmt.Sprintf("Bearer %s", acc.ApiToken),
	)

	res, err := client.Do(req)
	if err != nil {
		t.Error(err)
		return
	}

	response := handler.OrderResponse{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Error(err)
		return
	}

	if response.Success != true {
		t.Errorf("expected success to be true, got: %v", response.Success)
	}

	if response.Order.Status != "pending" {
		t.Errorf("expected order status to be pending, got: %v", response.Order.Status)
	}

	if response.Order.Direction != "sell" {
		t.Errorf("expected order direction to be sell, got: %v", response.Order.Direction)
	}

	if response.Order.FillPrice != 0 {
		t.Errorf("expected order fill price to be 0, got: %v", response.Order.FillPrice)
	}

	if response.Order.Symbol != payload.Symbol {
		t.Errorf(
			"expected order symbol to be %v, got: %v",
			payload.Symbol,
			response.Order.Symbol,
		)
	}

	if response.Order.Quantity != payload.Quantity {
		t.Errorf(
			"expected order quantity to be %d, got: %v",
			payload.Quantity,
			response.Order.Quantity,
		)
	}
}

func TestFetchOrders(t *testing.T) {
	cont := setupSuite(t)
	defer tearDown(t, cont)

	// Create an account
	acc, err := account.NewAccount("mihaib", "mihai@gmail.com", "1234")
	if err != nil {
		t.Error(err)
		return
	}

	err = cont.Invoke(func(accountRepo *account.AccountRepo, orderRepo *order.OrderRepo) {
		acc, err := accountRepo.Save(acc)
		if err != nil {
			t.Error(err)
			return
		}

		_, err = orderRepo.Save(order.NewBuyOrder(acc.ID, "limit", "AAPL", float64(1000)))
		if err != nil {
			t.Error(err)
			return
		}
	})
	if err != nil {
		t.Error(err)
		return
	}

	r := mux.NewRouter()
	r.Handle("/api/v1/orders", handler.OrdersHandler(cont)).Methods(http.MethodGet)

	ts := httptest.NewServer(r)
	defer ts.Close()

	client := http.Client{}
	req, err := http.NewRequest("GET", ts.URL+"/api/v1/orders", nil)
	if err != nil {
		t.Error(err)
		return
	}

	req.Header.Set(
		"Authorization",
		fmt.Sprintf("Bearer %s", acc.ApiToken),
	)

	res, err := client.Do(req)
	if err != nil {
		t.Error(err)
		return
	}

	response := handler.OrdersResponse{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Error(err)
		return
	}

	if response.Success != true {
		t.Errorf("expected success to be true, got: %v", response.Success)
	}

	if len(response.Orders) != 1 {
		t.Errorf("expected orders length to be 1, got: %v", len(response.Orders))
	}

	if response.Orders[0].Symbol != "AAPL" {
		t.Errorf("expected order symbol to be APPL, got: %v", response.Orders[0].Symbol)
	}

	if response.Orders[0].Title != appleSymbol.Title {
		t.Errorf(
			"expected order tite to be %s, got: %s",
			appleSymbol.Title,
			response.Orders[0].Title,
		)
	}
}

func TestCancelOrder(t *testing.T) {
	cont := setupSuite(t)
	defer tearDown(t, cont)

	// Create an account
	acc, err := account.NewAccount("mihaib", "mihai@gmail.com", "1234")
	if err != nil {
		t.Error(err)
		return
	}

	err = cont.Invoke(func(accountRepo *account.AccountRepo, orderRepo *order.OrderRepo) {
		acc, err := accountRepo.Save(acc)
		if err != nil {
			t.Error(err)
			return
		}

		_, err = orderRepo.Save(order.NewBuyOrder(acc.ID, "limit", "AAPL", float64(1000)))
		if err != nil {
			t.Error(err)
			return
		}
	})
	if err != nil {
		t.Error(err)
		return
	}

	r := mux.NewRouter()
	r.Handle("/api/v1/order/cancel", handler.CancelOrderHandler(cont)).Methods(http.MethodPut)

	ts := httptest.NewServer(r)
	defer ts.Close()

	payload := handler.CancelOrderRequest{
		OrderID: 1,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		t.Error(err)
		return
	}

	client := http.Client{}
	req, err := http.NewRequest("PUT", ts.URL+"/api/v1/order/cancel", bytes.NewBuffer(b))
	if err != nil {
		t.Error(err)
		return
	}

	req.Header.Set(
		"Authorization",
		fmt.Sprintf("Bearer %s", acc.ApiToken),
	)

	res, err := client.Do(req)
	if err != nil {
		t.Error(err)
		return
	}

	response := handler.OrdersResponse{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Error(err)
		return
	}

	if response.Success != true {
		t.Errorf("expected success to be true, got: %v", response.Success)
	}
}
