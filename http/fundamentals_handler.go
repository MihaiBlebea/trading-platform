package http

import (
	"errors"
	"net/http"

	yahooapi "github.com/MihaiBlebea/yahoo-api-go"
	"github.com/gorilla/mux"
)

var (
	ModuleBalanceSheetHistoryQuarterly    = "balance-sheet-history-quarterly"
	ModuleEarningsHistory                 = "earnings-history"
	ModuleIncomeStatementHistoryQuarterly = "income-statement-history-quarterly"
)

func fundamentalsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		symbol := vars["symbol"]
		module := vars["module"]

		client := yahooapi.NewClient(false)
		if module == ModuleBalanceSheetHistoryQuarterly {
			balanceSheetHistoryQuarterly(client, w, symbol)
			return
		}

		if module == ModuleEarningsHistory {
			earningsHistory(client, w, symbol)
			return
		}

		if module == ModuleIncomeStatementHistoryQuarterly {
			incomeStatementHistoryQuarterly(client, w, symbol)
			return
		}

		serverError(w, errors.New("could not find module"))
	})
}

func balanceSheetHistoryQuarterly(client *yahooapi.Client, w http.ResponseWriter, symbol string) {
	res, err := client.BalanceSheetHistoryQuarterly(symbol)
	if err != nil {
		serverError(w, err)
		return
	}

	response := struct {
		Data    []yahooapi.BalanceSheet `json:"data"`
		Success bool                    `json:"success"`
	}{
		Data:    res,
		Success: true,
	}

	sendResponse(w, &response, http.StatusOK)
}

func earningsHistory(client *yahooapi.Client, w http.ResponseWriter, symbol string) {
	res, err := client.EarningsHistory(symbol)
	if err != nil {
		serverError(w, err)
		return
	}

	response := struct {
		Data    []yahooapi.History `json:"data"`
		Success bool               `json:"success"`
	}{
		Data:    res,
		Success: true,
	}

	sendResponse(w, &response, http.StatusOK)
}

func incomeStatementHistoryQuarterly(client *yahooapi.Client, w http.ResponseWriter, symbol string) {
	res, err := client.IncomeStatementHistoryQuarterly(symbol)
	if err != nil {
		serverError(w, err)
		return
	}

	response := struct {
		Data    []yahooapi.IncomeStatement `json:"data"`
		Success bool                       `json:"success"`
	}{
		Data:    res,
		Success: true,
	}

	sendResponse(w, &response, http.StatusOK)
}
