package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	_zeroRune byte = 48
	_nineRune byte = 57
)

var (
	numSpellings = []map[string]string{
		{"zero": "0"},
		{"one": "1"},
		{"two": "2"},
		{"three": "3"},
		{"four": "4"},
		{"five": "5"},
		{"six": "6"},
		{"seven": "7"},
		{"eight": "8"},
		{"nine": "9"},
	}
)

func main() {
	var inputLoc string

	flag.StringVar(&inputLoc, "loc", "", "specify location of input file")
	flag.Parse()

	if inputLoc == "" {
		log.Fatalf("no input file found")
	}

	f, err := os.ReadFile(inputLoc)
	if err != nil {
		log.Fatalf("opening file: %s", err)
	}

	fmt.Printf("calibration value of %q is %d\n", inputLoc, calibration(string(f)))
}

// calibration attempts to determine a calibration rate from a garbled series of lines, summing
// together all numbers found across the lines.
func calibration(input string) int {
	inputScanner := bufio.NewScanner(strings.NewReader(input))

	var total int
	for inputScanner.Scan() {
		foundNums := parseNumbers(inputScanner.Text())
		if len(foundNums) == 0 {
			continue
		}

		total += calibrationValue(foundNums)
	}

	return total
}

// parseNumbers returns a collection of numbers if any are found within the given line.
func parseNumbers(input string) []string {
	// handle empty or whitespace only
	line := strings.TrimSpace(input)
	if line == "" {
		return nil
	}

	var nums []string
	for i := 0; i < len(line); i++ {
		if line[i] >= _zeroRune && line[i] <= _nineRune {
			// if we find a digit, add it
			nums = append(nums, string(line[i]))
			continue
		}

		// check for spelled numbers
		for _, ns := range numSpellings {
			for k, v := range ns {
				if !strings.HasPrefix(input[i:], k) {
					break
				}

				nums = append(nums, v)
				// some spelled numbers overlap so setting index back by two
				// moves us back one back from the end of the found spelling
				// examples: twone & threeight
				i += len(k) - 2
			}
		}
	}

	if len(nums) == 0 {
		return nil
	}

	return nums
}

// calibrationValue is responsible for returning a two-digit value to be used in calibration. If it
// is unable to construct a suitable number or the number would be 00, it will return 0 instead.
func calibrationValue(nums []string) int {
	var found string

	switch len(nums) {
	case 0:
		return 0
	case 1:
		// duplicate single numbers
		found += nums[0]
		found += nums[0]
	default:
		found += nums[0]
		found += nums[len(nums)-1]
	}

	out, err := strconv.Atoi(found)
	if err != nil {
		return 0
	}

	return out
}
