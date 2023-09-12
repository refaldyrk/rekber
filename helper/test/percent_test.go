package test

import (
	"rekber/helper"
	"testing"
)

func Test_Percent(t *testing.T) {
	tests := []struct {
		input    int
		expected float64
	}{
		{
			input:    5,
			expected: 0.05,
		},
		{
			input:    10,
			expected: 0.1,
		},
		{
			input:    20,
			expected: 0.2,
		},
		{
			input:    100,
			expected: 1,
		},
		{
			input:    90,
			expected: 0.9,
		},
		{

			input:    32,
			expected: 0.32,
		},
	}
	for _, tt := range tests {
		t.Run("Convert", func(t *testing.T) {
			res := helper.ConvertPercent(tt.input)
			if res != tt.expected {
				t.Error("not expected: got ", res)
			} else {
				t.Log("Expected: ", res)
			}
		})
	}
}
