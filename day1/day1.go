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

func toWhichFloorSantaGoes(filename string) (floor int, basement int, ok bool) {
	floor = 0
	basement = -1
	computator := func(line string) {
		for pos, op := range line {
			switch op {
			case '(':
				floor++
			case ')':
				floor--
			}
			if floor == -1 && basement == -1 {
				basement = pos + 1
			}
		}
	}

	iterateOverLinesInTextFile(filename, computator)
	ok = true
	return
}

func main() {
	var floor, basement, ok = toWhichFloorSantaGoes("input.text")
	if ok {
		fmt.Printf("Santa went to the %v floor\n", floor)
		fmt.Printf("Santa first at the basement %v\n", basement)
	} else {
		fmt.Printf("Could not process the input")
	}
}
