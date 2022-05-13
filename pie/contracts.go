package pie

import (
	"github.com/MihaiBlebea/trading-platform/account"
	"github.com/MihaiBlebea/trading-platform/order"
)

type AccountRepo interface {
	WithId(id int) (*account.Account, error)
	Update(account *account.Account) error
}

type OrderRepo interface {
	Save(order *order.Order) (*order.Order, error)
}
