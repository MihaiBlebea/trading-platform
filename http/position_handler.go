package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/MihaiBlebea/trading-platform/di"
	"github.com/MihaiBlebea/trading-platform/pos"
)

type PositionsResponse struct {
	Success   bool           `json:"success"`
	Error     string         `json:"error,omitempty"`
	Positions []pos.Position `json:"positions"`
}

func positionsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			serverError(w, errors.New("could not find authorization header"))
			return
		}
		apiToken := strings.Split(header, "Bearer ")[1]

		di, err := di.NewContainer()
		if err != nil {
			serverError(w, err)
			return
		}

		accountRepo := di.GetAccountRepo()

		account, err := accountRepo.WithToken(apiToken)
		if err != nil {
			serverError(w, err)
			return
		}

		positionRepo := di.GetPositionRepo()

		positions, err := positionRepo.WithAccountId(account.ID)
		if err != nil {
			serverError(w, err)
			return
		}

		resp := PositionsResponse{
			Success:   true,
			Positions: positions,
		}
		sendResponse(w, resp, http.StatusOK)
	})
}
