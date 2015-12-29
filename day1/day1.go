package main

import (
	"bufio"
	"fmt"
	"os"
)

func iterateOverLinesInTextFile(filename string, action func(string)) {
	// Open the file.
	f, _ := os.Open(filename)
	defer f.Close()

	// Create a new Scanner for the file.
	scanner := bufio.NewScanner(f)

	// Loop over all lines in the file and print them.
	for scanner.Scan() {
		line := scanner.Text()
		action(line)
	}
}

func toWhichFloorSantaGoes(filename string) (floor int, ok bool) {
	floor = 0
	computator := func(line string) {
		for _, op := range line {
			switch op {
			case '(':
				floor++
			case ')':
				floor--
			}
		}
	}

	iterateOverLinesInTextFile(filename, computator)
	ok = true
	return
}

func main() {
	var floor, ok = toWhichFloorSantaGoes("input.text")
	if ok {
		fmt.Printf("Santa went to the %v floor", floor)
	} else {
		fmt.Printf("Could not process the input")
	}
}
