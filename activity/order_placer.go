package activity

import (
	"errors"

	"github.com/MihaiBlebea/trading-platform/account"
	"github.com/MihaiBlebea/trading-platform/order"
)

type OrderPlacer struct {
	accountRepo  AccountRepo
	orderRepo    OrderRepo
	positionRepo PositionRepo
}

func NewOrderPlacer(
	accountRepo AccountRepo,
	orderRepo OrderRepo,
	positionRepo PositionRepo) *OrderPlacer {

	return &OrderPlacer{
		accountRepo:  accountRepo,
		orderRepo:    orderRepo,
		positionRepo: positionRepo,
	}
}

func (op *OrderPlacer) PlaceOrder(
	apiToken,
	orderType,
	direction,
	symbol string,
	amount,
	quantity,
	stopLoss,
	takeProfit float64) (*order.Order, error) {

	// Get the user account
	account, err := op.accountRepo.WithToken(apiToken)
	if err != nil {
		return &order.Order{}, err
	}

	// Check if this is a buy or sell order
	if direction == string(order.DirectionBuy) {
		o, err := op.PlaceBuyOrder(account, orderType, symbol, amount)
		if err != nil {
			return &order.Order{}, err
		}

		// Create extra stop loss order
		if stopLoss != 0 {
			_, err = op.orderRepo.Save(
				order.NewStopLossOrder(account.ID, o.ID, symbol, stopLoss),
			)
			if err != nil {
				return &order.Order{}, err
			}
		}

		// Create extra take profit order
		if takeProfit != 0 {
			_, err = op.orderRepo.Save(
				order.NewTakeProfitOrder(account.ID, o.ID, symbol, takeProfit),
			)
			if err != nil {
				return &order.Order{}, err
			}
		}

		return o, nil
	}

	return op.PlaceSellOrder(account, orderType, symbol, quantity)
}

func (op *OrderPlacer) PlaceBuyOrder(
	account *account.Account,
	orderType,
	symbol string,
	amount float64) (*order.Order, error) {

	if amount == 0 {
		return &order.Order{}, errors.New("need to specify an amount greater than 0")
	}

	if !account.HasEnoughBalance(amount) {
		return &order.Order{}, errors.New("insufficient balance to place this order")
	}

	o := order.NewBuyOrder(account.ID, orderType, symbol, amount)
	o, err := op.orderRepo.Save(o)
	if err != nil {
		return &order.Order{}, err
	}

	account.PendingBalance = amount
	err = op.accountRepo.Update(account)
	if err != nil {
		return &order.Order{}, err
	}

	return o, nil
}

func (op *OrderPlacer) PlaceSellOrder(
	account *account.Account,
	orderType,
	symbol string,
	quantity float64) (*order.Order, error) {

	if quantity == 0 {
		return &order.Order{}, errors.New("need to specify a quantity greater than 0")
	}

	position, err := op.positionRepo.WithAccountAndSymbol(account.ID, symbol)
	if err != nil {
		return &order.Order{}, err
	}

	if !position.IsFound() {
		return &order.Order{}, errors.New("could not find position")
	}

	if position.Quantity < quantity {
		return &order.Order{}, errors.New("position quantity is too low")
	}

	// Check if there is already a pending order for this symbol
	orders, err := op.orderRepo.WithDirectionStatusSymbolAndAccountId(
		order.DirectionSell,
		order.StatusPending,
		account.ID,
		symbol,
	)
	if err != nil {
		return &order.Order{}, err
	}
	if len(orders) > 0 {
		return &order.Order{}, errors.New("selling pending order already exists")
	}

	o := order.NewSellOrder(account.ID, orderType, symbol, quantity)
	o, err = op.orderRepo.Save(o)
	if err != nil {
		return &order.Order{}, err
	}

	return o, nil
}
