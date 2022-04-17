package symbols

import (
	"strings"
	"time"

	"github.com/MihaiBlebea/trading-platform/symbols/yahoofin"
)

type Symbol struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	LongTitle   string     `json:"long_title"`
	Industry    string     `json:"industry"`
	Currency    string     `json:"currency"`
	Symbol      string     `json:"symbol" gorm:"uniqueIndex"`
	Description string     `json:"description,omitempty" gorm:"-:all"`
	MarketCap   int        `json:"market_cap,omitempty" gorm:"-:all"`
	Bid         float64    `json:"bid,omitempty" gorm:"-:all"`
	Ask         float64    `json:"ask,omitempty" gorm:"-:all"`
	MarketPrice float64    `json:"market_price,omitempty" gorm:"-:all"`
	MarketState string     `json:"market_state,omitempty" gorm:"-:all"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"-"`
}

func NewSymbol(title, longTitle, industry, currency, symbol string) *Symbol {
	return &Symbol{
		Title:     title,
		LongTitle: longTitle,
		Industry:  industry,
		Currency:  currency,
		Symbol:    strings.ToUpper(symbol),
	}
}

func (s *Symbol) IsMarketOpen() bool {
	return s.MarketState == "REGULAR"
}

func (s *Symbol) fromQuote(quote *yahoofin.Quote) {
	s.Ask = quote.Ask
	s.Bid = quote.Bid
	s.MarketPrice = quote.RegularMarketPrice
	s.MarketState = string(quote.MarketState)
}
