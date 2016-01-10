package main

/*
--- Day 19: Medicine for Rudolph ---

Rudolph the Red-Nosed Reindeer is sick! His nose isn't shining very brightly,
and he needs medicine.

Red-Nosed Reindeer biology isn't similar to regular reindeer biology; Rudolph
is going to need custom-made medicine. Unfortunately, Red-Nosed Reindeer
chemistry isn't similar to regular reindeer chemistry, either.

The North Pole is equipped with a Red-Nosed Reindeer nuclear fusion/fission
plant, capable of constructing any Red-Nosed Reindeer molecule you need.
It works by starting with some input molecule and then doing a series of
replacements, one per step, until it has the right molecule.

However, the machine has to be calibrated before it can be used. Calibration
involves determining the number of molecules that can be generated in one step
from a given starting point.

For example, imagine a simpler machine that supports only the following
replacements:

H => HO
H => OH
O => HH
Given the replacements above and starting with HOH, the following molecules
could be generated:

HOOH (via H => HO on the first H).
HOHO (via H => HO on the second H).
OHOH (via H => OH on the first H).
HOOH (via H => OH on the second H).
HHHH (via O => HH).
So, in the example above, there are 4 distinct molecules (not five, because
HOOH appears twice) after one replacement from HOH. Santa's favorite molecule,
HOHOHO, can become 7 distinct molecules (over nine replacements: six from H,
and three from O).

The machine replaces without regard for the surrounding characters.
For example, given the string H2O, the transition H => OO would result in OO2O.

Your puzzle input describes all of the possible replacements and, at the bottom,
the medicine molecule for which you need to calibrate the machine. How many
distinct molecules can be created after all the different ways you can do one
replacement on the medicine molecule?
*/

import (
	"bufio"
	"fmt"
	"os"
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

type molecule struct {
	sequence string
	replacer string
}

func (m *molecule) deserialize(str string) {
	// Example line:
	//      Al => ThF

	parts := strings.Split(str, "=>")
	m.sequence = strings.Trim(parts[0], " ")
	m.replacer = strings.Trim(parts[1], " ")
}

func printReplacements(r []*molecule) {
	for _, n := range r {
		fmt.Printf("%s => %s\n", n.sequence, n.replacer)
	}
}

func readInputFromFile(filename string) (medicine string, replacements []*molecule) {
	replacements = make([]*molecule, 0, 100)

	medicine = ""
	state := "replacements"
	computator := func(text string, line int) {
		if text == "" {
			state = "medicine"
		} else {
			if state == "replacements" {
				m := &molecule{sequence: "", replacer: ""}
				m.deserialize(text)
				replacements = append(replacements, m)
			} else {
				medicine = text
			}
		}
	}
	iterateOverLinesInTextFile(filename, computator)

	return
}

func buildMedicine(medicinestr string, replaceat int, mol *molecule) string {
	sl := len(mol.sequence)
	rl := len(mol.replacer)
	l := rl - sl
	scratch := make([]rune, len(medicinestr)+l)
	for i, c := range medicinestr {
		if i >= replaceat {
			if i >= replaceat+sl {
				scratch[i+l] = c
			}
		} else {
			scratch[i] = c
		}
	}
	for i, c := range mol.replacer {
		scratch[replaceat+i] = c
	}
	return string(scratch)
}

func computeMedicine(medicinestr string, replacements []*molecule) (distinct int) {
	distinct = 0

	medicines := make(map[string]int)

	for _, mol := range replacements {
		for mci := range medicinestr {
			found := true
			for msi := range mol.sequence {
				mcie := mci + msi
				if mcie == len(medicinestr) || mol.sequence[msi] != medicinestr[mcie] {
					found = false
					break
				}
			}
			if found {
				medicine := buildMedicine(medicinestr, mci, mol)
				if mc, exists := medicines[medicine]; exists {
					medicines[medicine] = mc + 1
				} else {
					distinct++
					medicines[medicine] = 0
				}
			}
		}
	}
	return
}

func main() {
	medicine, replacements := readInputFromFile("input.text")
	fmt.Println(medicine)
	printReplacements(replacements)
	distinct := computeMedicine(medicine, replacements)
	fmt.Printf("Found %v distinct medicine molecules\n", distinct)
}
