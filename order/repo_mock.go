package order

import "errors"

type OrderRepoMock struct {
	orders []Order
}

func (or *OrderRepoMock) Save(order *Order) (*Order, error) {
	or.orders = append(or.orders, *order)
	order.ID = len(or.orders) + 1

	return order, nil
}

func (or *OrderRepoMock) Update(order *Order) error {
	if len(or.orders) < order.ID-1 {
		return errors.New("could not find index")
	}

	or.orders[order.ID-1] = *order

	return nil
}

func (or *OrderRepoMock) WithPendingStatus() ([]Order, error) {
	orders := []Order{}
	for _, order := range or.orders {
		if order.Status == StatusPending {
			orders = append(orders, order)
		}
	}

	return orders, nil
}

func (or *OrderRepoMock) WithAccountId(accountId int) ([]Order, error) {
	orders := []Order{}
	for _, order := range or.orders {
		if order.AccountID == accountId {
			orders = append(orders, order)
		}
	}

	return orders, nil
}
