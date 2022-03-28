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

		di, err := di.NewContainer()
		if err != nil {
			resp := AccountResponse{
				Success: false,
				Error:   err.Error(),
			}
			sendResponse(w, resp, http.StatusInternalServerError)
			return
		}

		accountRepo := di.GetAccountRepo()

		account, err := accountRepo.Save(account.NewAccount())
		if err != nil {
			resp := AccountResponse{
				Success: false,
				Error:   err.Error(),
			}
			sendResponse(w, resp, http.StatusInternalServerError)
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

		di, err := di.NewContainer()
		if err != nil {
			resp := AccountResponse{
				Success: false,
				Error:   err.Error(),
			}
			sendResponse(w, resp, http.StatusInternalServerError)
			return
		}

		accountRepo := di.GetAccountRepo()

		account, err := accountRepo.WithToken(apiToken)
		if err != nil {
			resp := AccountResponse{
				Success: false,
				Error:   err.Error(),
			}
			sendResponse(w, resp, http.StatusInternalServerError)
			return
		}

		resp := AccountResponse{
			Success: true,
			Account: account,
		}
		sendResponse(w, resp, http.StatusOK)
	})
}
