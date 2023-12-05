package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalibration(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name: "example case",
			input: `1abc2
			pqr3stu8vwx
			a1b2c3d4e5f
			treb7uchet`,
			expected: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, calibration(tc.input))
		})
	}
}
