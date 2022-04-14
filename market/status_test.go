package market_test

import (
	"testing"
	"time"

	"github.com/MihaiBlebea/trading-platform/market"
)

var cases = []struct {
	now      time.Time
	expected bool
}{
	{now: time.Date(2022, 4, 10, 20, 34, 58, 0, time.UTC), expected: false},
	{now: time.Date(2022, 4, 13, 6, 12, 58, 0, time.UTC), expected: false},
	{now: time.Date(2022, 4, 13, 10, 0, 58, 0, time.UTC), expected: true},
}

type TimeKeeper struct {
	now time.Time
}

func (tk *TimeKeeper) Now() time.Time {
	return tk.now
}

func TestMarketDuringWeekday(t *testing.T) {
	for _, c := range cases {
		tk := TimeKeeper{c.now}

		status := market.New(&tk)

		isOpen := status.IsOpen()

		if isOpen != c.expected {
			t.Errorf("market is open expected %v, received %v", c.expected, isOpen)
		}
	}
}
