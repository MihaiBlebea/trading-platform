package cmd

import (
	"errors"
	"fmt"

	"github.com/MihaiBlebea/trading-platform/di"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

func init() {
	rootCmd.AddCommand(dropTableCmd)
}

var dropTableCmd = &cobra.Command{
	Use: "drop-table",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("provide at least one argument for the table name")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		tableName := args[0]

		err := di.BuildContainer().Invoke(func(conn *gorm.DB) {
			fmt.Printf("Dropping table %s \n", tableName)
			conn.Migrator().DropTable(tableName)
		})

		return err
	},
}
