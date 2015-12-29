package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

type sortRunes []rune

func (s sortRunes) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s sortRunes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortRunes) Len() int {
	return len(s)
}

func sortString(s string) string {
	r := []rune(s)
	sort.Sort(sortRunes(r))
	return string(r)
}

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

func isNiceString(str string) (nice bool) {
	nice = false

	// Count the number of vowels
	vowels := 0
	for _, char := range str {
		if char == 'a' || char == 'e' || char == 'i' || char == 'o' || char == 'u' {
			vowels++
		}
	}

	naughty := false
	if strings.Contains(str, "ab") || strings.Contains(str, "cd") || strings.Contains(str, "pq") || strings.Contains(str, "xy") {
		naughty = true
	}

	// It contains at least one letter that appears twice in a row
	twiceInRow := false
	prevchar := rune(0)
	for index, char := range str {
		if index > 0 {
			if prevchar == char {
				twiceInRow = true
				break
			}
		}
		prevchar = char
	}

	nice = vowels >= 3 && twiceInRow && !naughty
	return
}

func howManyNiceStringsInFile(filename string) (nice int) {
	nice = 0
	computator := func(line string) {
		if isNiceString(line) {
			nice++
		}
	}

	iterateOverLinesInTextFile(filename, computator)
	return
}

func main() {
	var howManyNiceStrings = howManyNiceStringsInFile("input.text")
	fmt.Printf("There are %v nice strings", howManyNiceStrings)
}
