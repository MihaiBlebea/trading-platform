package order

import (
	"gorm.io/gorm"
)

type OrderRepo struct {
	conn *gorm.DB
}

func NewOrderRepo(conn *gorm.DB) (*OrderRepo, error) {
	if err := conn.AutoMigrate(&Order{}); err != nil {
		return &OrderRepo{}, err
	}

	return &OrderRepo{conn: conn}, nil
}

func (or *OrderRepo) Save(order *Order) (*Order, error) {
	resp := or.conn.Create(order)

	if resp.Error != nil {
		return &Order{}, resp.Error
	}
	return order, nil
}

func (or *OrderRepo) Update(order *Order) error {
	return or.conn.Save(order).Error
}

func (or *OrderRepo) WithPendingStatus() ([]Order, error) {
	orders := []Order{}
	err := or.conn.Where("status = 'pending'").Find(&orders).Error
	if err != nil {
		return []Order{}, err
	}

	return orders, err
}

func (or *OrderRepo) WithAccountId(accountId int) ([]Order, error) {
	orders := []Order{}
	err := or.conn.Where("account_id = ?", accountId).Find(&orders).Error
	if err != nil {
		return []Order{}, err
	}

	return orders, err
}

func (or *OrderRepo) WithId(id int) (*Order, error) {
	order := Order{}
	err := or.conn.Where("id = ?", id).Find(&order).Error
	if err != nil {
		return &Order{}, err
	}

	return &order, err
}
