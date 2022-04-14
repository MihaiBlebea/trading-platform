package market

import (
	"time"
)

type MarketStatus struct {
	timeKeeper TimeKeeperInterface
}

type TimeKeeper struct {
}

type TimeKeeperInterface interface {
	Now() time.Time
}

func New(timeKeeper TimeKeeperInterface) *MarketStatus {
	if timeKeeper == nil {
		timeKeeper = NewTimeKeeper()
	}
	return &MarketStatus{timeKeeper}
}

func (m *MarketStatus) IsOpen() bool {
	// london, _ := time.LoadLocation("Europe/London")
	now := m.timeKeeper.Now()

	weekday := now.Weekday()
	if int(weekday) == 0 || int(weekday) == 6 {
		return false
	}

	hour, min, _ := now.Clock()

	if hour < 9 && min < 30 {
		return false
	}

	if hour > 16 && min > 0 {
		return false
	}

	return true
}

func NewTimeKeeper() *TimeKeeper {
	return &TimeKeeper{}
}

func (tk *TimeKeeper) Now() time.Time {
	// london, _ := time.LoadLocation("Europe/London")
	newyork, _ := time.LoadLocation("America/New_York")

	return time.Now().In(newyork)
}
