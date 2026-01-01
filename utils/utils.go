package utils

import (
	"math"
	"strconv"
)

func TotalPage(limit int, totalData int64) int {
	if totalData <= 0 {
		return 0
	}

	flimit := float64(limit)
	fdata := float64(totalData)

	res := math.Ceil(fdata / flimit)

	return int(res)
}

func StringToBool(name string) bool {
	result, err := strconv.ParseBool(name)
	if err != nil {
		return false
	}
	return result
}

func StringToInt(num string) int {
	result, err := strconv.Atoi(num)
	if err != nil {
		return 0
	}
	return result
}
