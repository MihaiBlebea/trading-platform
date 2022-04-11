package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/MihaiBlebea/trading-platform/account"
	"github.com/MihaiBlebea/trading-platform/di"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AccountResponse struct {
	Success bool             `json:"success"`
	Error   string           `json:"error,omitempty"`
	Account *account.Account `json:"account,omitempty"`
}

func RegisterAccountHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req RegisterRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			serverError(w, err)
			return
		}

		if req.Username == "" {
			serverError(w, errors.New("invalid username"))
			return
		}

		if req.Password == "" {
			serverError(w, errors.New("invalid password"))
			return
		}

		if req.Email == "" {
			serverError(w, errors.New("invalid email"))
			return
		}

		di := di.NewContainer()

		accountRepo, err := di.GetAccountRepo()
		if err != nil {
			serverError(w, err)
			return
		}

		account, err := account.NewAccount(req.Username, req.Email, req.Password)
		if err != nil {
			serverError(w, err)
			return
		}

		if _, err := accountRepo.Save(account); err != nil {
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

func LoginAccountHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			serverError(w, err)
			return
		}

		di := di.NewContainer()

		accountRepo, err := di.GetAccountRepo()
		if err != nil {
			serverError(w, err)
			return
		}

		account, err := accountRepo.WithEmail(req.Email)
		if err != nil {
			serverError(w, errors.New("invalid credentials"))
			return
		}

		if account.CheckPasswordHash(req.Password) == false {
			serverError(w, errors.New("invalid credentials"))
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
