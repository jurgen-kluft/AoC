package main

/*

 */

import (
	"bufio"
	"fmt"
	"os"
)

func iterateOverLinesInTextFile(filename string, action func(string, int)) {
	// Open the file.
	f, _ := os.Open(filename)
	defer f.Close()

	// Create a new Scanner for the file.
	scanner := bufio.NewScanner(f)

	// Loop over all lines in the file and print them.
	lineNumber := 1
	for scanner.Scan() {
		line := scanner.Text()
		action(line, lineNumber)
		lineNumber++
	}
}

type eState int

const (
	cReadChar        eState = 2
	cEscapeOpen      eState = 3
	cEscapeHexDigit1 eState = 5
	cEscapeHexDigit2 eState = 6
)

func countNumberOfInMemoryChars(str string) (total int32, inmemory int32) {
	state := cReadChar

	total = 0
	inmemory = 0
	for _, char := range str {
		total++
		switch state {
		case cReadChar:
			if char == '"' {
				// Start or End Quote ?
			} else if char == '\\' {
				state = cEscapeOpen
			} else {
				inmemory++
			}
		case cEscapeOpen:
			if char == 'x' {
				state = cEscapeHexDigit1
			} else {
				inmemory++
				state = cReadChar
			}
		case cEscapeHexDigit1:
			state = cEscapeHexDigit2
		case cEscapeHexDigit2:
			state = cReadChar
			inmemory++
		}
	}
	return
}

func reEncodeString(str string) (newstr string) {
	newstr = ""
	escaped := false
	for _, char := range str {
		if escaped {
			escaped = false
			if char == 'x' {
				newstr = newstr + "\\\\" + string(char)
			} else {
				newstr = newstr + "\\\\\\" + string(char)
			}
		} else if char == '"' {
			newstr = newstr + "\"\\\""
		} else if char == '\\' {
			escaped = true
		} else {
			newstr = newstr + string(char)
		}
	}
	return
}

func evalLiteralsFromFile(filename string) (result int32) {
	total := int32(0)
	inmemory := int32(0)

	encodedTotal := int32(0)
	encodedInmemory := int32(0)

	computator := func(text string, line int) {
		t, i := countNumberOfInMemoryChars(text)
		fmt.Printf("Line(%d): total:%v, in-memory:%v", line, t, i)
		total += t
		inmemory += i

		text = reEncodeString(text)
		et, ei := countNumberOfInMemoryChars(text)
		fmt.Printf(", encoded: total:%v, in-memory:%v\n", et, ei)
		encodedTotal += et
		encodedInmemory += ei
	}
	iterateOverLinesInTextFile(filename, computator)

	result = encodedTotal - total
	return
}

func main() {
	result := evalLiteralsFromFile("input.text")
	fmt.Printf("The result is: %v", result)
}
