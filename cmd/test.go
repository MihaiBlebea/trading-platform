package cmd

import (
	"fmt"

	"github.com/MihaiBlebea/trading-platform/di"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(testCmd)
}

var testCmd = &cobra.Command{
	Use: "test",
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Testing...")
		// q := quotes.New()
		// bidAsk, err := q.GetCurrentPrice("AAPL")
		// if err != nil {
		// 	return err
		// }

		// fmt.Printf("%+v", bidAsk)

		service, err := di.NewContainer().GetSymbolService()
		if err != nil {
			return err
		}

		s, err := service.GetSymbols([]string{"AAPL", "TSLA", "FB"})
		if err != nil {
			return err
		}

		fmt.Printf("%+v", s)

		return nil
	},
}
