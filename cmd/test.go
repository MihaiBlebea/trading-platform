package cmd

import (
	"fmt"

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
		fmt.Println("Starting to migrate symbols")

		c := symbols.NewClient(true)
		s, err := c.MakeRequest("aapl")
		if err != nil {
			return err
		}
		fmt.Printf("%+v", s)

		return nil
	},
}
