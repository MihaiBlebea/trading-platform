package order

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type OrderType string

const (
	TypeLimit      OrderType = "limit"
	TypeStopLoss   OrderType = "stop-loss"
	TypeTakeProfit OrderType = "take-profit"
)

type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusFilled    OrderStatus = "filled"
	StatusCancelled OrderStatus = "cancelled"
)

type OrderDirection string

const (
	DirectionBuy  OrderDirection = "buy"
	DirectionSell OrderDirection = "sell"
)

type Order struct {
	ID          int            `json:"id"`
	AccountID   int            `json:"-"`
	Type        OrderType      `json:"type"`
	Status      OrderStatus    `json:"status"`
	Direction   OrderDirection `json:"direction"`
	Amount      float32        `json:"amount"`
	FillPrice   float32        `json:"fill_price"`
	Symbol      string         `json:"symbol"`
	Quantity    int            `json:"quantity"`
	FilledAt    time.Time      `json:"filled_at,omitempty"`
	CancelledAt time.Time      `json:"cancelled_at,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"-"`
}

type OrderRepo struct {
	conn *gorm.DB
}

func NewOrder(accountId int, orderType, direction string, amount float32, symbol string) *Order {
	return &Order{
		AccountID: accountId,
		Type:      OrderType(orderType),
		Status:    StatusPending,
		Direction: OrderDirection(direction),
		Amount:    amount,
		Symbol:    strings.ToUpper(symbol),
	}
}

func NewOrderRepo() (*OrderRepo, error) {
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
		return &OrderRepo{}, err
	}

	if err := conn.AutoMigrate(&Order{}); err != nil {
		return &OrderRepo{}, err
	}

	return &OrderRepo{conn: conn}, nil
}

func (or *OrderRepo) Save(order *Order) (*Order, error) {
	resp := or.conn.Create(order)

	if resp.Error != nil {
		return &Order{}, resp.Error
	}
	return order, nil
}

func (or *OrderRepo) Update(order *Order) error {
	return or.conn.Save(order).Error
}

func (or *OrderRepo) WithPendingStatus() ([]Order, error) {
	orders := []Order{}
	err := or.conn.Where("status = 'pending'").Find(&orders).Error
	if err != nil {
		return []Order{}, err
	}

	return orders, err
}

func (or *OrderRepo) WithAccountId(accountId int) ([]Order, error) {
	orders := []Order{}
	err := or.conn.Where("account_id = ?", accountId).Find(&orders).Error
	if err != nil {
		return []Order{}, err
	}

	return orders, err
}

func (o *Order) FillOrder(price float32) {
	o.FillPrice = price
	o.FilledAt = time.Now()
	o.Quantity = int(o.Amount / price)
	o.Status = StatusFilled
}
