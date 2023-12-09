package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseGames(t *testing.T) {
	testCases := []struct {
		name        string
		gamesIn     string
		expected    int
		expectedErr error
	}{
		{
			name:        "no input",
			gamesIn:     "",
			expectedErr: fmt.Errorf("no input received"),
		},
		{
			name:     "two games",
			gamesIn:  "Game 1: 1 red, 1 blue, 1 green; 11 green, 11 blue, 11 red\nGame 2: 2 red, 2 blue, 2 green",
			expected: 3,
		},
		{
			name:     "one game not possible",
			gamesIn:  "Game 3: 1 red, 1 blue, 1 green\nGame 4: 14 red, 2 blue, 2 green",
			expected: 3,
		},
		{
			name:     "sum of second game not possible",
			gamesIn:  "Game 5: 1 red, 1 blue, 1 green\nGame 10: 3 red, 7 blue, 2 green; 4 red, 8 blue, 3 green\nGame 15: 3 red, 3 blue, 3 green",
			expected: 20,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := parseGames(tc.gamesIn)

			assert.Equal(t, tc.expected, actual)
			checkErr(t, tc.expectedErr, err)
		})
	}
}

func TestParseGame(t *testing.T) {
	testCases := []struct {
		name        string
		gameIn      string
		expected    *game
		expectedErr error
	}{
		{
			name:     "no input",
			gameIn:   "",
			expected: nil,
		},
		{
			name:     "no sets found",
			gameIn:   "Game 1:",
			expected: &game{id: 1},
		},
		{
			name:     "single set",
			gameIn:   "Game 2: 1 red, 2 blue, 3 green",
			expected: &game{id: 2, sets: []*set{{red: 1, blue: 2, green: 3}}},
		},
		{
			name:   "three sets - all colors in order",
			gameIn: "Game 3: 1 red, 2 blue, 3 green; 4 red, 5 blue, 6 green; 7 red, 8 blue, 9 green",
			expected: &game{id: 3, sets: []*set{
				{red: 1, blue: 2, green: 3},
				{red: 4, blue: 5, green: 6},
				{red: 7, blue: 8, green: 9},
			}},
		},
		{
			name:   "trailing comma after second set",
			gameIn: "Game 4: 8 blue, 9 red, 7 green; 5 blue, 4 green, 6 red,; 1 green, 2 blue, 3 red",
			expected: &game{id: 4, sets: []*set{
				{red: 9, blue: 8, green: 7},
				{red: 6, blue: 5, green: 4},
				{red: 3, blue: 2, green: 1},
			}},
		},
		{
			name:   "five sets - single color in each",
			gameIn: "Game 5: 1 green; 2 blue; 3 red; 4 red; 6 blue",
			expected: &game{id: 5, sets: []*set{
				{green: 1},
				{blue: 2},
				{red: 3},
				{red: 4},
				{blue: 6},
			}},
		},
		{
			name:   "empty set within multiples",
			gameIn: "Game 6: 3 red, 3 blue, 3 green;; 4 red, 4 blue, 4 green",
			expected: &game{id: 6, sets: []*set{
				{red: 3, blue: 3, green: 3},
				{red: 4, blue: 4, green: 4},
			}},
		},
		{
			name:   "no space in game name",
			gameIn: "Game7: 5 red, 5 blue, 5 green",
			expected: &game{id: 7, sets: []*set{
				{red: 5, blue: 5, green: 5},
			}},
		},
		{
			name:     "no spaces in sets",
			gameIn:   "Game 8:1red,1blue,1green;2red,2blue,2green",
			expected: &game{id: 8},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := parseGame(tc.gameIn)

			assert.Equal(t, tc.expected, actual)
			checkErr(t, tc.expectedErr, err)
		})
	}
}

