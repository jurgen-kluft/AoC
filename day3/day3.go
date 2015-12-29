package main

import (
	"bufio"
	"fmt"
	"os"
)

func iterateOverCharsInTextFile(filename string, action func(rune)) {
	// Open the file.
	f, _ := os.Open(filename)
	defer f.Close()

	// Create a new Scanner for the file.
	scanner := bufio.NewScanner(f)

	// Loop over all lines in the file and print them.
	for scanner.Scan() {
		line := scanner.Text()
		for _, c := range line {
			// do something with c
			action(c)
		}
	}
}

func handoutPresentToCurrentLocation(X int, Y int, worldmap map[uint64]int) {
	var loc = (uint64(X) << 32) | (uint64(Y) & 0xFFFFFFFF)
	if value, ok := worldmap[loc]; ok {
		worldmap[loc] = value + 1
		fmt.Printf("Loc[X,Y] %v : %v = %v presents (loc: 0x%x)\n", X, Y, value+1, loc)
	} else {
		worldmap[loc] = 1
		fmt.Printf("Loc[X,Y] %v : %v = %v present (loc: 0x%x)\n", X, Y, 1, loc)
	}
}

func howManyHousesReceiveOnePresent(filename string) (numberOfHousesWithOnlyOnePresent int, ok bool) {
	numberOfHousesWithOnlyOnePresent = 0

	X := 0
	Y := 0
	worldmap := make(map[uint64]int)
	handoutPresentToCurrentLocation(X, Y, worldmap)

	computator := func(move rune) {
		switch move {
		case '<':
			X++
		case '>':
			X--
		case '^':
			Y++
		case 'v':
			Y--
		default:
			return
		}
		handoutPresentToCurrentLocation(X, Y, worldmap)
	}
	iterateOverCharsInTextFile(filename, computator)

	for _, nrpresents := range worldmap {
		if nrpresents >= 1 {
			numberOfHousesWithOnlyOnePresent++
		}
	}

	ok = true
	return
}

func main() {
	var numberOfHousesWithOnlyOnePresent, ok = howManyHousesReceiveOnePresent("input.text")
	if ok {
		fmt.Printf("There are %v houses that only received one present", numberOfHousesWithOnlyOnePresent)
	} else {
		fmt.Printf("Could not process the input")
	}
}
