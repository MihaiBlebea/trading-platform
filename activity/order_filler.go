package activity

import (
	"fmt"

	"github.com/MihaiBlebea/trading-platform/account"
	"github.com/MihaiBlebea/trading-platform/order"
	"github.com/MihaiBlebea/trading-platform/quotes"
	"github.com/sirupsen/logrus"
)

type Filler struct {
	accountRepo *account.AccountRepo
	orderRepo   *order.OrderRepo
	logger      *logrus.Logger
}

func NewFiller(accountRepo *account.AccountRepo, orderRepo *order.OrderRepo, logger *logrus.Logger) *Filler {
	return &Filler{
		accountRepo: accountRepo,
		orderRepo:   orderRepo,
		logger:      logger,
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

	return nil
}
