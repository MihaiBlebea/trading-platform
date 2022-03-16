package main

import (
	"github.com/MihaiBlebea/trading-platform/cmd"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load("./.env")
}

func main() {
	cmd.Execute()
}
