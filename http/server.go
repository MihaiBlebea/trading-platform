package http

import (
	"fmt"

	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
)

const prefix = "/api/v1/"

func New(container *dig.Container) {

	r := mux.NewRouter()

	api := r.PathPrefix(prefix).Subrouter()

	// Handle api calls
	api.Handle("/health-check", healthHandler(container)).
		Methods(http.MethodGet)

	// Account endpoints
	api.Handle("/login", LoginAccountHandler(container)).
		Methods(http.MethodPost)

	api.Handle("/register", RegisterAccountHandler(container)).
		Methods(http.MethodPost)

	api.Handle("/account", AccountHandler(container)).
		Methods(http.MethodGet)

	// Order endpoints
	api.Handle("/order", PlaceOrderHandler(container)).
		Methods(http.MethodPost)

	api.Handle("/order/cancel", CancelOrderHandler(container)).
		Methods(http.MethodPut)

	api.Handle("/orders", OrdersHandler(container)).
		Methods(http.MethodGet)

	// Position endpoints
	api.Handle("/positions", PositionsHandler(container)).
		Methods(http.MethodGet)

	// Symbols endpoints
	api.Handle("/symbol", symbolHandler(container)).
		Methods(http.MethodGet)

	api.Handle("/symbols", symbolsHandler(container)).
		Methods(http.MethodGet)

	api.Handle("/chart", chartHandler(container)).
		Methods(http.MethodGet)

	container.Invoke(func(logger *logrus.Logger) {
		r.Use(loggerMiddleware(logger))

		logger.Info(fmt.Sprintf("Started server on port %s", os.Getenv("HTTP_PORT")))

		srv := &http.Server{
			Handler:      cors.AllowAll().Handler(r),
			Addr:         fmt.Sprintf("0.0.0.0:%s", os.Getenv("HTTP_PORT")),
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}

		logger.Fatal(srv.ListenAndServe())
	})
}
