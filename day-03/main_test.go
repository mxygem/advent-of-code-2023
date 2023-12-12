package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPartNumberSum(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name: "no part numbers",
			input: `....#..987
			123.#.....
			....#.3...`,
			expected: 0,
		},
		{
			name: "four part numbers, not gears",
			input: `....10....
			..10##10..
			....10....`,
			expected: 0,
		},
		{
			name: "single gear",
			input: `....10....
			....*10...`,
			expected: 100,
		},
		{
			name: "example",
			input: `467..114..
			...*......
			..35..633.
			......#...
			617*......
			.....+.58.
			..592.....
			......755.
			...$.*....
			.664.598..`,
			expected: 467835,
		},
		{
			name: "dedupe",
			input: `....10....
			..10*10...`,
			expected: 0,
		},
		{
			name: "dedupe with extra",
			input: `...10...99
			        .10*10..*1`,
			expected: 99,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, partNumberSum(tc.input))
		})
	}
}

func TestParts(t *testing.T) {
	testCases := []struct {
		name     string
		lines    []string
		expected []*part
	}{
		{
			name:  "no input",
			lines: nil,
		},
		{
			name:  "no lines",
			lines: []string{},
		},
		{
			name: "single line, no numbers",
			lines: []string{
				`..........`,
			},
		},
		{
			name: "single line, with two numbers, no parts",
			lines: []string{
				`.2...38...`,
			},
		},
		{
			name: "single numbers separated by single symbols",
			lines: []string{
				`0*2@4!6^8.`,
			},
			expected: []*part{
				{val: 0, pos: pos{line: 0, start: 0, end: 0}, symbol: &symbol{kind: "*", pos: pos{line: 0, start: 1}}},
				{val: 2, pos: pos{line: 0, start: 2, end: 2}, symbol: &symbol{kind: "*", pos: pos{line: 0, start: 1}}},
				{val: 4, pos: pos{line: 0, start: 4, end: 4}, symbol: &symbol{kind: "@", pos: pos{line: 0, start: 3}}},
				{val: 6, pos: pos{line: 0, start: 6, end: 6}, symbol: &symbol{kind: "!", pos: pos{line: 0, start: 5}}},
				{val: 8, pos: pos{line: 0, start: 8, end: 8}, symbol: &symbol{kind: "^", pos: pos{line: 0, start: 7}}},
			},
		},
		{
			name: "multiple lines",
			lines: []string{
				`$...&*2113`,
				`33!.....99`,
			},
			expected: []*part{
				{val: 2113, pos: pos{line: 0, start: 6, end: 9}, symbol: &symbol{kind: "*", pos: pos{line: 0, start: 5}}},
				{val: 33, pos: pos{line: 1, start: 0, end: 1}, symbol: &symbol{kind: "$", pos: pos{line: 0, start: 0}}},
			},
		},
		{
			name: "tens",
			lines: []string{
				`...10...99`,
				`.10*10..*1`,
			},
			expected: []*part{
				{val: 10, pos: pos{line: 0, start: 3, end: 4}, symbol: &symbol{kind: "*", pos: pos{line: 1, start: 3}}},
				{val: 99, pos: pos{line: 0, start: 8, end: 9}, symbol: &symbol{kind: "*", pos: pos{line: 1, start: 8}}},
				{val: 10, pos: pos{line: 1, start: 1, end: 2}, symbol: &symbol{kind: "*", pos: pos{line: 1, start: 3}}},
				{val: 10, pos: pos{line: 1, start: 4, end: 5}, symbol: &symbol{kind: "*", pos: pos{line: 1, start: 3}}},
				{val: 1, pos: pos{line: 1, start: 9, end: 9}, symbol: &symbol{kind: "*", pos: pos{line: 1, start: 8}}},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, parts(tc.lines))
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
		name       string
		lines      []string
		part       *part
		lineNum    int
		pos        []int
		expected   *part
		expectedOK bool
	}{
		{
			name: "line number greater than lines",
			lines: []string{
				"..........",
				"..........",
			},
			part:       &part{pos: pos{line: 2}},
			expectedOK: false,
		},
		{
			name:       "start less than 0",
			part:       &part{pos: pos{start: -1}},
			expectedOK: false,
		},
		{
			name: "end greater than line length",
			lines: []string{
				"..........",
			},
			part:       &part{pos: pos{end: 200}},
			expectedOK: false,
		},
		{
			name: "middle area part number via prev row start diag",
			lines: []string{
				".....*....",
				"......33..",
				"..........",
			},
			part: &part{
				pos: pos{line: 1, start: 6, end: 7},
			},
			expected: &part{
				pos:    pos{line: 1, start: 6, end: 7},
				symbol: &symbol{kind: "*", pos: pos{line: 0, start: 5}},
			},
			expectedOK: true,
		},
		{
			name: "middle area number via prev row end diag",
			lines: []string{
				"........!.",
				"......33..",
				"..........",
			},
			part: &part{
				pos: pos{line: 1, start: 6, end: 7},
			},
			expected: &part{
				pos:    pos{line: 1, start: 6, end: 7},
				symbol: &symbol{kind: "!", pos: pos{line: 0, start: 8}},
			},
			expectedOK: true,
		},
		{
			name: "middle area number via same row before num",
			lines: []string{
				"..........",
				".....$33..",
				"..........",
			},
			part: &part{
				pos: pos{line: 1, start: 6, end: 7},
			},
			expected: &part{
				pos:    pos{line: 1, start: 6, end: 7},
				symbol: &symbol{kind: "$", pos: pos{line: 1, start: 5}},
			},
			expectedOK: true,
		},
		{
			name: "middle area number via same row after num",
			lines: []string{
				"..........",
				"..23(.....",
				"..........",
			},
			part: &part{
				pos: pos{line: 1, start: 2, end: 3},
			},
			expected: &part{
				pos:    pos{line: 1, start: 2, end: 3},
				symbol: &symbol{kind: "(", pos: pos{line: 1, start: 4}},
			},
			expectedOK: true,
		},
		{
			name: "middle area number prev row start diag",
			lines: []string{
				".#........",
				"..23......",
				"..........",
			},
			part: &part{
				pos: pos{line: 1, start: 2, end: 3},
			},
			expected: &part{
				pos:    pos{line: 1, start: 2, end: 3},
				symbol: &symbol{kind: "#", pos: pos{line: 0, start: 1}},
			},
			expectedOK: true,
		},
		{
			name: "middle area number prev row end diag",
			lines: []string{
				"....#.....",
				"..23......",
				"..........",
			},
			part: &part{
				pos: pos{line: 1, start: 2, end: 3},
			},
			expected: &part{
				pos:    pos{line: 1, start: 2, end: 3},
				symbol: &symbol{kind: "#", pos: pos{line: 0, start: 4}},
			},
			expectedOK: true,
		},
		{
			name: "middle area number next row start diag",
			lines: []string{
				"..........",
				"..23......",
				".%........",
			},
			part: &part{
				pos: pos{line: 1, start: 2, end: 3},
			},
			expected: &part{
				pos:    pos{line: 1, start: 2, end: 3},
				symbol: &symbol{kind: "%", pos: pos{line: 2, start: 1}},
			},
			expectedOK: true,
		},
		{
			name: "middle area number next row end diag",
			lines: []string{
				"..........",
				"..23......",
				"....@.....",
			},
			part: &part{
				pos: pos{line: 1, start: 2, end: 3},
			},
			expected: &part{
				pos:    pos{line: 1, start: 2, end: 3},
				symbol: &symbol{kind: "@", pos: pos{line: 2, start: 4}},
			},
			expectedOK: true,
		},
		{
			name: "first line number symbol after",
			lines: []string{
				"....2113#.",
				"..........",
			},
			part: &part{
				pos: pos{line: 0, start: 4, end: 7},
			},
			expected: &part{
				pos:    pos{line: 0, start: 4, end: 7},
				symbol: &symbol{kind: "#", pos: pos{line: 0, start: 8}},
			},
			expectedOK: true,
		},
		{
			name: "first line number symbol before",
			lines: []string{
				"^499......",
				"..........",
			},
			part: &part{
				pos: pos{line: 0, start: 1, end: 3},
			},
			expected: &part{
				pos:    pos{line: 0, start: 1, end: 3},
				symbol: &symbol{kind: "^", pos: pos{line: 0, start: 0}},
			},
			expectedOK: true,
		},
		{
			name: "first line end number symbol before",
			lines: []string{
				"........&0",
				"..........",
			},
			part: &part{
				pos: pos{line: 0, start: 9, end: 9},
			},
			expected: &part{
				pos:    pos{line: 0, start: 9, end: 9},
				symbol: &symbol{kind: "&", pos: pos{line: 0, start: 8}},
			},
			expectedOK: true,
		},
		{
			name: "first line almost end number symbol after",
			lines: []string{
				"........1*",
				"..........",
			},
			part: &part{
				pos: pos{line: 0, start: 8, end: 8},
			},
			expected: &part{
				pos:    pos{line: 0, start: 8, end: 8},
				symbol: &symbol{kind: "*", pos: pos{line: 0, start: 9}},
			},
			expectedOK: true,
		},
		{
			name: "first line next row start diag",
			lines: []string{
				"...123....",
				"..).......",
			},
			part: &part{
				pos: pos{line: 0, start: 3, end: 5},
			},
			expected: &part{
				pos:    pos{line: 0, start: 3, end: 5},
				symbol: &symbol{kind: ")", pos: pos{line: 1, start: 2}},
			},
			expectedOK: true,
		},
		{
			name: "first line next row after diag",
			lines: []string{
				"...123....",
				"......_...",
			},
			part: &part{
				pos: pos{line: 0, start: 3, end: 5},
			},
			expected: &part{
				pos:    pos{line: 0, start: 3, end: 5},
				symbol: &symbol{kind: "_", pos: pos{line: 1, start: 6}},
			},
			expectedOK: true,
		},
		{
			name: "end line mid number symbol after",
			lines: []string{
				"..........",
				"..........",
				"....4123(.",
			},
			part: &part{
				pos: pos{line: 2, start: 4, end: 7},
			},
			expected: &part{
				pos:    pos{line: 2, start: 4, end: 7},
				symbol: &symbol{kind: "(", pos: pos{line: 2, start: 8}},
			},
			expectedOK: true,
		},
		{
			name: "end line mid number prev line start diag",
			lines: []string{
				"..........",
				"^.........",
				".9874123..",
			},
			part: &part{
				pos: pos{line: 2, start: 1, end: 7},
			},
			expected: &part{
				pos:    pos{line: 2, start: 1, end: 7},
				symbol: &symbol{kind: "^", pos: pos{line: 1, start: 0}},
			},
			expectedOK: true,
		},
		{
			name: "end line mid number prev line end dig",
			lines: []string{
				"..........",
				"........~.",
				".9874123..",
			},
			part: &part{
				pos: pos{line: 2, start: 1, end: 7},
			},
			expected: &part{
				pos:    pos{line: 2, start: 1, end: 7},
				symbol: &symbol{kind: "~", pos: pos{line: 1, start: 8}},
			},
			expectedOK: true,
		},
		{
			name: "not part number symbols all around with one gap",
			lines: []string{
				"&&&&&&&&&&",
				"~~......@@",
				"!!.1111.##",
				"$$......^^",
				"**********",
			},
			part: &part{
				pos: pos{line: 2, start: 3, end: 6},
			},
			expectedOK: false,
		},
		{
			name: "not part number at start of line",
			lines: []string{
				"..........",
				"1111.@@@##",
				"..........",
			},
			part: &part{
				pos: pos{line: 1, start: 0, end: 3},
			},
			expectedOK: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, ok := isPartNumber(tc.lines, tc.part)

			assert.Equal(t, tc.expected, actual)
			assert.Equal(t, tc.expectedOK, ok)
		})
	}
}

