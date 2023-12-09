package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
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

	_ = partNumberSum(string(f))
}

func partNumberSum(in string) int {
	inputScanner := bufio.NewScanner(strings.NewReader(in))
	var partNumberSum int

	var lines []string
	for inputScanner.Scan() {
		lines = append(lines, inputScanner.Text())
	}

	numberPositions := partNumberPositions(lines)
	fmt.Printf("read %d lines and found %d lines of numbers\n", len(lines), len(numberPositions))
	for i, l := range numberPositions {
		fmt.Printf("line: %d\n\t%s\n\t%+v\n", i, lines[i], l)
	}

	return partNumberSum
}

func partNumberPositions(lines []string) [][][]int {
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

	//check previous line
	line := lines[lineNum]

	// setup positions -1 and +1 from number to check
	cs := start - 1
	if cs < 0 {
		cs = 0
	}

	ce := end + 1
	if ce >= len(line) {
		ce = len(line) - 1
	}

	// check previous line
	if lineNum > 0 {
		cs--
		if cs < 0 {
			cs = 0
		}
		for _, c := range lines[lineNum-1][cs:ce] {
			if c != 46 && (c < 48 || c > 58) {
				return true
			}
		}
	}

	// check same row before and after
	fmt.Printf("same row cs: %d ce: %d\n", cs, ce)
	fmt.Printf("line[cs]: %d:%s line[ce]: %d:%s\n", line[cs], string(line[cs]), line[ce], string(line[ce]))
	if (line[cs] != 46 && (line[cs] < 48 || line[cs] > 58)) || (line[ce] != 46 && (line[ce] < 48 || line[ce] > 58)) {
		return true
	}

	// check next row if not last
	// check previous line
	fmt.Printf("lineNum: %d len(lines)-1: %d\n", lineNum, len(lines)-1)
	if lineNum < len(lines)-1 {
		ce++
		if ce >= len(line) {
			ce = len(line) - 1
		}

		for _, c := range lines[lineNum+1][cs:ce] {
			if c != 46 && (c < 48 || c > 58) {
				return true
			}
		}
	}

	return false
}
