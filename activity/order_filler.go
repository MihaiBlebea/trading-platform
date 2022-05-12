package activity

import (
	"errors"
	"fmt"

	"github.com/MihaiBlebea/trading-platform/order"
	"github.com/MihaiBlebea/trading-platform/pos"
	"github.com/sirupsen/logrus"
)

type Filler struct {
	accountRepo   AccountRepo
	orderRepo     OrderRepo
	positionRepo  PositionRepo
	symbolService SymbolService
	logger        *logrus.Logger
}

func NewFiller(
	accountRepo AccountRepo,
	orderRepo OrderRepo,
	positionRepo PositionRepo,
	symbolService SymbolService,
	logger *logrus.Logger) *Filler {

	return &Filler{
		accountRepo:   accountRepo,
		orderRepo:     orderRepo,
		positionRepo:  positionRepo,
		symbolService: symbolService,
		logger:        logger,
	}
}

func (f *Filler) FillPendingOrders() error {
	orders, err := f.orderRepo.WithPendingStatus()
	if err != nil {
		return err
	}

	if len(orders) == 0 {
		f.logger.Info("no orders to fill")
		return nil
	}

	for _, o := range orders {
		ask, bid, isOpen, err := f.symbolService.GetCurrentMarketStatus(o.Symbol)
		if err != nil {
			f.logger.Error(err)
			continue
		}

		if !isOpen {
			f.logger.Info("market is not open")
			continue
		}

		if o.Direction == order.DirectionBuy {
			o.FillOrder(ask)
		} else {
			o.FillOrder(bid)
		}

		if err := f.orderRepo.Update(&o); err != nil {
			f.logger.Error(err)
			continue
		}

		if err := f.updateAccount(&o); err != nil {
			f.logger.Error((err))
			continue
		}

		if err := f.updatePosition(&o); err != nil {
			f.logger.Error(err)
			continue
		}
	}

	return nil
}

func (f *Filler) updateAccount(o *order.Order) error {
	f.logger.Info(fmt.Sprintf("order filled id: %d", o.ID))

	account, err := f.accountRepo.WithId(o.AccountID)
	if err != nil {
		return err
	}

	account.UpdateBalance(o)

	if err := f.accountRepo.Update(account); err != nil {
		return err
	}
	f.logger.Info(
		fmt.Sprintf(
			"account balance updated id %d with amount: %v",
			account.ID,
			o.AmountAfterFill,
		),
	)

	return nil
}

func (f *Filler) updatePosition(o *order.Order) error {
	position, err := f.positionRepo.WithAccountAndSymbol(o.AccountID, o.Symbol)
	if err != nil {
		return err
	}

	if o.Direction == order.DirectionBuy {
		if position.ID == 0 {
			// There is no position yet
			position = pos.NewPosition(o.AccountID, o.Symbol, o.Quantity, o.FillPrice)
			_, err := f.positionRepo.Save(position)
			if err != nil {
				return err
			}
		} else {
			position.IncrementQuantity(o.Quantity, o.FillPrice)
			err := f.positionRepo.Update(position)
			if err != nil {
				return err
			}
		}
	} else {
		if position.ID == 0 {
			// There is no position yet
			// This should not be possible
			return errors.New("position that you want to sell is not found")
		}

		position.DecrementQuantity(o.Quantity)

		if position.Quantity == 0 {
			err := f.positionRepo.Delete(position)
			if err != nil {
				return err
			}
		} else {
			err := f.positionRepo.Update(position)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
