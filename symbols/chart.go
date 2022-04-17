package symbols

import "github.com/MihaiBlebea/trading-platform/symbols/yahoofin"

type Chart struct {
	Timestamp int     `json:"timestamp"`
	Open      float64 `json:"open"`
	Close     float64 `json:"close"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Volume    int     `json:"volume"`
}

func fromYahooCharts(charts []yahoofin.Chart) []Chart {
	res := []Chart{}
	for _, c := range charts {
		res = append(res, Chart{
			c.Timestamp,
			c.Open,
			c.Close,
			c.High,
			c.Low,
			c.Volume,
		})
	}

	return res
}
