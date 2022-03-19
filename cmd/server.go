package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/MihaiBlebea/trading-platform/account"
	"github.com/MihaiBlebea/trading-platform/activity"
	"github.com/MihaiBlebea/trading-platform/http"
	"github.com/MihaiBlebea/trading-platform/order"
	"github.com/MihaiBlebea/trading-platform/pos"
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

		accountRepo, err := account.NewAccountRepo()
		if err != nil {
			return err
		}

		positionRepo, err := pos.NewPositionRepo()
		if err != nil {
			return err
		}

		filler := activity.NewFiller(accountRepo, orderRepo, positionRepo, l)

		go func(orderRepo *order.OrderRepo) {
			for {
				err := filler.FillPendingOrders()
				if err != nil {
					continue
				}
				time.Sleep(60 * time.Second)
			}
		}(orderRepo)

		http.New(l)

		return nil
	},
}
