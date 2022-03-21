package activity_test

import (
	"fmt"
	"testing"

	"github.com/MihaiBlebea/trading-platform/account"
	"github.com/MihaiBlebea/trading-platform/activity"
	"github.com/MihaiBlebea/trading-platform/order"
	"github.com/MihaiBlebea/trading-platform/pos"
)

type PositionRepoMock struct {
}

func TestCanPlaceBuyOrder(t *testing.T) {
	orderPlacer := activity.NewOrderPlacer(
		&account.AccountRepoMock{},
		&order.OrderRepoMock{},
		&pos.PositionRepo{},
	)

	order, err := orderPlacer.PlaceOrder(
		"api_token",
		"limit",
		"buy",
		"AAPL",
		1000.00,
		0,
	)
	if err != nil {
		t.Errorf("Error message: %s", err)
		return
	}

	fmt.Println(order)

}
