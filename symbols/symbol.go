package symbols

import (
	"encoding/json"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Symbol struct {
	ID                      int                 `json:"id"`
	Title                   string              `json:"title"`
	LongTitle               string              `json:"long_title"`
	Industry                string              `json:"industry"`
	Currency                string              `json:"currency"`
	Symbol                  string              `json:"symbol" gorm:"uniqueIndex"`
	LongBusinessSummary     string              `json:"longBusinessSummary"`
	CashflowStatements      []CashFlowStatement `json:"cashflowStatements" gorm:"-"`
	CashflowStatementsRaw   string              `json:"-"`
	ProfitMargins           float64             `json:"profitMargins"`
	SharesOutstanding       int                 `json:"sharesOutstanding"`
	Beta                    float64             `json:"beta"`
	BookValue               float64             `json:"bookValue"`
	PriceToBook             float64             `json:"priceToBook"`
	EarningsQuarterlyGrowth float64             `json:"earningsQuarterlyGrowth"`
	CreatedAt               *time.Time          `json:"created_at"`
	UpdatedAt               *time.Time          `json:"-"`
}

func NewSymbol(title, longTitle, industry, currency, symbol string) *Symbol {
	return &Symbol{
		Title:     title,
		LongTitle: longTitle,
		Industry:  industry,
		Currency:  currency,
		Symbol:    strings.ToUpper(symbol),
	}
}

func (s *Symbol) BeforeCreate(tx *gorm.DB) error {
	b, err := json.Marshal(s.CashflowStatements)
	if err != nil {
		return err
	}

	s.CashflowStatementsRaw = string(b)

	return nil
}

func (s *Symbol) AfterFind(tx *gorm.DB) error {
	if s.CashflowStatementsRaw == "" {
		return nil
	}

	cfs := []CashFlowStatement{}

	if err := json.Unmarshal([]byte(s.CashflowStatementsRaw), &cfs); err != nil {
		return err
	}

	s.CashflowStatements = cfs

	return nil
}
