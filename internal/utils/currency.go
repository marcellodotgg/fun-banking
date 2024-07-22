package utils

import (
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"
)

func FormatCurrency(amount float64) string {
	return formatUSD(amount)
}

func FormatNumber(amount int64) string {
	return humanize.Comma(amount)
}

func formatUSD(amount float64) string {
	formattedAmount := humanize.CommafWithDigits(amount, 2)

	currencyParts := strings.Split(formattedAmount, ".")

	if len(currencyParts) == 1 {
		return fmt.Sprintf("$%s.00", formattedAmount)
	}
	if len(currencyParts) == 2 && len(currencyParts[1]) < 2 {
		return fmt.Sprintf("$%s0", formattedAmount)
	}
	return fmt.Sprintf("$%s", formattedAmount)
}
