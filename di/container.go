package di

import (
	"fmt"
	"log"
	"os"

	"github.com/MihaiBlebea/trading-platform/account"
	"github.com/MihaiBlebea/trading-platform/activity"
	"github.com/MihaiBlebea/trading-platform/order"
	"github.com/MihaiBlebea/trading-platform/pos"
	"github.com/MihaiBlebea/trading-platform/symbols"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var instance *Container

func init() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatal(err)
	}

	container, err := NewContainer()
	if err != nil {
		log.Fatal(err)
	}

	instance = container
}

type Container struct {
	conn           *gorm.DB
	logger         *logrus.Logger
	accountRepo    *account.AccountRepo
	orderRepo      *order.OrderRepo
	positionRepo   *pos.PositionRepo
	symbolRepo     *symbols.SymbolRepo
	orderFiller    *activity.Filler
	orderPlacer    *activity.OrderPlacer
	orderCanceller *activity.OrderCanceller
}

func NewContainer() (*Container, error) {
	if instance != nil {
		return instance, nil
	}

	logger := logrus.New()

	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)

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
		return &Container{}, err
	}

	accountRepo, err := account.NewAccountRepo(conn)
	if err != nil {
		return &Container{}, err
	}

	orderRepo, err := order.NewOrderRepo(conn)
	if err != nil {
		return &Container{}, err
	}

	positionRepo, err := pos.NewPositionRepo(conn)
	if err != nil {
		return &Container{}, err
	}

	symbolRepo, err := symbols.NewSymbolRepo(conn)
	if err != nil {
		return &Container{}, err
	}

	filler := activity.NewFiller(accountRepo, orderRepo, positionRepo, logger)

	orderPlacer := activity.NewOrderPlacer(accountRepo, orderRepo, positionRepo)

	orderCanceller := activity.NewOrderCanceller(accountRepo, orderRepo)

	return &Container{
		conn:           conn,
		logger:         logger,
		accountRepo:    accountRepo,
		orderRepo:      orderRepo,
		positionRepo:   positionRepo,
		symbolRepo:     symbolRepo,
		orderFiller:    filler,
		orderPlacer:    orderPlacer,
		orderCanceller: orderCanceller,
	}, nil
}

func (c *Container) GetAccountRepo() *account.AccountRepo {
	return c.accountRepo
}

func (c *Container) GetOrderRepo() *order.OrderRepo {
	return c.orderRepo
}

func (c *Container) GetPositionRepo() *pos.PositionRepo {
	return c.positionRepo
}

func (c *Container) GetSymbolRepo() *symbols.SymbolRepo {
	return c.symbolRepo
}

func (c *Container) GetOrderFiller() *activity.Filler {
	return c.orderFiller
}

func (c *Container) GetOrderPlacer() *activity.OrderPlacer {
	return c.orderPlacer
}

func (c *Container) GetOrderCanceller() *activity.OrderCanceller {
	return c.orderCanceller
}

func (c *Container) GetLogger() *logrus.Logger {
	return c.logger
}
