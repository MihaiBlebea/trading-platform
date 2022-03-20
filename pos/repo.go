package pos

import (
	"strings"

	"gorm.io/gorm"
)

type PositionRepo struct {
	conn *gorm.DB
}

func NewPositionRepo(conn *gorm.DB) (*PositionRepo, error) {
	if err := conn.AutoMigrate(&Position{}); err != nil {
		return &PositionRepo{}, err
	}

	return &PositionRepo{conn: conn}, nil
}

func (or *PositionRepo) Save(pos *Position) (*Position, error) {
	resp := or.conn.Create(pos)

	if resp.Error != nil {
		return &Position{}, resp.Error
	}
	return pos, nil
}

func (or *PositionRepo) Update(pos *Position) error {
	return or.conn.Save(pos).Error
}

func (or *PositionRepo) WithAccountAndSymbol(accountId int, symbol string) (*Position, error) {
	pos := Position{}
	err := or.conn.Where(
		"account_id = ? AND symbol = ?",
		accountId,
		strings.ToUpper(symbol),
	).Find(&pos).Error
	if err != nil {
		return &Position{}, err
	}

	return &pos, err
}

func (or *PositionRepo) WithAccountId(accountId int) ([]Position, error) {
	pos := []Position{}
	err := or.conn.Where(
		"account_id = ?",
		accountId,
	).Find(&pos).Error
	if err != nil {
		return []Position{}, err
	}

	return pos, err
}

func (or *PositionRepo) Delete(pos *Position) error {
	return or.conn.Delete(pos).Error
}
