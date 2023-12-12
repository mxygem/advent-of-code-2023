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

	allParts := parts(lines)
	gears := gears(allParts)

	return calc(gears)
}

type part struct {
	val int
	pos
	symbol *symbol
}

type symbol struct {
	kind string
	pos
}

type pos struct {
	line       int
	start, end int
}

func parts(lines []string) []*part {
	if len(lines) == 0 {
		return nil
	}

	var posParts []*part
	for lp, line := range lines {
		for i := 0; i < len(line); i++ {
			start, end := numberPos(line[i:])
			// no numbers found
			if start < 0 || end < 0 {
				break
			}

			num, err := strconv.Atoi(line[start+i : end+i+1])
			if err != nil {
				continue
			}

			posPart := &part{
				val: num,
				pos: pos{
					line:  lp,
					start: start + i,
					end:   end + i,
				},
			}

			p, ok := isPartNumber(lines, posPart)
			if !ok {
				continue
			}

			posParts = append(posParts, p)
			i = end + i
		}
	}

	return posParts
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

func isPartNumber(lines []string, possPart *part) (*part, bool) {
	if possPart == nil {
		return nil, false
	}
	if possPart.line >= len(lines) || possPart.start < 0 || possPart.end > len(lines[possPart.line]) {
		return nil, false
	}

	pp := *possPart
	line := lines[pp.line]

	// setup positions to index behavior and diags
	cs := pp.start - 1
	if cs < 0 {
		cs = 0
	}

	ce := pp.end + 2
	if ce >= len(line) {
		ce = len(line)
	}

	// check previous line if not first
	if pp.line > 0 {
		for i, c := range lines[pp.line-1][cs:ce] {
			if c != 46 && (c < 48 || c > 58) {
				pp.symbol = &symbol{
					kind: string(c),
					pos: pos{
						line:  pp.line - 1,
						start: cs + i,
					},
				}

				return &pp, true
			}
		}
	}

	// check same line
	for i, c := range line[cs:ce] {
		if c != 46 && (c < 48 || c > 58) {
			pp.symbol = &symbol{
				kind: string(c),
				pos: pos{
					line:  pp.line,
					start: cs + i,
				},
			}

			return &pp, true
		}
	}

	// check next row if not last
	if pp.line < len(lines)-1 {
		for i, c := range lines[pp.line+1][cs:ce] {
			if c != 46 && (c < 48 || c > 58) {
				pp.symbol = &symbol{
					kind: string(c),
					pos: pos{
						line:  pp.line + 1,
						start: cs + i,
					},
				}

				return &pp, true
			}
		}
	}

	return nil, false
}

func gears(parts []*part) []*part {
	var gs []*part
	fs := map[pos]int{}

	for _, p := range parts {
		if p.symbol == nil || p.symbol.kind != "*" {
			continue
		}

		for _, pg := range parts {
			if p == pg {
				continue
			}

			if pg.symbol.kind != "*" ||
				p.symbol.kind != pg.symbol.kind ||
				p.symbol.line != pg.symbol.line ||
				p.symbol.pos.start != pg.symbol.pos.start {
				continue
			}

			if fs[p.symbol.pos] == 0 {
				gs = append(gs, p, pg)
			}
			fs[p.symbol.pos]++
		}
	}

	// remove invalid matches
	var dd []*part
	for k, f := range fs {
		if f > 2 {
			continue
		}

		for _, p := range gs {
			if p.symbol.line != k.line || p.symbol.pos.start != k.start {
				continue
			}

			dd = append(dd, p)
		}
	}
	if len(dd) == 0 {
		return nil
	}

	return dd
}

func calc(ps []*part) int {
	var sum int
	for i, p := range ps {
		if i == 0 || i%2 == 0 {
			continue
		}

		sum += p.val * ps[i-1].val
	}
	return sum
}
