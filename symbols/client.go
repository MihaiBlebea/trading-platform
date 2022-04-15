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
)

type Client struct {
	baseUrl string
	cache   *redis.Client
	ttl     int
}

func NewClient(cache *redis.Client) *Client {
	return &Client{
		baseUrl: "https://query2.finance.yahoo.com/v7/finance/quote",
		cache:   cache,
		ttl:     60,
	}
}

func (c *Client) makeCacheRequest(symbols []string) ([]Quote, error) {
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

	quotes, err := c.makeRequest(notFound)
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

func (c *Client) makeRequest(symbols []string) ([]Quote, error) {
	url := fmt.Sprintf(
		"%s?symbols=%s",
		c.baseUrl,
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
