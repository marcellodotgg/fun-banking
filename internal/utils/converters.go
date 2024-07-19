package utils

import "strconv"

func ConvertToUintPointer(id string) (*uint, error) {
	returnValue := new(uint)
	idAsInt, err := strconv.Atoi(id)
	if err != nil {
		return returnValue, err
	}
	*returnValue = uint(idAsInt)
	return returnValue, nil
}
