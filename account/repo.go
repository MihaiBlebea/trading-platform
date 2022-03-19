package account

import (
	"errors"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AccountRepo struct {
	conn *gorm.DB
}

func NewAccountRepo() (*AccountRepo, error) {
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
		return &AccountRepo{}, err
	}

	if err := conn.AutoMigrate(&Account{}); err != nil {
		return &AccountRepo{}, err
	}

	return &AccountRepo{conn: conn}, nil
}

func (ar *AccountRepo) Save(account *Account) (*Account, error) {
	resp := ar.conn.Create(account)

	if resp.Error != nil {
		return &Account{}, resp.Error
	}
	return account, nil
}

func (ar *AccountRepo) Update(account *Account) error {
	return ar.conn.Save(account).Error
}

func (ar *AccountRepo) WithToken(token string) (*Account, error) {
	account := Account{}
	err := ar.conn.Where("api_token = ?", token).Find(&account).Error
	if err != nil {
		return &account, err
	}

	if account.ID == 0 {
		return &account, errors.New("could not find record")
	}

	return &account, err
}

func (ar *AccountRepo) WithId(id int) (*Account, error) {
	account := Account{}
	err := ar.conn.Where("id = ?", id).Find(&account).Error
	if err != nil {
		return &account, err
	}

	if account.ID == 0 {
		return &account, errors.New("could not find record")
	}

	return &account, err
}
