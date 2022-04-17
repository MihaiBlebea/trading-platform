package symbols

import (
	"errors"
	"strings"

	"github.com/MihaiBlebea/trading-platform/symbols/yahoofin"
)

type YahooClient interface {
	GetQuotes(symbols []string) ([]yahoofin.Quote, error)
	GetChart(symbol string) ([]yahoofin.Chart, error)
}

type Service struct {
	client YahooClient
	repo   *SymbolRepo
}

func NewService(client YahooClient, repo *SymbolRepo) *Service {
	return &Service{client, repo}
}

func (s *Service) Exists(symbol string) bool {
	if _, err := s.repo.WithSymbol(symbol); err != nil {
		return false
	}

	return true
}

func (s *Service) GetCurrentMarketStatus(symbol string) (ask float64, bid float64, _ bool, _ error) {
	symb, err := s.GetSymbol(symbol)
	if err != nil {
		return 0, 0, false, err
	}

	return symb.Ask, symb.Bid, symb.IsMarketOpen(), nil
}

func (s *Service) GetSymbol(symbol string) (*Symbol, error) {
	symbols, err := s.getSymbols([]string{symbol})
	if err != nil {
		return &Symbol{}, err
	}

	if len(symbols) == 0 {
		return &Symbol{}, errors.New("could not find symbol")
	}

	return &symbols[0], nil
}

func (s *Service) GetSymbols(symbols []string) ([]Symbol, error) {
	tickers, err := s.getSymbols(symbols)
	if err != nil {
		return []Symbol{}, err
	}

	return tickers, nil
}

func (s *Service) GetChart(symbol string) ([]Chart, error) {
	charts, err := s.client.GetChart(symbol)
	if err != nil {
		return []Chart{}, err
	}

	return fromYahooCharts(charts), nil
}

func (s *Service) getSymbols(symbols []string) ([]Symbol, error) {
	upperSymbols := []string{}
	for _, s := range symbols {
		upperSymbols = append(upperSymbols, strings.ToUpper(s))
	}

	tickers, err := s.repo.WithSymbols(upperSymbols)
	if err != nil {
		return []Symbol{}, err
	}

	quotes, err := s.client.GetQuotes(upperSymbols)
	if err != nil {
		return []Symbol{}, err
	}

	res := []Symbol{}
	for i, t := range tickers {
		t.fromQuote(&quotes[i])
		res = append(res, t)
	}

	return res, nil
}
