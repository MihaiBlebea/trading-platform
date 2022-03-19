package activity

import (
	"errors"

	"github.com/MihaiBlebea/trading-platform/account"
	"github.com/MihaiBlebea/trading-platform/order"
)

type OrderPlacer struct {
	accountRepo *account.AccountRepo
	orderRepo   *order.OrderRepo
}

func NewOrderPlacer(accountRepo *account.AccountRepo, orderRepo *order.OrderRepo) *OrderPlacer {
	return &OrderPlacer{
		accountRepo: accountRepo,
		orderRepo:   orderRepo,
	}
}

func (op *OrderPlacer) PlaceOrder(apiToken, orderType, direction, symbol string, amount float32) (*order.Order, error) {
	account, err := op.accountRepo.WithToken(apiToken)
	if err != nil {
		return &order.Order{}, err
	}

	if account.Balance < amount {
		return &order.Order{}, errors.New("insufficient balance to place this order")
	}

	o := order.NewOrder(account.ID, orderType, direction, amount, symbol)
	o, err = op.orderRepo.Save(o)
	if err != nil {
		return &order.Order{}, err
	}

	return o, nil
}
