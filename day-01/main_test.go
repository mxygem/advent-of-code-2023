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
			name: "part 1 example case",
			input: `1abc2
			pqr3stu8vwx
			a1b2c3d4e5f
			treb7uchet`,
			expected: 142,
		},
		{
			name: "part 2 example case",
			input: `two1nine
			eightwothree
			abcone2threexyz
			xtwone3four
			4nineeightseven2
			zoneight234
			7pqrstsixteen`,
			expected: 281,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, calibration(tc.input))
		})
	}
}

func TestParseNumbers(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "empty",
			input:    "",
			expected: nil,
		},
		{
			name:     "newline",
			input:    "\n",
			expected: nil,
		},
		{
			name:     "no number(s) found",
			input:    "abcdefghi",
			expected: nil,
		},
		{
			name:     "single number",
			input:    "6bcde",
			expected: []string{"6"},
		},
		{
			name:     "single number at end",
			input:    "abcdefg8",
			expected: []string{"8"},
		},
		{
			name:     "two single numbers found",
			input:    "1bcd5fgh",
			expected: []string{"1", "5"},
		},
		{
			name:     "three single numbers found",
			input:    "ab3de6gh",
			expected: []string{"3", "6"},
		},
		{
			name:     "ensure 0 and 9 are read",
			input:    "x9x1x0x",
			expected: []string{"9", "1", "0"},
		},
		{
			name:     "only numbers",
			input:    "234",
			expected: []string{"2", "3", "4"},
		},
		{
			name:     "single number spelled out",
			input:    "three",
			expected: []string{"3"},
		},
		{
			name:     "four numbers spelled out",
			input:    "twooneonethree",
			expected: []string{"2", "1", "1", "3"},
		},
		{
			name:     "mixed numbers and spellings",
			input:    "eight6sevenfive30nine",
			expected: []string{"8", "6", "7", "5", "3", "0", "9"},
		},
		{
			name:     "partial spellings",
			input:    "threfour",
			expected: []string{"4"},
		},
		{
			name:     "all mixed",
			input:    "0ontwo3fourfivesi7",
			expected: []string{"0", "2", "3", "4", "5", "7"},
		},
		{
			name:     "puzzle input #1",
			input:    "sevenntgvnrrqfvxh2ttnkgffour8fiveone",
			expected: []string{"7", "2", "4", "8", "5", "1"},
		},
		{
			name:     "overlapping numbers",
			input:    "twonexxthreeight",
			expected: []string{"2", "1", "3", "8"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, parseNumbers(tc.input))
		})
	}
}

func TestCalibrationValue(t *testing.T) {
	testCases := []struct {
		name     string
		input    []string
		expected int
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: 0,
		},
		{
			name:     "two numbers",
			input:    []string{"2", "3"},
			expected: 23,
		},
		{
			name:     "single number is duplicated",
			input:    []string{"6"},
			expected: 66,
		},
		{
			name:     "single zero returned for double zeros",
			input:    []string{"0", "0"},
			expected: 0,
		},
		{
			// this may need to return the remaining value duplicated, aka 11.
			name:     "leading zero",
			input:    []string{"0", "1"},
			expected: 1,
		},
		{
			name:     "only the first and last numbers are used",
			input:    []string{"4", "1", "2"},
			expected: 42,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, calibrationValue(tc.input))
		})
	}
}
