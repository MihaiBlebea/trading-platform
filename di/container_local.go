package di

import (
	"fmt"
	"os"

	"github.com/MihaiBlebea/trading-platform/account"
	"github.com/MihaiBlebea/trading-platform/activity"
	"github.com/MihaiBlebea/trading-platform/order"
	"github.com/MihaiBlebea/trading-platform/pos"
	"github.com/MihaiBlebea/trading-platform/symbols"
	"github.com/MihaiBlebea/trading-platform/symbols/yahoofin"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// https://blog.drewolson.org/dependency-injection-in-go
func BuildContainer() *dig.Container {
	container := dig.New()

	if os.Getenv("APP_ENV") == "prod" {
		container.Provide(func() (*gorm.DB, error) {
			dsn := fmt.Sprintf(
				"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/London",
				os.Getenv("POSTGRES_HOST"),
				os.Getenv("POSTGRES_USER"),
				os.Getenv("POSTGRES_PASSWORD"),
				os.Getenv("POSTGRES_DB"),
				os.Getenv("POSTGRES_PORT"),
			)

			conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err != nil {
				return &gorm.DB{}, err
			}

			return conn, nil
		})

		container.Provide(func(client *yahoofin.ClientCache, repo *symbols.SymbolRepo) *symbols.Service {
			return symbols.NewService(client, repo)
		})
	} else {
		container.Provide(func() (*gorm.DB, error) {
			file := "file::memory:?cache=shared"
			// file := "gorm.db"
			conn, err := gorm.Open(sqlite.Open(file), &gorm.Config{})
			if err != nil {
				return &gorm.DB{}, err
			}

			return conn, nil
		})

		container.Provide(yahoofin.NewStubClient)

		container.Provide(func(client *yahoofin.ClientStub, repo *symbols.SymbolRepo) *symbols.Service {
			return symbols.NewService(client, repo)
		})
	}

	container.Provide(func() *logrus.Logger {
		logger := logrus.New()

		logger.SetFormatter(&logrus.JSONFormatter{})
		logger.SetOutput(os.Stdout)
		logger.SetLevel(logrus.InfoLevel)

		return logger
	})

	container.Provide(func() *redis.Client {
		return redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf(
				"%s:%s",
				os.Getenv("REDIS_HOST"),
				os.Getenv("REDIS_PORT"),
			),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       0, // use default DB
		})
	})

	container.Provide(account.NewAccountRepo, dig.As(new(activity.AccountRepo)))
	container.Provide(account.NewAccountRepo)

	container.Provide(order.NewOrderRepo, dig.As(new(activity.OrderRepo)))
	container.Provide(order.NewOrderRepo)

	container.Provide(pos.NewPositionRepo, dig.As(new(activity.PositionRepo)))
	container.Provide(pos.NewPositionRepo)

	container.Provide(symbols.NewSymbolRepo)

	container.Provide(yahoofin.NewClient)
	container.Provide(yahoofin.NewClientCache)

	container.Provide(activity.NewFiller)

	container.Provide(activity.NewOrderPlacer)

	container.Provide(activity.NewOrderCanceller)

	return container
}
