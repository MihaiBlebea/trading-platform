package http

import (
	"net/http"

	"github.com/MihaiBlebea/trading-platform/account"
)

type CreateAccountResponse struct {
	Success bool             `json:"success"`
	Error   string           `json:"error,omitempty"`
	Account *account.Account `json:"account,omitempty"`
}

func createAccountHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		repo, err := account.NewAccountRepo()
		if err != nil {
			resp := CreateAccountResponse{
				Success: false,
				Error:   err.Error(),
			}
			sendResponse(w, resp, http.StatusInternalServerError)
			return
		}

		account, err := repo.Save(account.NewAccount())
		if err != nil {
			resp := CreateAccountResponse{
				Success: false,
				Error:   err.Error(),
			}
			sendResponse(w, resp, http.StatusInternalServerError)
			return
		}

		resp := CreateAccountResponse{
			Success: true,
			Account: account,
		}
		sendResponse(w, resp, http.StatusOK)
	})
}
