package quotes

import (
	"time"

	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
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

	return q.addTrends(), nil
}

func (q *Quotes) addTrends() []Quote {
	return q.quotes
}

func toTimeDate(val string) (time.Time, error) {
	layout := "2006-01-02 15:04"
	t, err := time.Parse(layout, val)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}
