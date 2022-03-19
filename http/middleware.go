package http

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

func loggerMiddleware(logger *logrus.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info(fmt.Sprintf("Incoming %s request %s", r.Method, r.URL.Path))
			next.ServeHTTP(w, r)
		})
	}
}