func TestGameID(t *testing.T) {
	testCases := []struct {
		name        string
		headerIn    string
		expected    int
		expectedErr error
	}{
		{
			name:     "empty",
			headerIn: "",
			expected: 0,
		},
		{
			name:     "expected formatting",
			headerIn: "Game 1",
			expected: 1,
		},
		{
			name:     "no space",
			headerIn: "Game2",
			expected: 2,
		},
		{
			name:        "only id",
			headerIn:    "2",
			expectedErr: fmt.Errorf(`invalid game title found: "2"`),
		},
		{
			name:        "id is not number",
			headerIn:    "Game A",
			expectedErr: fmt.Errorf(`converting game id " a" to int: strconv.Atoi: parsing "a": invalid syntax`),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := gameID(tc.headerIn)

			assert.Equal(t, tc.expected, actual)
			checkErr(t, tc.expectedErr, err)
		})
	}
}

func TestParseSets(t *testing.T) {
	testCases := []struct {
		name     string
		setsIn   string
		expected []*set
	}{
		{
			name:   "single set",
			setsIn: " 1 red, 2 blue, 3 green",
			expected: []*set{
				{red: 1, blue: 2, green: 3},
			},
		},
		{
			name:   "multiple sets in order",
			setsIn: " 1 red, 2 blue, 3 green; 4 red, 5 blue, 6 green; 7 red, 8 blue, 9 green",
			expected: []*set{
				{red: 1, blue: 2, green: 3},
				{red: 4, blue: 5, green: 6},
				{red: 7, blue: 8, green: 9},
			},
		},
		{
			name:   "trailing comma after second set",
			setsIn: " 8 blue, 9 red, 7 green; 5 blue, 4 green, 6 red,; 1 green, 2 blue, 3 red",
			expected: []*set{
				{red: 9, blue: 8, green: 7},
				{red: 6, blue: 5, green: 4},
				{red: 3, blue: 2, green: 1},
			},
		},
		{
			name:   "empty sets are not returned",
			setsIn: "1 red;; 1 blue;; 1 green, 2 blue, 3 red",
			expected: []*set{
				{red: 1},
				{blue: 1},
				{red: 3, blue: 2, green: 1},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, parseSets(tc.setsIn))
		})
	}
}

func TestParseSet(t *testing.T) {
	testCases := []struct {
		name     string
		setIn    string
		expected *set
	}{
		{
			name:     "empty",
			setIn:    "",
			expected: nil,
		},
		{
			name:     "expected formatting and in order",
			setIn:    " 1 red, 2 blue, 3 green",
			expected: &set{red: 1, blue: 2, green: 3},
		},
		{
			name:     "expected formatting and in order",
			setIn:    "4 blue, 5 green, 6 red",
			expected: &set{red: 6, blue: 4, green: 5},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := parseSet(tc.setIn)

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestPossibleGame(t *testing.T) {
	testCases := []struct {
		name     string
		totals   set
		check    *game
		expected bool
	}{
		{
			name:     "totals zero",
			totals:   set{},
			expected: false,
		},
		{
			name:     "red too high for only set",
			totals:   set{red: 10},
			check:    &game{sets: []*set{{red: 20}}},
			expected: false,
		},
		{
			name:     "blue too high for last set",
			totals:   set{blue: 10},
			check:    &game{sets: []*set{{blue: 5}, {blue: 10}, {blue: 20}}},
			expected: false,
		},
		{
			name:     "sum of green across games too high",
			totals:   set{green: 20},
			check:    &game{sets: []*set{{green: 5}, {green: 10}, {green: 10}}},
			expected: false,
		},
		{
			name:   "all values valid",
			totals: set{red: 15, blue: 15, green: 15},
			check: &game{sets: []*set{
				{red: 2, blue: 2, green: 2},
				{red: 3, blue: 3, green: 3},
				{red: 4, blue: 4, green: 4},
			}},
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, possibleGame(tc.totals, tc.check))
		})
	}
}

func checkErr(t *testing.T, expected, actual error) {
	if expected != nil {
		require.Error(t, actual)
		assert.Equal(t, expected.Error(), actual.Error())
	} else {
		assert.Nil(t, actual)
	}
}
