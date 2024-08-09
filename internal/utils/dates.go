package utils

import (
	"fmt"
	"time"
)

func LastTwelveMonths() [][]string {
	const monthsInYear = 12

	months := make([][]string, 0, 12)

	currentMonth := time.Now().Month() - 1
	currentYear := time.Now().Year()

	for i := currentMonth; i > currentMonth-monthsInYear; i-- {
		if i < time.January {
			year := currentYear - 1
			month := monthsInYear + i

			months = append(months, []string{
				fmt.Sprintf("%d-%02d", year, month),
				fmt.Sprintf("%s-%d", time.Month(month).String(), year),
			})
			continue
		}

		months = append(months, []string{
			fmt.Sprintf("%d-%02d", currentYear, i),
			fmt.Sprintf("%s-%d", time.Month(i).String(), currentYear),
		})
	}
	return months
}

func LastMonth() (time.Month, int) {
	month := time.Now().Month()
	year := time.Now().Year()

	if month-1 < time.January {
		month = time.December
		year = year - 1
	} else {
		month = month - 1
	}

	return month, year
}
