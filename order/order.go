package order

import (
	"strings"
	"time"
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
	ID              int            `json:"id"`
	AccountID       int            `json:"-"`
	Type            OrderType      `json:"type"`
	Status          OrderStatus    `json:"status"`
	Direction       OrderDirection `json:"direction"`
	Amount          float32        `json:"amount"`
	FillPrice       float32        `json:"fill_price"`
	AmountAfterFill float32        `json:"amount_after_fill"`
	Symbol          string         `json:"symbol"`
	Quantity        int            `json:"quantity"`
	FilledAt        *time.Time     `json:"filled_at,omitempty"`
	CancelledAt     *time.Time     `json:"cancelled_at,omitempty"`
	CreatedAt       *time.Time     `json:"created_at"`
	UpdatedAt       *time.Time     `json:"-"`
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

func (o *Order) FillOrder(price float32) {
	now := time.Now()
	o.FillPrice = price
	o.FilledAt = &now
	o.Quantity = int(o.Amount / price)
	o.AmountAfterFill = o.FillPrice * float32(o.Quantity)
	if o.Direction == DirectionBuy {
		o.AmountAfterFill = -o.AmountAfterFill
	}

	o.Status = StatusFilled
}

func (o *Order) GetTotalFillAmount() float32 {
	return o.AmountAfterFill
}

func (o *Order) GetAmount() float32 {
	return o.Amount
}

func (o *Order) GetDirectionString() string {
	return string(o.Direction)
}
