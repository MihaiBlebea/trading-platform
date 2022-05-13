package pie

import (
	"errors"

	"github.com/MihaiBlebea/trading-platform/order"
)

type Service struct {
	accountRepo AccountRepo
	pieRepo     *PieRepo
	orderRepo   OrderRepo
}

func NewService(
	accountRepo AccountRepo,
	pieRepo *PieRepo,
	orderRepo OrderRepo) *Service {

	return &Service{
		accountRepo,
		pieRepo,
		orderRepo,
	}
}

func (s *Service) Buy(accountID, pieID int, amount float64) error {
	if amount == 0 {
		return errors.New("need to specify an amount greater than 0")
	}

	acc, err := s.accountRepo.WithId(accountID)
	if err != nil {
		return err
	}

	if !acc.HasEnoughBalance(amount) {
		return errors.New("insufficient balance to place this order")
	}

	pie, err := s.pieRepo.WithId(pieID)
	if err != nil {
		return err
	}

	if acc.ID != pie.AccountID {
		return errors.New("account does not own the pie")
	}

	acc.PendingBalance += amount
	err = s.accountRepo.Update(acc)
	if err != nil {
		return err
	}

	for _, slice := range pie.Slices {
		sliceAmount := slice.Size * amount

		o := order.NewBuyOrder(acc.ID, string(order.TypePie), slice.Symbol, sliceAmount)
		if _, err := s.orderRepo.Save(o); err != nil {
			return err
		}
	}

	err = s.accountRepo.Update(acc)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Sell(accountID, pieID int, amount float64) error {
	if amount == 0 {
		return errors.New("need to specify an amount greater than 0")
	}

	acc, err := s.accountRepo.WithId(accountID)
	if err != nil {
		return err
	}

	pie, err := s.pieRepo.WithId(pieID)
	if err != nil {
		return err
	}

	if acc.ID != pie.AccountID {
		return errors.New("account does not own the pie")
	}

	for _, slice := range pie.Slices {
		sliceAmount := slice.Size * amount

		// o := order.NewSellOrder(acc.ID, string(order.TypePie), slice.Symbol, sliceAmount)
		// if _, err := s.orderRepo.Save(o); err != nil {
		// 	return err
		// }
	}

	return nil
}

func (s *Service) Recalibrate() {

}
