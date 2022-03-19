package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/MihaiBlebea/trading-platform/account"
)

type AccountResponse struct {
	Success bool             `json:"success"`
	Error   string           `json:"error,omitempty"`
	Account *account.Account `json:"account,omitempty"`
}

func createAccountHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		repo, err := account.NewAccountRepo()
		if err != nil {
			resp := AccountResponse{
				Success: false,
				Error:   err.Error(),
			}
			sendResponse(w, resp, http.StatusInternalServerError)
			return
		}

		account, err := repo.Save(account.NewAccount())
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

func accountHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			serverError(w, errors.New("could not find authorization header"))
			return
		}
		apiToken := strings.Split(header, "Bearer ")[1]

		repo, err := account.NewAccountRepo()
		if err != nil {
			resp := AccountResponse{
				Success: false,
				Error:   err.Error(),
			}
			sendResponse(w, resp, http.StatusInternalServerError)
			return
		}

		account, err := repo.WithToken(apiToken)
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
