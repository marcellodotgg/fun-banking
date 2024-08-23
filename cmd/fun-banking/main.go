package main

import (
	"fmt"
	"time"

	"github.com/bytebury/fun-banking/internal/api"
	"github.com/bytebury/fun-banking/internal/domain"
	"github.com/bytebury/fun-banking/internal/infrastructure/persistence"
	"github.com/bytebury/fun-banking/internal/service"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic("unable to read .env file")
	}
	persistence.Connect()
	persistence.RunMigrations()

	// Cron Jobs
	c := cron.New()

	// Every day at midnight
	c.AddFunc("0 0 * * *", func() {
		fulfillAutoPay()
	})

	c.Start()

	api.Start()
}

func fulfillAutoPay() {
	persistence.DB.Transaction(func(tx *gorm.DB) error {
		var autoPays []domain.AutoPay

		persistence.DB.Find(&autoPays, "active = 1 AND strftime('%Y-%m-%d', start_date) <= ?", time.Now().Format("2006-01-02"))

		if len(autoPays) == 0 {
			return nil
		}

		transactionService := service.NewTransactionService()

		for _, autoPay := range autoPays {
			if err := transactionService.AutoPay(autoPay); err != nil {
				fmt.Println("Error! Auto Pay failed...")
			}
		}

		return nil
	})
}
