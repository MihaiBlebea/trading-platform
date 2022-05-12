package account

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Order interface {
	GetTotalFillAmount() float64
	GetAmount() float64
	IsBuyOrder() bool
}

type Account struct {
	ID             int       `json:"-"`
	Username       string    `json:"username"`
	Email          string    `json:"email" gorm:"uniqueIndex"`
	Password       string    `json:"-"`
	ApiToken       string    `json:"api_token"`
	Balance        float64   `json:"balance"`
	PendingBalance float64   `json:"pending_balance"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"-"`
}

func NewAccount(username, email, password string) (*Account, error) {
	hash, err := hashPassword(password)
	if err != nil {
		return &Account{}, err
	}

	return &Account{
		Username: username,
		Email:    email,
		Password: hash,
		ApiToken: genApiKey(10),
		Balance:  10000.00,
	}, nil
}

func genApiKey(n int) string {
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	rand.Seed(time.Now().UnixNano())

	b := make([]byte, n)
	for i, _ := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(b)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (a *Account) UpdateBalance(order Order) {
	if order.IsBuyOrder() {
		a.Balance -= order.GetTotalFillAmount()
		a.PendingBalance -= order.GetAmount()
	} else {
		a.Balance += order.GetTotalFillAmount()
	}
}

func (a *Account) HasEnoughBalance(amount float64) bool {
	return a.Balance-a.PendingBalance > amount
}

func (a *Account) CheckPasswordHash(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password))
	return err == nil
}
