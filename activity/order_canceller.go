package activity

import (
	"errors"
	"time"

	"github.com/MihaiBlebea/trading-platform/order"
)

type OrderCanceller struct {
	accountRepo AccountRepo
	orderRepo   OrderRepo
}

func NewOrderCanceller(
	accountRepo AccountRepo,
	orderRepo OrderRepo) *OrderCanceller {

	return &OrderCanceller{
		accountRepo: accountRepo,
		orderRepo:   orderRepo,
	}
}

func (oc *OrderCanceller) CancelOrder(
	apiToken string,
	orderId int) (*order.Order, error) {

	if apiToken == "" {
		return &order.Order{}, errors.New("api token cannot be null")
	}

	if orderId == 0 {
		return &order.Order{}, errors.New("order id cannot be null")
	}

	account, err := oc.accountRepo.WithToken(apiToken)
	if err != nil {
		return &order.Order{}, err
	}

	o, err := oc.orderRepo.WithId(orderId)
	if err != nil {
		return &order.Order{}, err
	}

	if o.AccountID != account.ID {
		return &order.Order{}, errors.New("order account id does not match")
	}

	if o.Status != order.StatusPending {
		return &order.Order{}, errors.New("order status can not be cancelled")
	}

	now := time.Now()
	o.Status = order.StatusCancelled
	o.CancelledAt = &now

	if err := oc.orderRepo.Update(o); err != nil {
		return &order.Order{}, err
	}

	return o, nil
}
