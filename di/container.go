package di

import (
	"fmt"
	"os"

	"github.com/MihaiBlebea/trading-platform/account"
	"github.com/MihaiBlebea/trading-platform/activity"
	"github.com/MihaiBlebea/trading-platform/market"
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
	if err := godotenv.Load(); err != nil {
		fmt.Println(err)
	}

	container := NewContainer()

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
	marketStatus   *market.MarketStatus
}

func NewContainer() *Container {
	if instance != nil {
		return instance
	}

	logger := logrus.New()

	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)

	return &Container{
		logger: logger,
	}
}

func (c *Container) GetDatabaseConn() (*gorm.DB, error) {
	if c.conn != nil {
		return c.conn, nil
	}
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

	c.conn = conn

	return conn, nil
}

func (c *Container) GetAccountRepo() (*account.AccountRepo, error) {
	if c.accountRepo != nil {
		return c.accountRepo, nil
	}

	conn, err := c.GetDatabaseConn()
	if err != nil {
		return &account.AccountRepo{}, err
	}

	accountRepo, err := account.NewAccountRepo(conn)
	if err != nil {
		return &account.AccountRepo{}, err
	}

	c.accountRepo = accountRepo

	return c.accountRepo, nil
}

func (c *Container) GetOrderRepo() (*order.OrderRepo, error) {
	if c.orderRepo != nil {
		return c.orderRepo, nil
	}

	conn, err := c.GetDatabaseConn()
	if err != nil {
		return &order.OrderRepo{}, err
	}

	orderRepo, err := order.NewOrderRepo(conn)
	if err != nil {
		return &order.OrderRepo{}, err
	}

	c.orderRepo = orderRepo

	return c.orderRepo, nil
}

func (c *Container) GetPositionRepo() (*pos.PositionRepo, error) {
	if c.positionRepo != nil {
		return c.positionRepo, nil
	}

	conn, err := c.GetDatabaseConn()
	if err != nil {
		return &pos.PositionRepo{}, err
	}

	positionRepo, err := pos.NewPositionRepo(conn)
	if err != nil {
		return &pos.PositionRepo{}, err
	}

	c.positionRepo = positionRepo

	return c.positionRepo, nil
}

func (c *Container) GetSymbolRepo() (*symbols.SymbolRepo, error) {
	if c.symbolRepo != nil {
		return c.symbolRepo, nil
	}

	conn, err := c.GetDatabaseConn()
	if err != nil {
		return &symbols.SymbolRepo{}, err
	}

	symbolRepo, err := symbols.NewSymbolRepo(conn)
	if err != nil {
		return &symbols.SymbolRepo{}, err
	}

	c.symbolRepo = symbolRepo

	return c.symbolRepo, nil
}

func (c *Container) GetOrderFiller() (*activity.Filler, error) {
	if c.orderFiller != nil {
		return c.orderFiller, nil
	}

	accountRepo, err := c.GetAccountRepo()
	if err != nil {
		return &activity.Filler{}, err
	}

	orderRepo, err := c.GetOrderRepo()
	if err != nil {
		return &activity.Filler{}, err
	}

	positionRepo, err := c.GetPositionRepo()
	if err != nil {
		return &activity.Filler{}, err
	}

	marketStatus, err := c.GetMarketStatus()
	if err != nil {
		return &activity.Filler{}, err
	}

	orderFiller := activity.NewFiller(
		accountRepo,
		orderRepo,
		positionRepo,
		marketStatus,
		c.logger)

	c.orderFiller = orderFiller

	return c.orderFiller, nil
}

func (c *Container) GetOrderPlacer() (*activity.OrderPlacer, error) {
	if c.orderPlacer != nil {
		return c.orderPlacer, nil
	}

	accountRepo, err := c.GetAccountRepo()
	if err != nil {
		return &activity.OrderPlacer{}, err
	}

	orderRepo, err := c.GetOrderRepo()
	if err != nil {
		return &activity.OrderPlacer{}, err
	}

	positionRepo, err := c.GetPositionRepo()
	if err != nil {
		return &activity.OrderPlacer{}, err
	}

	orderPlacer := activity.NewOrderPlacer(accountRepo, orderRepo, positionRepo)

	c.orderPlacer = orderPlacer

	return c.orderPlacer, nil
}

func (c *Container) GetOrderCanceller() (*activity.OrderCanceller, error) {
	if c.orderCanceller != nil {
		return c.orderCanceller, nil
	}

	accountRepo, err := c.GetAccountRepo()
	if err != nil {
		return &activity.OrderCanceller{}, err
	}

	orderRepo, err := c.GetOrderRepo()
	if err != nil {
		return &activity.OrderCanceller{}, err
	}

	orderCanceller := activity.NewOrderCanceller(accountRepo, orderRepo)

	c.orderCanceller = orderCanceller

	return c.orderCanceller, nil
}

func (c *Container) GetLogger() *logrus.Logger {
	return c.logger
}

func (c *Container) GetMarketStatus() (*market.MarketStatus, error) {
	if c.marketStatus != nil {
		return c.marketStatus, nil
	}

	marketStatus := market.New(nil)

	c.marketStatus = marketStatus

	return c.marketStatus, nil
}
