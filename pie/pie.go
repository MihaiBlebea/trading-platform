package pie

import (
	"errors"
	"time"
)

type Pie struct {
	ID        int        `json:"id"`
	AccountID int        `json:"-"`
	Name      string     `json:"name"`
	Slices    []PieSlice `json:"slices" gorm:"-:all"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"-"`
}

type PieSlice struct {
	ID        int        `json:"id"`
	PieID     int        `json:"-"`
	Size      float64    `json:"size"`
	Quantity  float64    `json:"quantity"`
	Symbol    string     `json:"symbol"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"-"`
}

func NewPie(accountID int, name string, slices []PieSlice) (*Pie, error) {
	var total float64
	for _, slice := range slices {
		total += slice.Size
	}
	if total < 1 {
		return &Pie{}, errors.New("insufficient slices to make a pie")
	}
	if total > 1 {
		return &Pie{}, errors.New("too many slices for a single pie")
	}

	return &Pie{
		AccountID: accountID,
		Name:      name,
		Slices:    slices,
	}, nil
}
