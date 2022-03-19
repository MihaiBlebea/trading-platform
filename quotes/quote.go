package quotes

import (
	"strings"

	finance "github.com/piquette/finance-go"

	"github.com/shopspring/decimal"
)

type Quote struct {
	Symbol    string  `json:"symbol"`
	Open      float32 `json:"open"`
	Close     float32 `json:"close"`
	High      float32 `json:"high"`
	Low       float32 `json:"low"`
	Volume    int     `json:"volume"`
	Timestamp int     `json:"timestamp"`
}

type BidAsk struct {
	Symbol string  `json:"symbol"`
	Bid    float32 `json:"bid"`
	Ask    float32 `json:"ask"`
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

func NewQuoteFromRaw(symbol string, open, close, high, low float32, volume, timestamp int) *Quote {
	return &Quote{
		Symbol:    strings.ToUpper(symbol),
		Open:      open,
		Close:     close,
		High:      high,
		Low:       low,
		Volume:    volume,
		Timestamp: timestamp,
	}
}

func NewBidAsk(symbol string, bid, ask float32) *BidAsk {
	return &BidAsk{symbol, bid, ask}
}

func toFloat32(val decimal.Decimal) float32 {
	return float32(val.InexactFloat64())
}
