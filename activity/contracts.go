package activity

import (
	"github.com/MihaiBlebea/trading-platform/account"
	"github.com/MihaiBlebea/trading-platform/order"
	"github.com/MihaiBlebea/trading-platform/pos"
)

type AccountRepo interface {
	Save(account *account.Account) (*account.Account, error)
	Update(account *account.Account) error
	WithToken(token string) (*account.Account, error)
	WithId(id int) (*account.Account, error)
}

type OrderRepo interface {
	Save(order *order.Order) (*order.Order, error)
	Update(order *order.Order) error
	WithPendingStatus() ([]order.Order, error)
	WithAccountId(accountId int) ([]order.Order, error)
	WithId(id int) (*order.Order, error)
}

type PositionRepo interface {
	Save(pos *pos.Position) (*pos.Position, error)
	Update(pos *pos.Position) error
	WithAccountAndSymbol(accountId int, symbol string) (*pos.Position, error)
	WithAccountId(accountId int) ([]pos.Position, error)
	Delete(pos *pos.Position) error
}
