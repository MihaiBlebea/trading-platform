package http

import (
	"fmt"
	"log"

	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

const prefix = "/api/v1/"

func New(logger *logrus.Logger) {

	r := mux.NewRouter()

	api := r.PathPrefix(prefix).Subrouter()

	// Handle api calls
	api.Handle("/health-check", healthHandler()).
		Methods(http.MethodGet)

	// Account endpoints
	api.Handle("/login", LoginAccountHandler()).
		Methods(http.MethodPost)

	api.Handle("/register", RegisterAccountHandler()).
		Methods(http.MethodPost)

	api.Handle("/account", AccountHandler()).
		Methods(http.MethodGet)

	// Order endpoints
	api.Handle("/order", PlaceOrderHandler()).
		Methods(http.MethodPost)

	api.Handle("/order/cancel", CancelOrderHandler()).
		Methods(http.MethodPut)

	api.Handle("/orders", OrdersHandler()).
		Methods(http.MethodGet)

	// Position endpoints
	api.Handle("/positions", positionsHandler()).
		Methods(http.MethodGet)

	// Symbols endpoints
	api.Handle("/symbol", symbolHandler()).
		Methods(http.MethodGet)

	api.Handle("/symbols", symbolsHandler()).
		Methods(http.MethodGet)

	api.Handle("/chart", chartHandler()).
		Methods(http.MethodGet)

	r.Use(loggerMiddleware(logger))

	srv := &http.Server{
		Handler:      cors.AllowAll().Handler(r),
		Addr:         fmt.Sprintf("0.0.0.0:%s", os.Getenv("HTTP_PORT")),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Info(fmt.Sprintf("Started server on port %s", os.Getenv("HTTP_PORT")))

	log.Fatal(srv.ListenAndServe())
}
