package main

/*
--- Day 16: Aunt Sue ---

Your Aunt Sue has given you a wonderful gift, and you'd like to send her a
thank you card. However, there's a small problem: she signed it "From, Aunt Sue".

You have 500 Aunts named "Sue".

So, to avoid sending the card to the wrong person, you need to figure out which
Aunt Sue (which you conveniently number 1 to 500, for sanity) gave you the gift.
You open the present and, as luck would have it, good ol' Aunt Sue got you a
My First Crime Scene Analysis Machine! Just what you wanted. Or needed, as the
case may be.

The My First Crime Scene Analysis Machine (MFCSAM for short) can detect a few
specific compounds in a given sample, as well as how many distinct kinds of
those compounds there are. According to the instructions, these are what the
MFCSAM can detect:

children, by human DNA age analysis.
cats. It doesn't differentiate individual breeds.
Several seemingly random breeds of dog: samoyeds, pomeranians, akitas, and vizslas.
goldfish. No other kinds of fish.
trees, all in one group.
cars, presumably by exhaust or gasoline or something.
perfumes, which is handy, since many of your Aunts Sue wear a few kinds.
In fact, many of your Aunts Sue have many of these. You put the wrapping from
the gift into the MFCSAM. It beeps inquisitively at you a few times and then
prints out a message on ticker tape:

children: 3
cats: 7
samoyeds: 2
pomeranians: 3
akitas: 0
vizslas: 0
goldfish: 5
trees: 3
cars: 2
perfumes: 1
You make a list of the things you can remember about each Aunt Sue.
Things missing from your list aren't zero - you simply don't remember the
value.

What is the number of the Sue that got you the gift?
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

type aunt struct {
	name      string
	index     int
	compounds map[string]int
}

func printAll(aunts []*aunt) {

}

func (r *aunt) deserialize(str string) {
	// Example line:
	//      Sue 1: cars: 9, akitas: 3, goldfish: 0
	r.compounds = make(map[string]int)
	fmt.Sscanf(str, "%s %d", r.name, r.index)
}

func readInputFromFile(filename string) (aunts []*aunt) {
	aunts = make([]*aunt, 0)

	computator := func(text string, line int) {
		r := new(aunt)
		r.deserialize(text)
		aunts = append(aunts, r)
	}
	iterateOverLinesInTextFile(filename, computator)

	return
}

func main() {
	aunts := readInputFromFile("input.text")
	printAll(aunts)

}
