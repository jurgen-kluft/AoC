package main

import (
    "bufio"
    "fmt"
    "os"
)

func ToWhichFloorSantaGoes(filename string) (floor int, ok bool) {
	// Open the file.
    f, _ := os.Open(filename)
    defer f.Close()

    // Create a new Scanner for the file.
    scanner := bufio.NewScanner(f)

    // Initial floor is 0
    floor = 0

    // Loop over all lines in the file and print them.
    for scanner.Scan() {
		line := scanner.Text()
		for _, op := range line {
			switch op {
			case '(':
				floor += 1
			case ')':
				floor -= 1
			}
		}
    }

    ok = true
	return
}

func main() {
	var floor, ok = ToWhichFloorSantaGoes("input.text")
	if (ok) {
		fmt.Printf("Santa went to the %v floor", floor)
	} else {
		fmt.Printf("Could not process the input")
	}
}