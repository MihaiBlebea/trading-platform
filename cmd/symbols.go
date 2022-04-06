package cmd

import (
	"encoding/csv"
	"fmt"
	"io"
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
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Starting to migrate symbols")

		di := di.NewContainer()

		symbolRepo, err := di.GetSymbolRepo()
		if err != nil {
			return err
		}

		f, err := os.Open("freetrade_universe.csv")
		if err != nil {
			return err
		}
		defer f.Close()

		client := symbs.NewClient(true)

		symbols := []symbs.Symbol{}

		uniqSimbols := map[string]bool{}

		csvReader := csv.NewReader(f)
		for {
			rec, err := csvReader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}

			symbolName := rec[7]

			if symbolName == "Symbol" {
				continue
			}

			fmt.Println("Processing symbol " + symbolName)

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
