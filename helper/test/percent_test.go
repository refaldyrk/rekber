package test

import (
	"rekber/helper"
	"testing"
)

func Test_Percent(t *testing.T) {
	tests := []struct {
		input int
	}{
		{
			input: 5,
		},
		{
			input: 10,
		},
		{
			input: 20,
		},
		{
			input: 100,
		},
		{
			input: 90,
		},
	}
	for _, tt := range tests {
		t.Run("Convert", func(t *testing.T) {
			res := helper.ConvertPercent(tt.input)
			t.Log(res)
		})
	}
}
