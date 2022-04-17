package yahoofin

type ClientStub struct {
}

func NewStubClient() *ClientStub {
	return &ClientStub{}
}

func (c *ClientStub) GetQuotes(symbols []string) ([]Quote, error) {
	return []Quote{
		{Bid: 10.00, Ask: 12.00, RegularMarketPrice: 11.00, MarketState: "REGULAR", Symbol: "AAPL"},
		{Bid: 50.00, Ask: 52.00, RegularMarketPrice: 51.00, MarketState: "REGULAR", Symbol: "TSLA"},
	}, nil
}

func (c *ClientStub) GetChart(symbol string) ([]Chart, error) {

	return []Chart{
		{Timestamp: 1650211238, Open: 10.00, Close: 15.00, High: 20.00, Low: 9.00, Volume: 100},
		{Timestamp: 1650124838, Open: 10.00, Close: 15.00, High: 20.00, Low: 9.00, Volume: 100},
		{Timestamp: 1650038438, Open: 10.00, Close: 15.00, High: 20.00, Low: 9.00, Volume: 100},
	}, nil
}
