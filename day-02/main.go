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

	gs, err := parseGames(string(f))
	if err != nil {
		log.Fatalf("parsing games: %s", err)
	}

	fmt.Printf("found %d sum of IDs\n", gs)
}

const (
	_redMax   int = 12
	_blueMax  int = 14
	_greenMax int = 13
	_totalMax int = _redMax + _blueMax + _greenMax
)

type game struct {
	id   int
	sets []*set
}

type set struct {
	red, blue, green int
}

func parseGames(in string) (int, error) {
	if in == "" {
		return 0, fmt.Errorf("no input received")
	}

	inputScanner := bufio.NewScanner(strings.NewReader(in))
	var idTotal int
	var errs []error

	for inputScanner.Scan() {
		game, err := parseGame(inputScanner.Text())
		if err != nil {
			errs = append(errs, err)
			continue
		}

		// determine possible game here
		max := set{red: _redMax, blue: _blueMax, green: _greenMax}
		if !possibleGame(max, game) {
			continue
		}

		idTotal += game.id
	}

	if len(errs) > 0 {
		// update later to return errs
		for i, err := range errs {
			fmt.Printf("err %d:%v\n", i, err)
		}
	}

	return idTotal, nil
}

func parseGame(in string) (*game, error) {
	if in == "" {
		return nil, nil
	}

	ts := strings.Split(in, ":")
	if len(ts) < 2 {
		return nil, fmt.Errorf("not enough parts found. found %d", len(ts))
	}

	g := &game{}

	id, err := gameID(ts[0])
	if err != nil {
		return nil, fmt.Errorf("retrieving game id: %w", err)
	}
	g.id = id
	g.sets = parseSets(ts[1])

	return g, nil
}

func gameID(in string) (int, error) {
	if in == "" {
		return 0, nil
	}

	sp := strings.Split(strings.ToLower(in), "game")
	if len(sp) < 2 {
		return 0, fmt.Errorf("invalid game title found: %q", in)
	}

	id, err := strconv.Atoi(strings.TrimSpace(sp[1]))
	if err != nil {
		return 0, fmt.Errorf("converting game id %q to int: %w", sp[1], err)
	}

	return id, nil
}

func parseSets(in string) []*set {
	splitSets := strings.Split(in, ";")

	var sets []*set
	for _, s := range splitSets {
		set := parseSet(s)
		if set == nil {
			continue
		}

		sets = append(sets, set)
	}

	return sets
}

func parseSet(in string) *set {
	if in == "" {
		fmt.Println("found empty set")
		return nil
	}

	in = strings.TrimSpace(in)
	cs := strings.Split(in, ",")

	s := &set{}
	for _, c := range cs {
		sc := strings.Split(strings.TrimSpace(c), " ")
		if len(sc) < 2 {
			fmt.Printf("found invalid color spit: %q\n", sc)
			continue
		}

		color := strings.TrimSpace(sc[1])
		num := strings.TrimSpace(sc[0])

		n, err := strconv.Atoi(num)
		if err != nil {
			fmt.Printf("could not convert: %q to int: %v\n", n, err)
			continue
		}

		switch color {
		case "red":
			s.red = n
		case "blue":
			s.blue = n
		case "green":
			s.green = n
		}
	}

	if s.red == 0 && s.blue == 0 && s.green == 0 {
		return nil
	}

	return s
}

func possibleGame(totals set, check *game) bool {
	if totals.red == 0 && totals.blue == 0 && totals.green == 0 {
		return false
	}

	fmt.Printf("----- game %d -----\n", check.id)

	// var redTotal, blueTotal, greenTotal int
	for i, s := range check.sets {
		fmt.Printf("\tset %d:%+v\n", i, s)
		if s.red > totals.red || s.blue > totals.blue || s.green > totals.green {
			// fmt.Printf("\tsingle value too high - r:%d:%d, b:%d:%d, g:%d:%d\n", s.red, totals.red, s.blue, totals.blue, s.green, totals.green)
			return false
		}
		// redTotal += s.red
		// blueTotal += s.blue
		// greenTotal += s.green
	}

	// if redTotal > totals.red || blueTotal > totals.blue || greenTotal > totals.green {
	// 	fmt.Printf("\tgame total too high: - r:%d:%d, b:%d:%d, g:%d:%d\n", redTotal, totals.red, blueTotal, totals.blue, greenTotal, totals.green)
	// 	return false
	// }

	return true
}
