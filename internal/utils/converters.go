package utils

import (
	"math"
	"strconv"
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

func ConvertToIntPointer(id string) (*int, error) {
	returnValue := new(int)
	idAsInt, err := strconv.Atoi(id)
	if err != nil {
		return returnValue, err
	}
	*returnValue = int(idAsInt)
	return returnValue, nil
}

func ConvertMonthToNumeric(month time.Month) string {
	return months[month]
}

func GetDollarAmount(amount string) (float64, error) {
	floatAmount, err := strconv.ParseFloat(amount, 64)

	if err != nil {
		return 0, err
	}

	floatAmount = math.Round(floatAmount*100) / 100
	return floatAmount, nil
}

func SafelyAddDollars(a, b float64) float64 {
	sum := a + b
	return math.Round(sum*100) / 100
}

func SafelySubtractDollars(a, b float64) float64 {
	sum := a - b
	return math.Round(sum*100) / 100
}
