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
	ParentOrderID   int            `json:"parent_order_id,omitempty"`
	AccountID       int            `json:"-"`
	Type            OrderType      `json:"type"`
	Status          OrderStatus    `json:"status"`
	Direction       OrderDirection `json:"direction"`
	Amount          float64        `json:"amount"`
	FillPrice       float64        `json:"fill_price"`
	AmountAfterFill float64        `json:"amount_after_fill"`
	Symbol          string         `json:"symbol"`
	Quantity        int            `json:"quantity"`
	FilledAt        *time.Time     `json:"filled_at,omitempty"`
	CancelledAt     *time.Time     `json:"cancelled_at,omitempty"`
	CreatedAt       *time.Time     `json:"created_at"`
	UpdatedAt       *time.Time     `json:"-"`
}

func NewBuyOrder(accountId int, orderType, symbol string, amount float64) *Order {
	return &Order{
		AccountID: accountId,
		Type:      OrderType(orderType),
		Status:    StatusPending,
		Direction: DirectionBuy,
		Amount:    amount,
		Symbol:    strings.ToUpper(symbol),
	}
}

func NewSellOrder(accountID int, orderType, symbol string, quantity int) *Order {
	return &Order{
		AccountID: accountID,
		Type:      OrderType(orderType),
		Status:    StatusPending,
		Direction: DirectionSell,
		Quantity:  quantity,
		Symbol:    strings.ToUpper(symbol),
	}
}

func NewStopLossOrder(accountID, parentID int, symbol string, amount float64) *Order {
	return &Order{
		AccountID:     accountID,
		ParentOrderID: parentID,
		Type:          TypeStopLoss,
		Status:        StatusPending,
		Direction:     DirectionBuy,
		Amount:        amount,
		Symbol:        strings.ToUpper(symbol),
	}
}

func NewTakeProfitOrder(accountID, parentID int, symbol string, amount float64) *Order {
	return &Order{
		AccountID:     accountID,
		ParentOrderID: parentID,
		Type:          TypeTakeProfit,
		Status:        StatusPending,
		Direction:     DirectionBuy,
		Amount:        amount,
		Symbol:        strings.ToUpper(symbol),
	}
}

func (o *Order) FillOrder(price float64) {
	now := time.Now()
	o.FillPrice = price
	o.FilledAt = &now
	o.Status = StatusFilled
	if o.Direction == DirectionBuy {
		o.Quantity = int(o.Amount / price)
	}
	o.AmountAfterFill = o.FillPrice * float64(o.Quantity)
	if o.Direction == DirectionBuy {
		o.AmountAfterFill = -o.AmountAfterFill
	}
}

func (o *Order) GetTotalFillAmount() float64 {
	return o.AmountAfterFill
}

func (o *Order) GetAmount() float64 {
	return o.Amount
}

func (o *Order) GetDirectionString() string {
	return string(o.Direction)
}
