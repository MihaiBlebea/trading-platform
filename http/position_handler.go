package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/MihaiBlebea/trading-platform/account"
	"github.com/MihaiBlebea/trading-platform/pos"
	"github.com/MihaiBlebea/trading-platform/symbols"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
)

type PositionsResponse struct {
	Success   bool           `json:"success"`
	Error     string         `json:"error,omitempty"`
	Equity    float64        `json:"equity"`
	Cash      float64        `json:"cash"`
	Total     float64        `json:"total"`
	Positions []pos.Position `json:"positions"`
}

func PositionsHandler(cont *dig.Container) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			serverError(w, cont, errors.New("could not find authorization header"))
			return
		}
		apiToken := strings.Split(header, "Bearer ")[1]

		err := cont.Invoke(func(
			accountRepo *account.AccountRepo,
			positionRepo *pos.PositionRepo,
			symbolService *symbols.Service,
			logger *logrus.Logger) {

			account, err := accountRepo.WithToken(apiToken)
			if err != nil {
				serverError(w, cont, err)
				return
			}

			positions, err := positionRepo.WithAccountId(account.ID)
			if err != nil {
				serverError(w, cont, err)
				return
			}

			resp := PositionsResponse{
				Success:   true,
				Cash:      account.Balance,
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
				resp.Equity += p.TotalValue
			}

			resp.Total = resp.Cash + resp.Equity

			sendResponse(w, logger, resp, http.StatusOK)
		})
		if err != nil {
			serverError(w, cont, err)
		}
	})
}
