package test

import (
	"rekber/helper"
	"testing"

	"github.com/spf13/viper"
)

func Test_Fee(t *testing.T) {
	viper.SetConfigFile("../../.env")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	tests := []struct {
		input int
	}{
		{
			input: 3000,
		},
		{
			input: 25000,
		},
		{
			input: 35000,
		}, {
			input: 87000,
		},
	}
	for _, tt := range tests {
		t.Run("Fee", func(t *testing.T) {
			res := helper.CalculateFee(tt.input)
			t.Log(res)
		})
	}
}
