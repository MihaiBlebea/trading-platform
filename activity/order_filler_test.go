package activity_test

import (
	"io/ioutil"
	"testing"

	"github.com/MihaiBlebea/trading-platform/account"
	"github.com/MihaiBlebea/trading-platform/activity"
	"github.com/MihaiBlebea/trading-platform/order"
	"github.com/MihaiBlebea/trading-platform/pos"
	"github.com/sirupsen/logrus"
)

type SymbolServiceStub struct {
}

func (s *SymbolServiceStub) GetCurrentMarketStatus(symbol string) (float64, float64, bool, error) {
	return 20, 22, true, nil
}

func TestCanFillBuyOrder(t *testing.T) {
	logger := logrus.New()
	logger.Out = ioutil.Discard

	accountRepo := account.AccountRepoMock{}
	account := createAccount(t, &accountRepo)

	orderRepo := order.OrderRepoMock{}
	posRepo := pos.NewPositionRepoMock()

	symbolService := SymbolServiceStub{}

	amount := float64(1000.00)
	symbol := "aapl"
	o, err := orderRepo.Save(order.NewBuyOrder(
		account.ID,
		"limit",
		symbol,
		1000.00,
	))
	if err != nil {
		t.Errorf("error message: %s", err)
		return
	}

	if o.Amount != amount {
		t.Errorf("expected order amount %v, got %v", amount, o.Amount)
	}

	orderFiller := activity.NewFiller(
		&accountRepo,
		&orderRepo,
		posRepo,
		&symbolService,
		logger,
	)

	if err := orderFiller.FillPendingOrders(); err != nil {
		t.Errorf("error message: %s", err)
		return
	}

	o, err = orderRepo.WithId(o.ID)
	if err != nil {
		t.Errorf("error message: %s", err)
		return
	}

	if o.Status != order.StatusFilled {
		t.Errorf(
			"expected order amount %v, got %v",
			order.StatusFilled,
			o.Status,
		)
	}
}

func TestCanFillSellOrder(t *testing.T) {
	logger := logrus.New()
	logger.Out = ioutil.Discard

	accountRepo := account.AccountRepoMock{}
	account := createAccount(t, &accountRepo)

	orderRepo := order.OrderRepoMock{}
	posRepo := pos.NewPositionRepoMock()

	symbolService := SymbolServiceStub{}

	symbol := "aapl"
	quantity := 50
	posRepo.Save(pos.NewPosition(account.ID, symbol, quantity, 1077.44))

	o, err := orderRepo.Save(order.NewSellOrder(
		account.ID,
		"limit",
		symbol,
		quantity,
	))
	if err != nil {
		t.Errorf("error message: %s", err)
		return
	}

	if o.Quantity != quantity {
		t.Errorf(
			"expected order quantity %v, got %v",
			quantity,
			o.Quantity,
		)
	}

	orderFiller := activity.NewFiller(
		&accountRepo,
		&orderRepo,
		posRepo,
		&symbolService,
		logger,
	)

	if err := orderFiller.FillPendingOrders(); err != nil {
		t.Errorf("error message: %s", err)
		return
	}

	o, err = orderRepo.WithId(o.ID)
	if err != nil {
		t.Errorf("error message: %s", err)
		return
	}

	if o.Status != order.StatusFilled {
		t.Errorf(
			"expected order amount %v, got %v",
			order.StatusFilled,
			o.Status,
		)
	}
}
