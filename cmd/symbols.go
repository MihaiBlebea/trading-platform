package cmd

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"

	"github.com/MihaiBlebea/trading-platform/di"
	"github.com/MihaiBlebea/trading-platform/symbols"
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

		err := di.BuildContainer().Invoke(func(symbolRepo *symbols.SymbolRepo) error {
			f, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer f.Close()

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

				symbols = append(symbols, symbol)
			}

			if _, err := symbolRepo.SaveMany(symbols); err != nil {
				return err
			}

			return nil
		})

		return err
	},
}
