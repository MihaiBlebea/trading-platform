package http

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/MihaiBlebea/trading-platform/account"
	"github.com/MihaiBlebea/trading-platform/di"
	"github.com/MihaiBlebea/trading-platform/pos"
	"github.com/MihaiBlebea/trading-platform/symbols"
	"github.com/sirupsen/logrus"
)

type PositionsResponse struct {
	Success   bool           `json:"success"`
	Error     string         `json:"error,omitempty"`
	Positions []pos.Position `json:"positions"`
}

func PositionsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			serverError(w, errors.New("could not find authorization header"))
			return
		}
		apiToken := strings.Split(header, "Bearer ")[1]

		err := di.BuildContainer().Invoke(func(
			accountRepo *account.AccountRepo,
			positionRepo *pos.PositionRepo,
			symbolService *symbols.Service,
			logger *logrus.Logger) {

			account, err := accountRepo.WithToken(apiToken)
			if err != nil {
				serverError(w, err)
				return
			}

			positions, err := positionRepo.WithAccountId(account.ID)
			if err != nil {
				serverError(w, err)
				return
			}

			fmt.Println(positions)

			resp := PositionsResponse{
				Success:   true,
				Positions: []pos.Position{},
			}

			for _, p := range positions {
				s, err := symbolService.GetSymbol(p.Symbol)
				if err != nil {
					logger.Error(err)
					continue
				}
				p.TotalValue = s.MarketPrice * float64(p.Quantity)
				p.CalculateAverageBoughtPrice()

				resp.Positions = append(resp.Positions, p)
			}

			sendResponse(w, resp, http.StatusOK)
		})
		if err != nil {
			serverError(w, err)
		}
	})
}
