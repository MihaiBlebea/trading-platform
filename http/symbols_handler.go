package http

import (
	"errors"
	"net/http"

	"github.com/MihaiBlebea/trading-platform/symbols"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
)

type SymbolResponse struct {
	Success bool            `json:"success"`
	Error   string          `json:"error,omitempty"`
	Symbol  *symbols.Symbol `json:"symbol,omitempty"`
}

type SymbolsResponse struct {
	Success bool             `json:"success"`
	Error   string           `json:"error,omitempty"`
	Symbols []symbols.Symbol `json:"symbols,omitempty"`
}

type ChartResponse struct {
	Success bool            `json:"success"`
	Error   string          `json:"error,omitempty"`
	Chart   []symbols.Chart `json:"chart,omitempty"`
}

func symbolHandler(cont *dig.Container) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		symbol := query.Get("symbol")
		if symbol == "" {
			serverError(w, cont, errors.New("invalid symbol"))
			return
		}

		err := cont.Invoke(func(symbolService *symbols.Service, logger *logrus.Logger) {
			s, err := symbolService.GetSymbol(symbol)
			if err != nil {
				serverError(w, cont, err)
				return
			}

			resp := SymbolResponse{
				Success: true,
				Symbol:  s,
			}
			sendResponse(w, logger, resp, http.StatusOK)
		})
		if err != nil {
			serverError(w, cont, err)
			return
		}
	})
}

func symbolsHandler(cont *dig.Container) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		symbol := query.Get("search")
		if symbol == "" {
			serverError(w, cont, errors.New("invalid symbol"))
			return
		}

		err := cont.Invoke(func(symbolRepo *symbols.SymbolRepo, logger *logrus.Logger) {
			symbols, err := symbolRepo.LikeSymbol(symbol)
			if err != nil {
				serverError(w, cont, err)
				return
			}

			resp := SymbolsResponse{
				Success: true,
				Symbols: symbols,
			}
			sendResponse(w, logger, resp, http.StatusOK)
		})
		if err != nil {
			serverError(w, cont, err)
			return
		}
	})
}

func chartHandler(cont *dig.Container) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		symbol := query.Get("symbol")
		if symbol == "" {
			serverError(w, cont, errors.New("invalid symbol"))
			return
		}

		err := cont.Invoke(func(symbolService *symbols.Service, logger *logrus.Logger) {
			charts, err := symbolService.GetChart(symbol)
			if err != nil {
				serverError(w, cont, err)
				return
			}

			resp := ChartResponse{
				Success: true,
				Chart:   charts,
			}
			sendResponse(w, logger, resp, http.StatusOK)
		})
		if err != nil {
			serverError(w, cont, err)
			return
		}
	})
}
