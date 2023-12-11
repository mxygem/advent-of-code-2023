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

	sum := partNumberSum(string(f))

	fmt.Printf("Part number total: %d\n", sum)
}

func partNumberSum(in string) int {
	inputScanner := bufio.NewScanner(strings.NewReader(in))

	var lines []string
	for inputScanner.Scan() {
		lines = append(lines, strings.TrimSpace(inputScanner.Text()))
	}

	numberPositions := numberPositions(lines)
	partNumbers := partNumbers(lines, numberPositions)

	return sum(partNumbers)
}

func numberPositions(lines []string) [][][]int {
	if len(lines) == 0 {
		return nil
	}

	allLines := make([][][]int, 0, len(lines))
	for _, line := range lines {
		var nums [][]int
		for i := 0; i < len(line); i++ {
			start, end := numberPos(line[i:])
			// no numbers found
			if start < 0 || end < 0 {
				break
			}

			nums = append(nums, []int{start + i, end + i})
			i = end + i
		}

		allLines = append(allLines, nums)
	}

	return allLines
}

func numberPos(line string) (int, int) {
	start := -1
	for i, c := range line {
		if c < 48 || c > 58 {
			// finding number - symbol found, start unset - continue
			if start < 0 {
				continue
			}
			// number end found - symbol found, start set - return start and current index - 1
			return start, i - 1
		}
		if start < 0 {
			// number found - number found, start unset - set start to current index
			start = i
		}
		// looking for number end - number found, start set - continue
	}

	if start < 0 {
		return -1, -1
	}

	// number found and ran until end of line
	return start, len(line) - 1
}

func isPartNumber(lines []string, lineNum int, start, end int) bool {
	if lineNum >= len(lines) || start < 0 || end > len(lines[lineNum]) {
		return false
	}

	line := lines[lineNum]

	// setup positions to index behavior and diags
	cs := start - 1
	if cs < 0 {
		cs = 0
	}

	ce := end + 2
	if ce >= len(line) {
		ce = len(line)
	}

	// check previous line if not first
	if lineNum > 0 {
		for _, c := range lines[lineNum-1][cs:ce] {
			if c != 46 && (c < 48 || c > 58) {
				return true
			}
		}
	}

	// check same line
	for _, c := range line[cs:ce] {
		if c != 46 && (c < 48 || c > 58) {
			return true
		}
	}

	// check next row if not last
	if lineNum < len(lines)-1 {
		for _, c := range lines[lineNum+1][cs:ce] {
			if c != 46 && (c < 48 || c > 58) {
				return true
			}
		}
	}

	return false
}

func partNumbers(lines []string, positions [][][]int) []int {
	var pns []int

	for i, linePositions := range positions {
		for _, pos := range linePositions {
			if !isPartNumber(lines, i, pos[0], pos[1]) {
				continue
			}

			n, err := strconv.Atoi(lines[i][pos[0] : pos[1]+1])
			if err != nil {
				fmt.Printf("WARN - could not convert %q to int: %v\n", lines[i][pos[0]:pos[1]+1], err)
				continue
			}

			pns = append(pns, n)
		}
	}

	return pns
}

func sum(pns []int) int {
	var sum int
	for _, pn := range pns {
		sum += pn
	}
	return sum
}
