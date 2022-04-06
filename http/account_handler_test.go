package http_test

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/MihaiBlebea/trading-platform/account"
// 	"github.com/MihaiBlebea/trading-platform/di"
// 	handler "github.com/MihaiBlebea/trading-platform/http"
// 	"github.com/gorilla/mux"
// )

// func getAccountToken() string {
// 	container, _ := di.NewContainer()
// 	accountRepo := container.GetAccountRepo()
// 	account, _ := accountRepo.Save(account.NewAccount())

// 	return account.ApiToken
// }

// func TestCreateAccountSuccess(t *testing.T) {
// 	r := mux.NewRouter()
// 	r.Handle("/api/v1/account", handler.CreateAccountHandler()).Methods(http.MethodPost)

// 	ts := httptest.NewServer(r)
// 	defer ts.Close()

// 	res, err := http.Post(ts.URL+"/api/v1/account", "application/json", nil)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}

// 	response := handler.AccountResponse{}
// 	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
// 		t.Error(err)
// 		return
// 	}

// 	if response.Success != true {
// 		t.Errorf("expected success to be true, got: %v", response.Success)
// 	}

// 	accountBalance := float32(10000)
// 	if response.Account.Balance != accountBalance {
// 		t.Errorf(
// 			"expected account balance to be %v, got: %v",
// 			accountBalance,
// 			response.Account.Balance,
// 		)
// 	}

// 	if response.Account.ApiToken == "" {
// 		t.Error("expected account token to not be empty")
// 	}
// }

// func TestAccountHandlerSuccess(t *testing.T) {
// 	r := mux.NewRouter()
// 	r.Handle("/api/v1/account", handler.AccountHandler()).Methods(http.MethodGet)

// 	ts := httptest.NewServer(r)
// 	defer ts.Close()

// 	client := http.Client{}
// 	req, err := http.NewRequest("GET", ts.URL+"/api/v1/account", nil)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}

// 	req.Header.Set(
// 		"Authorization",
// 		fmt.Sprintf("Bearer %s", getAccountToken()),
// 	)

// 	res, err := client.Do(req)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}

// 	response := handler.AccountResponse{}
// 	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
// 		t.Error(err)
// 		return
// 	}

// 	if response.Success != true {
// 		t.Errorf("expected success to be true, got: %v", response.Success)
// 	}

// 	accountBalance := float32(10000)
// 	if response.Account.Balance != accountBalance {
// 		t.Errorf(
// 			"expected account balance to be %v, got: %v",
// 			accountBalance,
// 			response.Account.Balance,
// 		)
// 	}

// 	if response.Account.ApiToken == "" {
// 		t.Error("expected account token to not be empty")
// 	}
// }
