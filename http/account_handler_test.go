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
	"github.com/MihaiBlebea/trading-platform/di"
	handler "github.com/MihaiBlebea/trading-platform/http"
	"github.com/MihaiBlebea/trading-platform/order"
	"github.com/MihaiBlebea/trading-platform/pos"
	"github.com/MihaiBlebea/trading-platform/symbols"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func init() {
	os.Setenv("APP_ENV", "test")
}

func tearDown(t *testing.T) {
	err := di.BuildContainer().Invoke(func(conn *gorm.DB) {
		conn.Migrator().DropTable(
			account.Account{},
			pos.Position{},
			order.Order{},
			symbols.Symbol{},
		)
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestRegisterSuccess(t *testing.T) {
	defer tearDown(t)

	r := mux.NewRouter()
	r.Handle("/api/v1/register", handler.RegisterAccountHandler()).Methods(http.MethodPost)

	ts := httptest.NewServer(r)
	defer ts.Close()

	req := handler.RegisterRequest{
		Username: "mihaib",
		Email:    "mihai@gmail.com",
		Password: "1234",
	}
	b, err := json.Marshal(req)
	if err != nil {
		t.Error(err)
		return
	}

	res, err := http.Post(ts.URL+"/api/v1/register", "application/json", bytes.NewBuffer(b))
	if err != nil {
		t.Error(err)
		return
	}

	response := handler.AccountResponse{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Error(err)
		return
	}

	if response.Success != true {
		t.Errorf("expected success to be true, got: %v", response.Success)
	}

	accountBalance := float64(10000)
	if response.Account.Balance != accountBalance {
		t.Errorf(
			"expected account balance to be %v, got: %v",
			accountBalance,
			response.Account.Balance,
		)
	}

	if response.Account.ApiToken == "" {
		t.Error("expected account token to not be empty")
	}

	if response.Account.Username != req.Username {
		t.Errorf(
			"expected username to be %v, got: %v",
			req.Username,
			response.Account.Username,
		)
	}

	if response.Account.Email != req.Email {
		t.Errorf(
			"expected email to be %v, got: %v",
			req.Email,
			response.Account.Email,
		)
	}
}

func TestLoginSuccess(t *testing.T) {
	defer tearDown(t)

	password := "1234"
	acc, err := account.NewAccount("mihaib", "mihai@gmail.com", password)
	if err != nil {
		t.Error(err)
		return
	}

	cont := di.BuildContainer()
	err = cont.Invoke(func(accountRepo *account.AccountRepo) {
		accountRepo.Save(acc)
	})
	if err != nil {
		t.Error(err)
		return
	}

	r := mux.NewRouter()
	r.Handle("/api/v1/login", handler.LoginAccountHandler()).Methods(http.MethodPost)

	ts := httptest.NewServer(r)
	defer ts.Close()

	req := handler.LoginRequest{
		Email:    acc.Email,
		Password: password,
	}
	b, err := json.Marshal(req)
	if err != nil {
		t.Error(err)
		return
	}

	res, err := http.Post(ts.URL+"/api/v1/login", "application/json", bytes.NewBuffer(b))
	if err != nil {
		t.Error(err)
		return
	}

	response := handler.AccountResponse{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Error(err)
		return
	}

	if response.Success != true {
		t.Errorf("expected success to be true, got: %v", response.Success)
	}

	accountBalance := float64(10000)
	if response.Account.Balance != accountBalance {
		t.Errorf(
			"expected account balance to be %v, got: %v",
			accountBalance,
			response.Account.Balance,
		)
	}

	if response.Account.ApiToken == "" {
		t.Error("expected account token to not be empty")
	}

	if response.Account.Username != acc.Username {
		t.Errorf(
			"expected username to be %v, got: %v",
			acc.Username,
			response.Account.Username,
		)
	}

	if response.Account.Email != acc.Email {
		t.Errorf(
			"expected email to be %v, got: %v",
			acc.Email,
			response.Account.Email,
		)
	}
}

func TestFetchAccount(t *testing.T) {
	defer tearDown(t)

	acc, err := account.NewAccount("mihaib", "mihai@gmail.com", "1234")
	if err != nil {
		t.Error(err)
		return
	}

	cont := di.BuildContainer()
	err = cont.Invoke(func(accountRepo *account.AccountRepo) {
		accountRepo.Save(acc)
	})
	if err != nil {
		t.Error(err)
		return
	}

	r := mux.NewRouter()
	r.Handle("/api/v1/account", handler.AccountHandler()).Methods(http.MethodGet)

	ts := httptest.NewServer(r)
	defer ts.Close()

	client := http.Client{}
	req, err := http.NewRequest("GET", ts.URL+"/api/v1/account", nil)
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

	response := handler.AccountResponse{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Error(err)
		return
	}

	if response.Success != true {
		t.Errorf("expected success to be true, got: %v", response.Success)
	}

	accountBalance := float64(10000)
	if response.Account.Balance != accountBalance {
		t.Errorf(
			"expected account balance to be %v, got: %v",
			accountBalance,
			response.Account.Balance,
		)
	}

	if response.Account.ApiToken == "" {
		t.Error("expected account token to not be empty")
	}

	if response.Account.Username != acc.Username {
		t.Errorf(
			"expected username to be %v, got: %v",
			acc.Username,
			response.Account.Username,
		)
	}

	if response.Account.Email != acc.Email {
		t.Errorf(
			"expected email to be %v, got: %v",
			acc.Email,
			response.Account.Email,
		)
	}
}
