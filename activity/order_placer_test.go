package activity_test

import (
	"strings"
	"testing"

	"github.com/MihaiBlebea/trading-platform/account"
	"github.com/MihaiBlebea/trading-platform/activity"
	"github.com/MihaiBlebea/trading-platform/order"
	"github.com/MihaiBlebea/trading-platform/pos"
)

func TestCanCreateAccount(t *testing.T) {
	accountRepo := account.AccountRepoMock{}
	account := createAccount(t, &accountRepo)

	if account.ApiToken == "" {
		t.Error("could not create valid account")
	}
}

func TestCanPlaceBuyOrder(t *testing.T) {
	accountRepo := account.AccountRepoMock{}
	account := createAccount(t, &accountRepo)

	orderPlacer := activity.NewOrderPlacer(
		&accountRepo,
		&order.OrderRepoMock{},
		&pos.PositionRepo{},
	)

	amount := float64(1000.00)
	symbol := "aapl"
	o, err := orderPlacer.PlaceOrder(
		account.ApiToken,
		"limit",
		"buy",
		symbol,
		amount,
		0,
		0,
		0,
	)
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
}

func TestCanPlaceSellOrder(t *testing.T) {
	accountRepo := account.AccountRepoMock{}
	account := createAccount(t, &accountRepo)

	orderRepo := order.OrderRepoMock{}
	posRepo := pos.NewPositionRepoMock()

	orderPlacer := activity.NewOrderPlacer(
		&accountRepo,
		&orderRepo,
		posRepo,
	)

	// Update position before placing sell order
	symbol := "aapl"
	quantity := 50
	posRepo.Save(pos.NewPosition(account.ID, symbol, quantity))

	// Place sell order
	o, err := orderPlacer.PlaceOrder(
		account.ApiToken,
		"limit",
		"sell",
		symbol,
		0,
		quantity,
		0,
		0,
	)
	if err != nil {
		t.Errorf("error message: %s", err)
		return
	}

	if o.Quantity != quantity {
		t.Errorf("expected order quantity %d, got %v", quantity, o.Quantity)
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

	if o.Direction != order.DirectionSell {
		t.Errorf(
			"expected order direction %s, got %s",
			order.DirectionSell,
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
}

func TestSellOrderInsufficientQuantity(t *testing.T) {
	accountRepo := account.AccountRepoMock{}
	account := createAccount(t, &accountRepo)

	orderRepo := order.OrderRepoMock{}
	posRepo := pos.NewPositionRepoMock()

	orderPlacer := activity.NewOrderPlacer(
		&accountRepo,
		&orderRepo,
		posRepo,
	)

	// Update position before placing sell order
	symbol := "aapl"
	quantity := 50
	posRepo.Save(pos.NewPosition(account.ID, symbol, 10))

	// Place sell order
	_, err := orderPlacer.PlaceOrder(
		account.ApiToken,
		"limit",
		"sell",
		symbol,
		0,
		quantity,
		0,
		0,
	)
	errMessage := "position quantity is too low"
	if err.Error() != errMessage {
		t.Errorf("error message expected %s: got %s", errMessage, err)
	}
}
