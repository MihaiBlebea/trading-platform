package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	server "github.com/MihaiBlebea/trading-platform/server"
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

		server.New(l)

		return nil
	},
}
