package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/MihaiBlebea/trading-platform/order"
	"github.com/MihaiBlebea/trading-platform/quotes"
	server "github.com/MihaiBlebea/trading-platform/server"
)

func init() {
	rootCmd.AddCommand(startServerCmd)
}

var startServerCmd = &cobra.Command{
	Use:   "start-server",
	Short: "Run the API server.",
	Long:  "Run the API server",
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Server is starting")

		l := logrus.New()

		l.SetFormatter(&logrus.JSONFormatter{})
		l.SetOutput(os.Stdout)
		l.SetLevel(logrus.InfoLevel)

		orderRepo, err := order.NewOrderRepo()
		if err != nil {
			return err
		}

		go func(orderRepo *order.OrderRepo) {
			for {
				fillOrders(orderRepo)
				time.Sleep(60 * time.Second)
			}
		}(orderRepo)

		server.New(l)

		return nil
	},
}

func fillOrders(orderRepo *order.OrderRepo) error {
	orders, err := orderRepo.WithPendingStatus()
	if err != nil {
		return err
	}

	if len(orders) == 0 {
		fmt.Println("no orders to fill")
		return nil
	}

	for _, o := range orders {
		quotes := quotes.Quotes{}
		bidAsk, err := quotes.GetCurrentPrice(o.Symbol)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if o.Direction == order.DirectionBuy {
			o.FillOrder(bidAsk.Ask)
		} else {
			o.FillOrder(bidAsk.Bid)
		}

		orderRepo.Update(&o)
	}

	return nil
}
