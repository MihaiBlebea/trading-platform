package cmd

import (
	"fmt"
	"time"

	"github.com/MihaiBlebea/trading-platform/activity"
	"github.com/MihaiBlebea/trading-platform/di"
	"github.com/MihaiBlebea/trading-platform/http"
	"github.com/spf13/cobra"
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

		di, err := di.NewContainer()
		if err != nil {
			return err
		}

		orderFiller := di.GetOrderFiller()

		logger := di.GetLogger()

		go func(orderFiller *activity.Filler) {
			for {
				err := orderFiller.FillPendingOrders()
				if err != nil {
					continue
				}
				time.Sleep(60 * time.Second)
			}
		}(orderFiller)

		http.New(logger)

		return nil
	},
}
