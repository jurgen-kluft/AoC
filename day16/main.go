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
	"strconv"
	"strings"
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

type compound struct {
	name  string
	score int
}

const (
	eEvalEqual int = 0
	eEvalLess  int = -1
	eEvalMore  int = 1
)

var sCompoundEval = map[string]int{
	"children":    eEvalEqual,
	"cats":        eEvalMore,
	"samoyeds":    eEvalEqual,
	"pomeranians": eEvalLess,
	"akitas":      eEvalEqual,
	"vizslas":     eEvalEqual,
	"goldfish":    eEvalLess,
	"trees":       eEvalMore,
	"cars":        eEvalEqual,
	"perfumes":    eEvalEqual,
}

var sCompoundScore = map[string]int{
	"children":    1,
	"cats":        1,
	"samoyeds":    1000,
	"pomeranians": 1000,
	"akitas":      1000,
	"vizslas":     1000,
	"goldfish":    1000,
	"trees":       1,
	"cars":        1000,
	"perfumes":    1,
}

func isNameOfCompound(name string) bool {
	_, x := sCompoundScore[name]
	return x
}

func getScoreOfCompound(cmpName string, lhsValue int, rhsValue int) int {
	evalType := sCompoundEval[cmpName]
	score := 0
	switch evalType {
	case eEvalEqual:
		if lhsValue == rhsValue {
			score = 1
		}
	case eEvalLess:
		if lhsValue < rhsValue {
			score = 1
		}
	case eEvalMore:
		if lhsValue > rhsValue {
			score = 1
		}
	}
	score = sCompoundScore[cmpName] * score
	return score
}

type aunt struct {
	name      string
	index     int
	compounds map[string]int
}

func (r *aunt) print(eol string) {
	fmt.Printf("%s[%d]: ", r.name, r.index)
	for k, v := range r.compounds {
		fmt.Printf("{%s:%v}", k, v)
	}
	fmt.Printf(eol)
}

func printAll(aunts []*aunt) {
	for _, a := range aunts {
		a.print("\n")
	}
}

func (r *aunt) deserialize(str string) {
	// Example line:
	//      Sue 1: cars: 9, akitas: 3, goldfish: 0
	r.compounds = make(map[string]int)

	f := strings.FieldsFunc(str, func(r rune) bool {
		return r == ':' || r == ',' || r == ' '
	})

	r.name = f[0]
	r.index, _ = strconv.Atoi(f[1])

	doprint := false

	for i := 2; i < len(f)-1; i += 2 {
		if doprint {
			fmt.Printf("%s:%s ", f[i], f[i+1])
		}
		if isNameOfCompound(f[i]) {
			r.compounds[f[i]], _ = strconv.Atoi(f[i+1])
		}
	}
	if doprint {
		fmt.Print("\n")
	}
}

func readInputFromFile(filename string) (tofind *aunt, aunts []*aunt) {
	aunts = make([]*aunt, 0)

	computator := func(text string, line int) {
		if line == 1 {
			r := new(aunt)
			r.deserialize(text)
			tofind = r
		} else {
			r := new(aunt)
			r.deserialize(text)
			aunts = append(aunts, r)
		}
	}
	iterateOverLinesInTextFile(filename, computator)

	return
}

func computeScore(tofind *aunt, match *aunt) int {
	score := 0
	for n1, v1 := range tofind.compounds {
		if v2, x2 := match.compounds[n1]; x2 {
			if n1 == "trees" {

			} else if n1 == "cats" {

			} else {
				if v1 == v2 {
					score += getScoreOfCompound(n1, v1, v2)
				}
			}
		}
	}
	return score
}

func findBestMatchingAunt(tofind *aunt, aunts []*aunt) int {
	bestMatch := 0
	bestScore := 0
	for i, a := range aunts {
		score := computeScore(tofind, a)
		if score > bestScore {
			bestScore = score
			bestMatch = aunts[i].index
		}
	}
	return bestMatch
}

func main() {
	tofind, aunts := readInputFromFile("input.text")
	//tofind.print("\n")
	//printAll(aunts)
	match := findBestMatchingAunt(tofind, aunts)
	fmt.Printf("Our Aunt is number %v\n", match)
}
