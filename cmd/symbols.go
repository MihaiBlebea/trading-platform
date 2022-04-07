package cmd

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"

	"github.com/MihaiBlebea/trading-platform/di"
	symbs "github.com/MihaiBlebea/trading-platform/symbols"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(symbolsCmd)
}

var symbolsCmd = &cobra.Command{
	Use: "symbols",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("provide at least one argument for the file path")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Starting to migrate symbols")
		filePath := args[0]

		di := di.NewContainer()

		symbolRepo, err := di.GetSymbolRepo()
		if err != nil {
			return err
		}

		f, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer f.Close()

		client := symbs.NewClient(true)

		symbols := []symbs.Symbol{}

		uniqSimbols := map[string]bool{}

		filedata, err := csv.NewReader(f).ReadAll()
		if err != nil {
			return err
		}

		total := len(filedata)

		for i, rec := range filedata {
			symbolName := rec[7]

			if symbolName == "Symbol" {
				continue
			}

			fmt.Printf("%d/%d Processing symbol %s\n", i, total, symbolName)

			if _, ok := uniqSimbols[symbolName]; ok {
				fmt.Println("Key already exists: " + symbolName)
				continue
			}

			uniqSimbols[symbolName] = true

			symbol := *symbs.NewSymbol(rec[0], rec[1], rec[2], rec[3], symbolName)

			data, err := client.MakeRequest(symbolName)
			if err != nil {
				fmt.Printf("Error for symbol %s: %v", symbolName, err)
			} else {
				symbol.LongBusinessSummary = data.LongBusinessSummary
				symbol.CashflowStatements = data.CashflowStatements
				symbol.ProfitMargins = data.ProfitMargins
				symbol.SharesOutstanding = data.SharesOutstanding
				symbol.Beta = data.Beta
				symbol.BookValue = data.BookValue
				symbol.PriceToBook = data.PriceToBook
				symbol.EarningsQuarterlyGrowth = data.EarningsQuarterlyGrowth
			}

			symbols = append(symbols, symbol)
		}

		if _, err := symbolRepo.SaveMany(symbols); err != nil {
			return err
		}

		return nil
	},
}
