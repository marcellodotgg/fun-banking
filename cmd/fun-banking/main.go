package main

import (
	"github.com/bytebury/fun-banking/internal/api"
	"github.com/bytebury/fun-banking/internal/infrastructure/persistence"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic("unable to read .env file")
	}
	persistence.Connect()
	persistence.RunMigrations()
	api.Start()
}
