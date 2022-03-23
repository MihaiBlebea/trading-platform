package http

import (
	"net/http"

	yahooapi "github.com/MihaiBlebea/yahoo-api-go"
	"github.com/gorilla/mux"
)

func fundamentalsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		symbol := vars["symbol"]

		client := yahooapi.NewClient(false)
		res, err := client.BalanceSheetHistoryQuarterly(symbol)
		if err != nil {
			serverError(w, err)
			return
		}

		var response struct {
			Data []yahooapi.BalanceSheet
		}
		response.Data = res

		sendResponse(w, &response, http.StatusOK)
	})
}