func TestGears(t *testing.T) {
	testCases := []struct {
		name     string
		input    []*part
		expected []*part
	}{
		{
			name: "no gears - no kind match",
			input: []*part{
				{val: 1, symbol: &symbol{kind: "@", pos: pos{line: 1, start: 3}}},
				{val: 2, symbol: &symbol{kind: "&", pos: pos{line: 1, start: 3}}},
			},
			expected: nil,
		},
		{
			name: "match",
			input: []*part{
				{val: 1, symbol: &symbol{kind: "*", pos: pos{line: 1, start: 3}}},
				{val: 2, symbol: &symbol{kind: "*", pos: pos{line: 1, start: 3}}},
			},
			expected: []*part{
				{val: 1, symbol: &symbol{kind: "*", pos: pos{line: 1, start: 3}}},
				{val: 2, symbol: &symbol{kind: "*", pos: pos{line: 1, start: 3}}},
			},
		},
		{
			name: "mix of matches and unmatched",
			input: []*part{
				{val: 1, symbol: &symbol{kind: "*", pos: pos{line: 0, start: 2}}},
				{val: 2, symbol: &symbol{kind: "*", pos: pos{line: 2, start: 5}}},
				{val: 3, symbol: &symbol{kind: "*", pos: pos{line: 0, start: 2}}},
				{val: 5, symbol: &symbol{kind: "*", pos: pos{line: 1, start: 3}}},
				{val: 4, symbol: &symbol{kind: "*", pos: pos{line: 1, start: 3}}},
			},
			expected: []*part{
				{val: 1, symbol: &symbol{kind: "*", pos: pos{line: 0, start: 2}}},
				{val: 3, symbol: &symbol{kind: "*", pos: pos{line: 0, start: 2}}},
				{val: 5, symbol: &symbol{kind: "*", pos: pos{line: 1, start: 3}}},
				{val: 4, symbol: &symbol{kind: "*", pos: pos{line: 1, start: 3}}},
			},
		},
		{
			name: "three matches",
			input: []*part{
				{val: 1, symbol: &symbol{kind: "*", pos: pos{line: 0, start: 1}}},
				{val: 2, symbol: &symbol{kind: "*", pos: pos{line: 0, start: 1}}},
				{val: 3, symbol: &symbol{kind: "*", pos: pos{line: 0, start: 1}}},
			},
			expected: nil,
		},
		{
			name: "mix",
			input: []*part{
				{val: 1, symbol: &symbol{kind: "*", pos: pos{line: 0, start: 2}}},
				{val: 2, symbol: &symbol{kind: "*", pos: pos{line: 2, start: 5}}},
				{val: 3, symbol: &symbol{kind: "*", pos: pos{line: 0, start: 2}}},
				{val: 5, symbol: &symbol{kind: "*", pos: pos{line: 1, start: 3}}},
				{val: 4, symbol: &symbol{kind: "*", pos: pos{line: 1, start: 3}}},
			},
			expected: []*part{
				{val: 1, symbol: &symbol{kind: "*", pos: pos{line: 0, start: 2}}},
				{val: 3, symbol: &symbol{kind: "*", pos: pos{line: 0, start: 2}}},
				{val: 5, symbol: &symbol{kind: "*", pos: pos{line: 1, start: 3}}},
				{val: 4, symbol: &symbol{kind: "*", pos: pos{line: 1, start: 3}}},
			},
		},
		{
			name: "tens",
			input: []*part{
				{val: 10, symbol: &symbol{kind: "*", pos: pos{line: 1, start: 3}}},
				{val: 99, symbol: &symbol{kind: "*", pos: pos{line: 1, start: 8}}},
				{val: 10, symbol: &symbol{kind: "*", pos: pos{line: 1, start: 3}}},
				{val: 10, symbol: &symbol{kind: "*", pos: pos{line: 1, start: 3}}},
				{val: 1, symbol: &symbol{kind: "*", pos: pos{line: 1, start: 8}}},
			},
			expected: []*part{
				{val: 99, symbol: &symbol{kind: "*", pos: pos{line: 1, start: 8}}},
				{val: 1, symbol: &symbol{kind: "*", pos: pos{line: 1, start: 8}}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.ElementsMatch(t, tc.expected, gears(tc.input))
		})
	}
}
