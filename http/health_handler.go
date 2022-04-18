package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

type HealthResponse struct {
	Server   bool `json:"server"`
	Database bool `json:"database"`
	Redis    bool `json:"redis"`
}

func healthHandler(cont *dig.Container) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := HealthResponse{}
		response.Server = true

		cont.Invoke(func(
			conn *gorm.DB,
			redisClient *redis.Client,
			logger *logrus.Logger) {

			var tables []string
			if err := conn.Table("information_schema.tables").Where("table_schema = ?", "public").Pluck("table_name", &tables).Error; err != nil {
				response.Database = false
			}

			if _, err := redisClient.Keys(context.Background(), "*").Result(); err != nil {
				response.Redis = false
			}

			sendResponse(w, logger, &response, http.StatusOK)
		})
	})
}

func sendResponse(
	w http.ResponseWriter,
	logger *logrus.Logger,
	resp interface{},
	code int) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	if code < 300 {
		logger.Info(resp)
	} else {
		logger.Error(resp)
	}

	b, _ := json.Marshal(resp)

	w.Write(b)
}

func serverError(w http.ResponseWriter, cont *dig.Container, err error) {
	if err != nil {
		cont.Invoke(func(logger *logrus.Logger) {
			resp := OrderResponse{
				Success: false,
				Error:   err.Error(),
			}
			sendResponse(w, logger, resp, http.StatusInternalServerError)
		})
	}
}
