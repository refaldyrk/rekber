package helper

import (
	"rekber/constant"

	"github.com/spf13/viper"
)

func CalculateFee(num int) int {
	minimumPrice := constant.MINIMUM_PRICE

	keuntunganAwalPercentageS := viper.GetInt("FEE")
	keuntunganAwalPercentage := ConvertPercent(keuntunganAwalPercentageS)
	if num < minimumPrice {
		return 0
	}

	keuntunganAwal := (float64(num) * keuntunganAwalPercentage) / 100.0
	keuntunganAwalInteger := int(keuntunganAwal * 100)

	return keuntunganAwalInteger
}
