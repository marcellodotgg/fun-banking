package utils

import (
	"time"
)

var months = map[time.Month]string{
	time.January:   "01",
	time.February:  "02",
	time.March:     "03",
	time.April:     "04",
	time.May:       "05",
	time.June:      "06",
	time.July:      "07",
	time.August:    "08",
	time.September: "09",
	time.October:   "10",
	time.November:  "11",
	time.December:  "12",
}

func ConvertMonthToNumeric(month time.Month) string {
	return months[month]
}
