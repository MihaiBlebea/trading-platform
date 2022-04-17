package cmd

import (
	"fmt"

	"github.com/MihaiBlebea/trading-platform/di"
	"github.com/MihaiBlebea/trading-platform/symbols"
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

		container := di.BuildContainer()

		err := container.Invoke(func(ss *symbols.Service) {
			fmt.Println(ss.Exists("aapl"))
		})

		if err != nil {
			return err
		}

		return nil
	},
}
