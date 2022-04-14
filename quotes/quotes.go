package quotes

import (
	"fmt"
	"time"

	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
	"github.com/piquette/finance-go/quote"
)

type Quotes struct {
	quotes []Quote
}

func New() *Quotes {
	return &Quotes{}
}

func (q *Quotes) GetQuotes(symbol string, startDate string, interval string) ([]Quote, error) {
	startTime, err := toTimeDate(startDate)
	if err != nil {
		return []Quote{}, err
	}
	endTime := time.Now()

	params := &chart.Params{
		Symbol:   symbol,
		Start:    datetime.New(&startTime),
		End:      datetime.New(&endTime),
		Interval: datetime.OneHour,
	}
	iter := chart.Get(params)

	quotes := []Quote{}
	for iter.Next() {
		quotes = append(quotes, *NewQuote(iter.Bar(), symbol))
	}
	if err := iter.Err(); err != nil {
		return quotes, err
	}

	q.quotes = quotes

	return q.quotes, nil
}

func (q *Quotes) GetCurrentPrice(symbol string) (*BidAsk, error) {
	raw, err := quote.Get(symbol)
	fmt.Printf("%+v", raw)
	if err != nil {
		return &BidAsk{}, err
	}

	return NewBidAsk(symbol, float32(raw.Bid), float32(raw.Ask)), nil
}

func toTimeDate(val string) (time.Time, error) {
	layout := "2006-01-02 15:04"
	t, err := time.Parse(layout, val)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}
