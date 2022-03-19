package account

import (
	"math/rand"
	"time"
)

type Order interface {
	GetTotalFillAmount() float32
	GetAmount() float32
	GetDirectionString() string
}

type Account struct {
	ID             int       `json:"-"`
	ApiToken       string    `json:"api_token"`
	Balance        float32   `json:"balance"`
	PendingBalance float32   `json:"pending_balance"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"-"`
}

func NewAccount() *Account {
	return &Account{ApiToken: genApiKey(10), Balance: 10000.00}
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

func (a *Account) UpdateBalance(order Order) {
	a.Balance += order.GetTotalFillAmount()
	if order.GetDirectionString() == "buy" {
		a.PendingBalance -= order.GetAmount()
	}
}

func (a *Account) HasEnoughBalance(amount float32) bool {
	return a.Balance-a.PendingBalance > amount
}
