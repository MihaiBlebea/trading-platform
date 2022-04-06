package symbols

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

type SymbolRepo struct {
	conn *gorm.DB
}

func NewSymbolRepo(conn *gorm.DB) (*SymbolRepo, error) {
	if err := conn.AutoMigrate(&Symbol{}); err != nil {
		return &SymbolRepo{}, err
	}

	return &SymbolRepo{conn: conn}, nil
}

func (sr *SymbolRepo) Save(sym *Symbol) (*Symbol, error) {
	resp := sr.conn.Create(sym)

	if resp.Error != nil {
		return &Symbol{}, resp.Error
	}
	return sym, nil
}

func (sr *SymbolRepo) SaveMany(symbs []Symbol) ([]Symbol, error) {
	resp := sr.conn.Create(&symbs)

	if resp.Error != nil {
		return []Symbol{}, resp.Error
	}
	return symbs, nil
}

func (sr *SymbolRepo) WithSymbol(symbol string) (*Symbol, error) {
	sym := Symbol{}
	err := sr.conn.Where("symbol = ?", strings.ToUpper(symbol)).Find(&sym).Error
	if err != nil {
		return &Symbol{}, err
	}

	if sym.ID == 0 {
		return &Symbol{}, errors.New("could not find record")
	}

	return &sym, err
}
