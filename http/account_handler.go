package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/MihaiBlebea/trading-platform/account"
	"github.com/MihaiBlebea/trading-platform/di"
)

type AccountResponse struct {
	Success bool             `json:"success"`
	Error   string           `json:"error,omitempty"`
	Account *account.Account `json:"account,omitempty"`
}

func CreateAccountHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		di := di.NewContainer()

		accountRepo, err := di.GetAccountRepo()
		if err != nil {
			serverError(w, err)
			return
		}

		account, err := accountRepo.Save(account.NewAccount())
		if err != nil {
			serverError(w, err)
			return
		}

		resp := AccountResponse{
			Success: true,
			Account: account,
		}
		sendResponse(w, resp, http.StatusOK)
	})
}

func AccountHandler() http.Handler {
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

		account, err := accountRepo.WithToken(apiToken)
		if err != nil {
			serverError(w, err)
			return
		}

		resp := AccountResponse{
			Success: true,
			Account: account,
		}
		sendResponse(w, resp, http.StatusOK)
	})
}
