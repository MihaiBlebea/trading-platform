package yahoofin

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type ClientCache struct {
	client *Client
	cache  *redis.Client
}

func NewClientCache(client *Client, cache *redis.Client) *ClientCache {
	return &ClientCache{client, cache}
}

func (c *ClientCache) GetQuotes(symbols []string) ([]Quote, error) {
	ctx := context.Background()
	notFound := []string{}
	found := []Quote{}
	key := "quote_%s"
	ttl := time.Second * time.Duration(60)

	for _, symbol := range symbols {
		res, err := c.cache.Get(ctx, fmt.Sprintf(key, symbol)).Result()
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

	quotes, err := c.client.GetQuotes(notFound)
	if err != nil {
		return []Quote{}, err
	}

	for _, q := range quotes {
		b, err := json.Marshal(q)
		if err != nil {
			return []Quote{}, err
		}

		c.cache.Set(
			ctx,
			fmt.Sprintf(key, q.Symbol),
			string(b),
			ttl,
		)
	}

	return append(found, quotes...), nil
}

func (c *ClientCache) GetChart(symbol string) ([]Chart, error) {
	ctx := context.Background()
	key := "chart_%s"
	ttl := time.Hour * time.Duration(24)

	if res, err := c.cache.Get(ctx, fmt.Sprintf(key, symbol)).Result(); err == nil {
		// cache was a hit
		charts := []Chart{}
		if err := json.Unmarshal([]byte(res), &charts); err != nil {
			fmt.Println(err)
		}

		return charts, nil
	}

	charts, err := c.client.GetChart(symbol)
	if err != nil {
		return []Chart{}, err
	}

	b, err := json.Marshal(charts)
	if err != nil {
		return []Chart{}, err
	}

	c.cache.Set(
		ctx,
		fmt.Sprintf(key, symbol),
		string(b),
		ttl,
	)

	return charts, nil
}
