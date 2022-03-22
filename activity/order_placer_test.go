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
	accountRepo := account.AccountRepoMock{}
	account, _ := accountRepo.Save(account.NewAccount())
	apiToken := account.ApiToken

	orderPlacer := activity.NewOrderPlacer(
		&accountRepo,
		&order.OrderRepoMock{},
		&pos.PositionRepo{},
	)

	fmt.Println(accountRepo)

	order, err := orderPlacer.PlaceOrder(
		apiToken,
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
