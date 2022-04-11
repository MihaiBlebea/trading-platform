package activity_test

import (
	"strings"
	"testing"
	"time"

	"github.com/MihaiBlebea/trading-platform/account"
	"github.com/MihaiBlebea/trading-platform/activity"
	"github.com/MihaiBlebea/trading-platform/order"
)

func createAccount(t *testing.T, accountRepo activity.AccountRepo) account.Account {
	a, err := account.NewAccount("FakeUsername", "test@gmail.com", "1234")
	if err != nil {
		t.Fatal("could not create an account")
	}
	account, err := accountRepo.Save(a)
	if err != nil {
		t.Fatal("could not create an account")
	}

	return *account
}

func TestCanCancelBuyOrder(t *testing.T) {
	accountRepo := account.AccountRepoMock{}
	account := createAccount(t, &accountRepo)

	orderRepo := order.OrderRepoMock{}

	orderCanceller := activity.NewOrderCanceller(
		&accountRepo,
		&orderRepo,
	)

	amount := float32(1000.00)
	symbol := "aapl"
	o, err := orderRepo.Save(order.NewBuyOrder(
		account.ID,
		"limit",
		symbol,
		amount,
	))
	if err != nil {
		t.Errorf("error message: %s", err)
		return
	}

	if o.Amount != amount {
		t.Errorf("expected order amount %v, got %v", amount, o.Amount)
	}

	if o.Symbol != strings.ToUpper(symbol) {
		t.Errorf(
			"expected order symbol %s, got %s",
			strings.ToUpper(symbol),
			o.Symbol,
		)
	}

	if o.Status != order.StatusPending {
		t.Errorf(
			"expected order status %s, got %s",
			order.StatusPending,
			o.Status,
		)
	}

	if o.Direction != order.DirectionBuy {
		t.Errorf(
			"expected order direction %s, got %s",
			order.DirectionBuy,
			o.Direction,
		)
	}

	if o.Type != order.TypeLimit {
		t.Errorf(
			"expected order type %s, got %s",
			order.TypeLimit,
			o.Type,
		)
	}

	o, err = orderCanceller.CancelOrder(account.ApiToken, o.ID)
	if err != nil {
		t.Errorf("error message: %s", err)
		return
	}

	if o.CancelledAt.Before(time.Now()) == false {
		t.Errorf("expected cancelled time to be before now")
	}

	if o.Status != order.StatusCancelled {
		t.Errorf(
			"expected order status %s, got %s",
			order.StatusCancelled,
			o.Status,
		)
	}
}
