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

	// Data endpoints
	api.Handle("/data/historic", historicDataHandler()).
		Methods(http.MethodGet)

	// Account endpoints
	api.Handle("/account", createAccountHandler()).
		Methods(http.MethodPost)

	api.Handle("/account", accountHandler()).
		Methods(http.MethodGet)

	// Order endpoints
	api.Handle("/order", placeOrderHandler()).
		Methods(http.MethodPost)

	api.Handle("/orders", ordersHandler()).
		Methods(http.MethodGet)

	// Position endpoints
	api.Handle("/positions", positionsHandler()).
		Methods(http.MethodGet)

	// Fundamentals endpoints
	api.Handle("/fundamental/{symbol}", fundamentalsHandler()).
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
