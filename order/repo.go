package order

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type OrderRepo struct {
	conn *gorm.DB
}

func NewOrderRepo() (*OrderRepo, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/London",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
	)

	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return &OrderRepo{}, err
	}

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
