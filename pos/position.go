package pos

import "time"

type Position struct {
	ID                 int        `json:"id"`
	AccountID          int        `json:"-"`
	Symbol             string     `json:"symbol"`
	Quantity           float64    `json:"quantity"`                          // Up to date quantity of stocks in this position, updated when buying and selling stocks
	BoughtTotalPrice   float64    `json:"bought_total_price"`                // Only updated when stocks are bought, this does not get updated when selling stocks
	BoughtQuantity     float64    `json:"bought_quantity"`                   // Only the total quantity bought, this does not get updated when stocks are sold
	TotalValue         float64    `json:"total_value" gorm:"-:all"`          // Total value of the stocks at current market price, only populated when returned to user
	AverageBoughtPrice float64    `json:"average_bought_price" gorm:"-:all"` // Average price at which the stock has been bought, updated after each buy
	CreatedAt          *time.Time `json:"created_at"`
	UpdatedAt          *time.Time `json:"-"`
}

func NewPosition(accountId int, symbol string, quantity, fillPrice float64) *Position {
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

func (p *Position) IncrementQuantity(quantity, fillPrice float64) {
	p.Quantity += quantity
	p.BoughtQuantity += quantity
	p.BoughtTotalPrice += fillPrice * float64(quantity)
}

func (p *Position) DecrementQuantity(quantity float64) {
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
	p.AverageBoughtPrice = p.BoughtTotalPrice / p.BoughtQuantity
}
