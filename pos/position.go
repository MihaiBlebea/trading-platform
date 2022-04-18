package pos

import "time"

type Position struct {
	ID                 int        `json:"id"`
	AccountID          int        `json:"-"`
	Symbol             string     `json:"symbol"`
	Quantity           int        `json:"quantity"`
	BoughtTotalPrice   float64    `json:"bought_total_price"`
	BoughtQuantity     int        `json:"bought_quantity"`
	TotalValue         float64    `json:"total_value" gorm:"-:all"`
	AverageBoughtPrice float64    `json:"average_bought_price" gorm:"-:all"`
	CreatedAt          *time.Time `json:"created_at"`
	UpdatedAt          *time.Time `json:"-"`
}

func NewPosition(accountId int, symbol string, quantity int, fillPrice float64) *Position {
	pos := &Position{
		AccountID: accountId,
		Symbol:    symbol,
		Quantity:  quantity,
	}
	pos.IncrementQuantity(quantity, fillPrice)

	return pos
}

func (p *Position) IsFound() bool {
	return p.ID != 0
}

func (p *Position) IncrementQuantity(quantity int, fillPrice float64) {
	p.Quantity += quantity
	p.BoughtQuantity += quantity
	p.BoughtTotalPrice += fillPrice * float64(quantity)
}

func (p *Position) DecrementQuantity(quantity int) {
	p.Quantity -= quantity
	if p.Quantity < 0 {
		p.Quantity = 0
	}
}

func (p *Position) CalculateAverageBoughtPrice() {
	if p.BoughtTotalPrice == 0 || p.BoughtQuantity == 0 {
		p.AverageBoughtPrice = 0
		return
	}
	p.AverageBoughtPrice = p.BoughtTotalPrice / float64(p.BoughtQuantity)
}
