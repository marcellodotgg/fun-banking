package utils

import (
	"math"
	"strconv"
)

func ConvertToIntPointer(id string) (*int, error) {
	returnValue := new(int)
	idAsInt, err := strconv.Atoi(id)
	if err != nil {
		return returnValue, err
	}
	*returnValue = int(idAsInt)
	return returnValue, nil
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
