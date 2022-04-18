package http_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/MihaiBlebea/trading-platform/account"
	"github.com/MihaiBlebea/trading-platform/di"
	handler "github.com/MihaiBlebea/trading-platform/http"
	"github.com/MihaiBlebea/trading-platform/pos"
	"github.com/MihaiBlebea/trading-platform/symbols"
	"github.com/gorilla/mux"
)

func init() {
	os.Setenv("APP_ENV", "test")

	cont = di.BuildContainer()
}

func TestGetPortfolio(t *testing.T) {
	defer tearDown(t)

	// Create an account
	acc, err := account.NewAccount("mihaib", "mihai@gmail.com", "1234")
	if err != nil {
		t.Error(err)
		return
	}

	var position *pos.Position

	err = cont.Invoke(func(
		accountRepo *account.AccountRepo,
		symbolRepo *symbols.SymbolRepo,
		positionRepo *pos.PositionRepo) {

		acc, err := accountRepo.Save(acc)
		if err != nil {
			t.Error(err)
			return
		}

		_, err = symbolRepo.Save(symbols.NewSymbol("Apple", "Apple Inc.", "Computers", "USD", "AAPL"))
		if err != nil {
			t.Error(err)
			return
		}

		position = pos.NewPosition(acc.ID, "AAPL", 100)
		position.IncrementQuantity(10, float64(1077.44))

		_, err = positionRepo.Save(position)
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
	r.Handle("/api/v1/positions", handler.PositionsHandler(cont)).Methods(http.MethodGet)

	ts := httptest.NewServer(r)
	defer ts.Close()

	client := http.Client{}
	req, err := http.NewRequest("GET", ts.URL+"/api/v1/positions", nil)
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

	response := handler.PositionsResponse{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Error(err)
		return
	}

	if response.Success != true {
		t.Errorf("expected success to be true, got: %v", response.Success)
		return
	}

	if len(response.Positions) != 1 {
		t.Errorf("expected positions length to be 1, got: %v", len(response.Positions))
		return
	}

	if response.Positions[0].Symbol != "AAPL" {
		t.Errorf("expected position symbol to be AAPL, got: %v", response.Positions[0].Symbol)
	}

	if response.Positions[0].Quantity != position.Quantity {
		t.Errorf(
			"expected position quantity to be %v, got: %v",
			position.Quantity,
			response.Positions[0].Quantity,
		)
	}

	totalValue := float64(1210)
	if response.Positions[0].TotalValue != totalValue {
		t.Errorf(
			"expected position total value to be %v, got: %v",
			totalValue,
			response.Positions[0].TotalValue,
		)
	}
}
