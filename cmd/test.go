package cmd

import (
	"fmt"
	"log"

	"github.com/MihaiBlebea/trading-platform/di"
	"github.com/MihaiBlebea/trading-platform/pie"
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

		err := container.Invoke(func(pieRepo *pie.PieRepo) {
			// pie, err := pie.NewPie(1, "my new fancy pie", []pie.PieSlice{
			// 	{Size: 0.5, Symbol: "AAPL"},
			// 	{Size: 0.5, Symbol: "TSLA"},
			// })
			// if err != nil {
			// 	log.Fatal(err)
			// 	return
			// }

			// pieRepo.Save(pie)

			p, err := pieRepo.WithId(2)
			if err != nil {
				log.Fatal(err)
				return
			}

			fmt.Printf("%+v", p)
		})

		if err != nil {
			return err
		}

		return nil
	},
}
