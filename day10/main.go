package main

/*
--- Day 10, findShortestRouteFromFile(filename string)Elves Look, Elves Say ---

Today, the Elves are playing a game called look-and-say. They take turns making
sequences by reading aloud the previous sequence and using that reading as the
next sequence. For example, 211 is read as "one two, two ones", which becomes
1221 (1 2, 2 1s).

Look-and-say sequences are generated iteratively, using the previous value as
input for the next step. For each step, take the previous value, and replace
each run of digits (like 111) with the number of digits (3) followed by the
digit itself (1).

For example:

1 becomes 11 (1 copy of digit 1).
11 becomes 21 (2 copies of digit 1).
21 becomes 1211 (one 2 followed by one 1).
1211 becomes 111221 (one 1, one 2, and two 1s).
111221 becomes 312211 (three 1s, two 2s, and one 1).
Starting with the digits in your puzzle input, apply this process 40 times.
What is the length of the result?

Your puzzle input is 1113122113.

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

func applySequence(seq []rune) (sequence []rune) {
	count := 0
	prevchar := rune(0)

	sequence = make([]rune, 0, len(seq)*2)
	for pos, char := range seq {
		isSequenceBreak := char != prevchar
		if pos > 0 && isSequenceBreak {
			//sequence = sequence + strconv.Itoa(count) + string(prevchar)
			sequence = append(sequence, rune('0')+rune(count))
			sequence = append(sequence, prevchar)
			count = 0
		}
		prevchar = char
		count++
	}
	if prevchar != 0 {
		sequence = append(sequence, rune('0')+rune(count))
		sequence = append(sequence, prevchar)
	}
	return
}

func decodeSequence(seq string) (count int, sequence string) {
	fmt.Sscanf(seq, "%d*%s", &count, &sequence)
	return
}

func stringToRuneArray(str string) (runes []rune) {
	runes = []rune(str)
	return
}

func fromFile(filename string, linenumber int) (finalSequence string) {
	answer := ""
	computator := func(text string, line int) {
		if line == linenumber {
			count, sequence := decodeSequence(text)
			runes := stringToRuneArray(sequence)
			for i := 0; i < count; i++ {
				runes = applySequence(runes)
				fmt.Printf("Iteration(%d) = %d\n", i, len(runes))
			}
			fmt.Printf("  Final ==> %d\n", len(runes))
			answer = string(runes)
		}
	}
	iterateOverLinesInTextFile(filename, computator)

	finalSequence = answer
	return
}

func main() {
	result := fromFile("input.text", 2)
	fmt.Printf("The result is : %s, %v\n", result, len(result))
}
