package main

/*
--- Day 17: No Such Thing as Too Much ---

The elves bought too much eggnog again - 150 liters this time.
To fit it all into your refrigerator, you'll need to move it into smaller
containers. You take an inventory of the capacities of the available containers.

For example, suppose you have containers of size 20, 15, 10, 5, and 5 liters.
If you need to store 25 liters, there are four ways to do it:

15 and 10
20 and 5 (the first 5)
20 and 5 (the second 5)
15, 5, and 5

Filling all containers entirely, how many different combinations of containers
can exactly fit all 150 liters of eggnog?

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

func deserializeSize(str string) int {
	// Example line:
	//      11
	integer, _ := strconv.Atoi(str)
	return integer
}
func deserializeTotal(str string) int {
	// Example line:
	//      total: 11
	f := strings.FieldsFunc(str, func(r rune) bool {
		return r == ':'
	})

	total, _ := strconv.Atoi(f[1])
	return total
}

func readInputFromFile(filename string) (total int, sizes []int) {
	sizes = make([]int, 0, 64)

	computator := func(text string, line int) {
		if line == 1 {
			total = deserializeTotal(text)
		} else {
			size := deserializeSize(text)
			sizes = append(sizes, size)
		}
	}
	iterateOverLinesInTextFile(filename, computator)

	return
}

type largenumber struct {
	nset  int
	nbits int
	bits  []byte
}

func newLargeNumber(numbits int) *largenumber {
	l := &largenumber{nset: 0, nbits: numbits, bits: make([]byte, numbits, numbits)}
	return l
}

func (l *largenumber) reset() {
	l.nset = 0
	l.nbits = len(l.bits)
	for i := 0; i < l.nbits; i++ {
		l.bits[i] = 0
	}
}

func (l *largenumber) equals(rhs *largenumber) bool {
	if l.nbits != rhs.nbits {
		return false
	}
	for i := 0; i < l.nbits; i++ {
		if l.bits[i] != rhs.bits[i] {
			return false
		}
	}
	return true
}

func (l *largenumber) assign(rhs *largenumber) {
	l.nbits = rhs.nbits
	l.nset = rhs.nset
	for i := 0; i < l.nbits; i++ {
		l.bits[i] = rhs.bits[i]
	}
}

func (l *largenumber) increment() bool {
	n := l.nbits - 1
	for n >= 0 && l.bits[n] == 1 {
		l.bits[n] = 0
		l.nset--
		n--
	}
	if n >= 0 {
		l.nset++
		l.bits[n] = 1
	}
	return l.nset < l.nbits
}

func (l *largenumber) print() {
	for _, b := range l.bits {
		if b == 1 {
			fmt.Print("1")
		} else {
			fmt.Print("0")
		}
	}
	fmt.Print("\n")
}

// Compute total number of combinations
func compute(total int, sizes []int) (combinations int, minimumcount int) {
	number := newLargeNumber(len(sizes))

	minimum := len(sizes)
	minimumcount = 0

	combinations = 0
	current := newLargeNumber(len(sizes))
	computed := newLargeNumber(len(sizes))
	for true {
		sum := 0

		computed.reset()
		computed.nbits = 0
		for index, bit := range number.bits {
			computed.bits[computed.nbits] = bit
			computed.nbits++

			if bit == 1 {
				computed.nset++
				sum += sizes[index]
				if sum >= total {
					break
				}
			}
		}

		if sum == total {
			if current.equals(computed) == false {
				current.assign(computed)
				fmt.Printf("total:%v, formula:", sum)
				sep := ""
				for i, b := range current.bits {
					if i == current.nbits {
						break
					}
					if b == 1 {
						fmt.Printf("%s%v", sep, sizes[i])
						sep = "+"
					}
				}
				fmt.Print(", ")
				number.print()
				combinations++

				if computed.nset < minimum {
					minimum = computed.nset
					minimumcount = 1
				} else if computed.nset == minimum {
					minimumcount++
				}
			}
		}

		if number.increment() == false {
			break
		}
	}
	return
}

func main() {
	total, sizes := readInputFromFile("input.text")
	fmt.Println(total)
	fmt.Println(sizes)
	combinations, minimum := compute(total, sizes)
	fmt.Printf("Total number of combinations %v, minimum %v\n", combinations, minimum)
}
