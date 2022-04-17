package http

import (
	"errors"
	"fmt"
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

		positionRepo, err := di.GetPositionRepo()
		if err != nil {
			serverError(w, err)
			return
		}

		positions, err := positionRepo.WithAccountId(account.ID)
		if err != nil {
			serverError(w, err)
			return
		}

		symbolService, err := di.GetSymbolService()
		if err != nil {
			serverError(w, err)
			return
		}

		resp := PositionsResponse{
			Success: true,
		}

		for _, p := range positions {
			s, err := symbolService.GetSymbol(p.Symbol)
			if err != nil {
				fmt.Println(err)
				continue
			}
			p.TotalValue = s.MarketPrice * float64(p.Quantity)
			p.CalculateAverageBoughtPrice()

			resp.Positions = append(resp.Positions, p)
		}

		sendResponse(w, resp, http.StatusOK)
	})
}
