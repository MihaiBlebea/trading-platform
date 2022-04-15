package symbols

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/tidwall/gjson"
)

type Client struct {
	quoteBaseUrl string
	chartBaseUrl string
	cache        *redis.Client
	ttl          int
}

func NewClient(cache *redis.Client) *Client {
	return &Client{
		quoteBaseUrl: "https://query2.finance.yahoo.com/v7/finance/quote",
		chartBaseUrl: "https://query2.finance.yahoo.com/v8/finance/chart",
		cache:        cache,
		ttl:          60,
	}
}

func (c *Client) makeQuoteCacheRequest(symbols []string) ([]Quote, error) {
	ctx := context.Background()
	notFound := []string{}
	found := []Quote{}
	for _, symbol := range symbols {
		res, err := c.cache.Get(ctx, symbol).Result()
		if err != nil {
			notFound = append(notFound, symbol)
			continue
		}

		quote := Quote{}
		if err := json.Unmarshal([]byte(res), &quote); err != nil {
			fmt.Println(err)
		}

		found = append(found, quote)
	}

	quotes, err := c.makeQuoteRequest(notFound)
	if err != nil {
		return []Quote{}, err
	}

	for _, q := range quotes {
		b, err := json.Marshal(q)
		if err != nil {
			return []Quote{}, err
		}

		c.cache.Set(ctx, q.Symbol, string(b), time.Second*time.Duration(c.ttl))
	}

	return append(found, quotes...), nil
}

func (c *Client) makeQuoteRequest(symbols []string) ([]Quote, error) {
	url := fmt.Sprintf(
		"%s?symbols=%s",
		c.quoteBaseUrl,
		strings.Join(symbols, ","),
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []Quote{}, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []Quote{}, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []Quote{}, err
	}

	var result struct {
		QuoteResponse struct {
			Result []Quote `json:"result"`
		} `json:"quoteResponse"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return []Quote{}, err
	}

	return result.QuoteResponse.Result, nil
}

func (c *Client) makeChartRequest(symbol string) ([]Chart, error) {
	url := fmt.Sprintf(
		"%s/%s?range=%s",
		c.chartBaseUrl,
		strings.ToUpper(symbol),
		"1d",
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []Chart{}, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []Chart{}, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []Chart{}, err
	}

	timestamps := []int{}
	timestampRaw := gjson.Get(string(body), "chart.result.0.timestamp")
	if timestampRaw.Exists() {
		for _, ts := range timestampRaw.Array() {
			timestamps = append(timestamps, int(ts.Int()))
		}
	}

	opens := []float64{}
	openRaw := gjson.Get(string(body), "chart.result.0.indicators.quote.0.open")
	if openRaw.Exists() {
		for _, o := range openRaw.Array() {
			opens = append(opens, o.Float())
		}
	}

	closes := []float64{}
	closeRaw := gjson.Get(string(body), "chart.result.0.indicators.quote.0.close")
	if closeRaw.Exists() {
		for _, c := range closeRaw.Array() {
			closes = append(closes, c.Float())
		}
	}

	highs := []float64{}
	highRaw := gjson.Get(string(body), "chart.result.0.indicators.quote.0.high")
	if highRaw.Exists() {
		for _, h := range highRaw.Array() {
			highs = append(highs, h.Float())
		}
	}

	lows := []float64{}
	lowRaw := gjson.Get(string(body), "chart.result.0.indicators.quote.0.low")
	if lowRaw.Exists() {
		for _, l := range lowRaw.Array() {
			lows = append(lows, l.Float())
		}
	}

	volumes := []int{}
	volumeRaw := gjson.Get(string(body), "chart.result.0.indicators.quote.0.volume")
	if volumeRaw.Exists() {
		for _, v := range volumeRaw.Array() {
			volumes = append(volumes, int(v.Int()))
		}
	}

	charts := []Chart{}
	for i, ts := range timestamps {
		charts = append(charts, Chart{
			High:      highs[i],
			Low:       lows[i],
			Open:      opens[i],
			Close:     closes[i],
			Volume:    volumes[i],
			Timestamp: ts,
		})
	}

	return charts, nil
}
