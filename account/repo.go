package account

import (
	"errors"

	"gorm.io/gorm"
)

type AccountRepo struct {
	conn *gorm.DB
}

func NewAccountRepo(conn *gorm.DB) (*AccountRepo, error) {
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

func (ar *AccountRepo) WithEmail(email string) (*Account, error) {
	account := Account{}
	err := ar.conn.Where("email = ?", email).Find(&account).Error
	if err != nil {
		return &account, err
	}

	if account.ID == 0 {
		return &account, errors.New("could not find record")
	}

	return &account, err
}
