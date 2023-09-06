package helper

import (
	"fmt"
	"strconv"
)

func ConvertPercent(num int) float64 {
	if num > 100 {
		return 0
	}

	nilaiDesimal := float64(num) / 100.0

	f := fmt.Sprintf("%.2f", nilaiDesimal)
	floats, _ := strconv.ParseFloat(f, 64)
	return floats
}
