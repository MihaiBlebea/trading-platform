package quotes

import (
	"strings"

	finance "github.com/piquette/finance-go"

	"github.com/shopspring/decimal"
)

type Quote struct {
	Symbol    string
	Open      float32
	Close     float32
	High      float32
	Low       float32
	Volume    int
	Timestamp int
}

func NewQuote(bar *finance.ChartBar, symbol string) *Quote {
	return &Quote{
		Symbol:    strings.ToUpper(symbol),
		Open:      toFloat32(bar.Open),
		Close:     toFloat32(bar.Close),
		High:      toFloat32(bar.High),
		Low:       toFloat32(bar.Low),
		Volume:    bar.Volume,
		Timestamp: bar.Timestamp,
	}
}

func toFloat32(val decimal.Decimal) float32 {
	return float32(val.InexactFloat64())
}
