package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPartNumberPositions(t *testing.T) {
	testCases := []struct {
		name     string
		lines    []string
		expected [][][]int
	}{
		// {
		// 	name:     "no input",
		// 	lines:    nil,
		// 	expected: nil,
		// },
		// {
		// 	name:     "no lines",
		// 	lines:    []string{},
		// 	expected: nil,
		// },
		// {
		// 	name: "single line, no numbers",
		// 	lines: []string{
		// 		`..........`,
		// 	},
		// 	expected: [][][]int{nil},
		// },
		{
			name: "single line, with two numbers",
			lines: []string{
				`.2...38...`,
			},
			expected: [][][]int{
				{{1, 1}, {5, 6}},
			},
		},
		{
			name: "single numbers separated by single symbols",
			lines: []string{
				`0.2.4.6.8.`,
			},
			expected: [][][]int{
				{{0, 0}, {2, 2}, {4, 4}, {6, 6}, {8, 8}},
			},
		},
		{
			name: "multiple lines",
			lines: []string{
				`$...&*2113`,
				`33!.....99`,
			},
			expected: [][][]int{
				{{6, 9}},
				{{0, 1}, {8, 9}},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, partNumberPositions(tc.lines))
		})
	}
}

func TestNumberPos(t *testing.T) {
	testCases := []struct {
		name          string
		line          string
		expectedStart int
		expectedEnd   int
	}{
		{
			name:          "no number found",
			line:          `....#.....`,
			expectedStart: -1,
			expectedEnd:   -1,
		},
		{
			name:          "number at start",
			line:          `123.#.....`,
			expectedStart: 0,
			expectedEnd:   2,
		},
		{
			name:          "number at end",
			line:          `*#..#..*89`,
			expectedStart: 8,
			expectedEnd:   9,
		},
		{
			name:          "whole line is a single number",
			line:          `9876543210`,
			expectedStart: 0,
			expectedEnd:   9,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			start, end := numberPos(tc.line)

			assert.Equal(t, tc.expectedStart, start)
			assert.Equal(t, tc.expectedEnd, end)
		})
	}
}

func TestIsPartNumber(t *testing.T) {
	testCases := []struct {
		name     string
		lines    []string
		lineNum  int
		pos      []int
		expected bool
	}{
		// {
		// 	name: "line number greater than lines",
		// 	lines: []string{
		// 		"..........",
		// 		"..........",
		// 	},
		// 	lineNum:  2,
		// 	expected: false,
		// },
		// {
		// 	name:     "start less than 0",
		// 	pos:      []int{-1, 4},
		// 	expected: false,
		// },
		// {
		// 	name: "end greater than line length",
		// 	lines: []string{
		// 		"..........",
		// 	},
		// 	pos:      []int{5, 200},
		// 	expected: false,
		// },
		// {
		// 	name: "middle area part number via prev row start diag",
		// 	lines: []string{
		// 		".....*....",
		// 		"......33..",
		// 		"..........",
		// 	},
		// 	lineNum:  1,
		// 	pos:      []int{7, 8},
		// 	expected: true,
		// },
		// {
		// 	name: "middle area number number via prev row end diag",
		// 	lines: []string{
		// 		"........!.",
		// 		"......33..",
		// 		"..........",
		// 	},
		// 	lineNum:  1,
		// 	pos:      []int{7, 8},
		// 	expected: true,
		// },
		// {
		// 	name: "middle area number number via same row before num",
		// 	lines: []string{
		// 		"..........",
		// 		".....$33..",
		// 		"..........",
		// 	},
		// 	lineNum:  1,
		// 	pos:      []int{7, 8},
		// 	expected: true,
		// },
		// {
		// 	name: "middle area number number via same row after num",
		// 	lines: []string{
		// 		"..........",
		// 		"..23(.....",
		// 		"..........",
		// 	},
		// 	lineNum:  1,
		// 	pos:      []int{2, 3},
		// 	expected: true,
		// },
		// {
		// 	name: "middle area number number next row start diag",
		// 	lines: []string{
		// 		".#........",
		// 		"..23......",
		// 		"..........",
		// 	},
		// 	lineNum:  1,
		// 	pos:      []int{2, 3},
		// 	expected: true,
		// },
		// {
		// 	name: "middle area number number next row start diag",
		// 	lines: []string{
		// 		"..........",
		// 		"..23......",
		// 		".%........",
		// 	},
		// 	lineNum:  1,
		// 	pos:      []int{2, 3},
		// 	expected: true,
		// },
		{
			name: "middle area number number next row end diag",
			lines: []string{
				"..........",
				"..23......",
				"....@.....",
			},
			lineNum:  1,
			pos:      []int{2, 3},
			expected: true,
		},
		{
			name: "first line number symbol after",
			lines: []string{
				"....2113#.",
				"..........",
			},
			lineNum:  0,
			pos:      []int{4, 7},
			expected: true,
		},
		{
			name: "first line number symbol before",
			lines: []string{
				"^499......",
				"..........",
			},
			lineNum:  0,
			pos:      []int{1, 3},
			expected: true,
		},
		{
			name: "first line end number symbol before",
			lines: []string{
				"........&0",
				"..........",
			},
			lineNum:  0,
			pos:      []int{9, 9},
			expected: true,
		},
		{
			name: "first line almost end number symbol after",
			lines: []string{
				"........1*",
				"..........",
			},
			lineNum:  0,
			pos:      []int{8, 8},
			expected: true,
		},
		{
			name: "end line mid number symbol after",
			lines: []string{
				"..........",
				"..........",
				"....4123(.",
			},
			lineNum:  2,
			pos:      []int{5, 8},
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var start, end int
			if tc.pos != nil {
				start = tc.pos[0]
				end = tc.pos[1]
			}

			assert.Equal(t, tc.expected, isPartNumber(tc.lines, tc.lineNum, start, end))
		})
	}
}
