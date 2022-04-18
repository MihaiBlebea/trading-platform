package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/MihaiBlebea/trading-platform/account"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
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

func RegisterAccountHandler(cont *dig.Container) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req RegisterRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			serverError(w, cont, err)
			return
		}

		if req.Username == "" {
			serverError(w, cont, errors.New("invalid username"))
			return
		}

		if req.Password == "" {
			serverError(w, cont, errors.New("invalid password"))
			return
		}

		if req.Email == "" {
			serverError(w, cont, errors.New("invalid email"))
			return
		}

		err = cont.Invoke(func(accRepo *account.AccountRepo, logger *logrus.Logger) {
			account, err := account.NewAccount(req.Username, req.Email, req.Password)
			if err != nil {
				serverError(w, cont, err)
				return
			}

			if _, err := accRepo.Save(account); err != nil {
				serverError(w, cont, err)
				return
			}

			resp := AccountResponse{
				Success: true,
				Account: account,
			}
			sendResponse(w, logger, resp, http.StatusOK)
		})

		if err != nil {
			serverError(w, cont, err)
			return
		}
	})
}

func LoginAccountHandler(cont *dig.Container) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			serverError(w, cont, err)
			return
		}

		err = cont.Invoke(func(accountRepo *account.AccountRepo, logger *logrus.Logger) {
			account, err := accountRepo.WithEmail(req.Email)
			if err != nil {
				serverError(w, cont, errors.New("invalid credentials"))
				return
			}

			if !account.CheckPasswordHash(req.Password) {
				serverError(w, cont, errors.New("invalid credentials"))
				return
			}

			resp := AccountResponse{
				Success: true,
				Account: account,
			}
			sendResponse(w, logger, resp, http.StatusOK)
		})

		if err != nil {
			serverError(w, cont, err)
			return
		}
	})
}

func AccountHandler(cont *dig.Container) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			serverError(w, cont, errors.New("could not find authorization header"))
			return
		}
		apiToken := strings.Split(header, "Bearer ")[1]

		err := cont.Invoke(func(accountRepo *account.AccountRepo, logger *logrus.Logger) {
			account, err := accountRepo.WithToken(apiToken)
			if err != nil {
				serverError(w, cont, err)
				return
			}

			resp := AccountResponse{
				Success: true,
				Account: account,
			}
			sendResponse(w, logger, resp, http.StatusOK)
		})
		if err != nil {
			serverError(w, cont, err)
			return
		}
	})
}
