package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/MihaiBlebea/trading-platform/di"
)

type HealthResponse struct {
	Server   bool `json:"server"`
	Database bool `json:"database"`
	Redis    bool `json:"redis"`
}

func healthHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := HealthResponse{}
		response.Server = true

		conn, err := di.NewContainer().GetDatabaseConn()
		if err == nil {
			response.Database = true
		}

		var tables []string
		if err := conn.Table("information_schema.tables").Where("table_schema = ?", "public").Pluck("table_name", &tables).Error; err != nil {
			response.Database = false
		}

		redisClient, err := di.NewContainer().GetRedisClient()
		if err == nil {
			response.Redis = true
		}

		if _, err := redisClient.Keys(context.Background(), "*").Result(); err != nil {
			response.Redis = false
		}

		sendResponse(w, &response, http.StatusOK)
	})
}

func sendResponse(w http.ResponseWriter, resp interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	b, _ := json.Marshal(resp)

	w.Write(b)
}
