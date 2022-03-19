package pos

import "time"

type Position struct {
	ID        int        `json:"id"`
	AccountID int        `json:"-"`
	Symbol    string     `json:"symbol"`
	Quantity  int        `json:"quantity"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"-"`
}

func NewPosition(accountId int, symbol string, quantity int) *Position {
	return &Position{
		AccountID: accountId,
		Symbol:    symbol,
		Quantity:  quantity,
	}
}

func (p *Position) IsFound() bool {
	return p.ID != 0
}

func (p *Position) IncrementQuantity(quantity int) {
	p.Quantity += quantity
}

func (p *Position) DecrementQuantity(quantity int) {
	p.Quantity -= quantity
	if p.Quantity < 0 {
		p.Quantity = 0
	}
}
