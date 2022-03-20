package activity

import (
	"errors"
	"fmt"

	"github.com/MihaiBlebea/trading-platform/account"
	"github.com/MihaiBlebea/trading-platform/order"
	"github.com/MihaiBlebea/trading-platform/pos"
	"github.com/MihaiBlebea/trading-platform/quotes"
	"github.com/sirupsen/logrus"
)

type Filler struct {
	accountRepo  *account.AccountRepo
	orderRepo    *order.OrderRepo
	positionRepo *pos.PositionRepo
	logger       *logrus.Logger
}

func NewFiller(
	accountRepo *account.AccountRepo,
	orderRepo *order.OrderRepo,
	positionRepo *pos.PositionRepo,
	logger *logrus.Logger) *Filler {

	return &Filler{
		accountRepo:  accountRepo,
		orderRepo:    orderRepo,
		positionRepo: positionRepo,
		logger:       logger,
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
		quotes := quotes.Quotes{}
		bidAsk, err := quotes.GetCurrentPrice(o.Symbol)
		if err != nil {
			f.logger.Error(err)
			continue
		}

		if o.Direction == order.DirectionBuy {
			o.FillOrder(bidAsk.Ask)
		} else {
			o.FillOrder(bidAsk.Bid)
		}

		err = f.updateInternalRecords(&o)
		if err != nil {
			f.logger.Error(err)
			continue
		}
	}

	return nil
}

func (f *Filler) updateInternalRecords(o *order.Order) error {
	err := f.orderRepo.Update(o)
	if err != nil {
		return err
	}
	f.logger.Info(fmt.Sprintf("order filled id: %d", o.ID))

	account, err := f.accountRepo.WithId(o.AccountID)
	if err != nil {
		return err
	}

	account.UpdateBalance(o)
	if err != nil {
		return err
	}
	err = f.accountRepo.Update(account)
	if err != nil {
		return err
	}
	f.logger.Info(
		fmt.Sprintf(
			"account balance updated id %d with amount: %v",
			account.ID,
			o.AmountAfterFill,
		),
	)

	err = f.updatePosition(o)
	if err != nil {
		return err
	}

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
			position = pos.NewPosition(o.AccountID, o.Symbol, o.Quantity)
			_, err := f.positionRepo.Save(position)
			if err != nil {
				return err
			}
		} else {
			position.IncrementQuantity(o.Quantity)
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