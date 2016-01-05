package main

/*
--- Day 11: Corporate Policy ---

Santa's previous password expired, and he needs help choosing a new one.

To help him remember his new password after the old one expires, Santa has
devised a method of coming up with a password based on the previous one.
Corporate policy dictates that passwords must be exactly eight lowercase letters
 (for security reasons), so he finds his new password by incrementing his old
 password string repeatedly until it is valid.

Incrementing is just like counting with numbers: xx, xy, xz, ya, yb, and so on.
Increase the rightmost letter one step; if it was z, it wraps around to a, and
repeat with the next letter to the left until one doesn't wrap around.

Unfortunately for Santa, a new Security-Elf recently started, and he has imposed
some additional password requirements:

Passwords must include one increasing straight of at least three letters, like
abc, bcd, cde, and so on, up to xyz. They cannot skip letters; abd doesn't count.
Passwords may not contain the letters i, o, or l, as these letters can be
mistaken for other characters and are therefore confusing.
Passwords must contain at least two different, non-overlapping pairs of letters,
like aa, bb, or zz.

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

func passwordRuneIsValid(r rune) (isValid bool) {
	return r != 'i' && r != 'o' && r != 'l'
}

func passwordPairIsValid(ph rune, pl rune) bool {
	return passwordRuneIsValid(ph) && passwordRuneIsValid(pl)
}

func passwordHasLetterPairs(password []rune, N int) (has bool) {
	pairs := 0
	prev := rune(0)
	for _, cur := range password {
		if prev == cur {
			pairs++
			prev = rune(0)
		} else {
			prev = cur
		}
	}

	has = pairs >= N
	return
}

func passwordHasIncreasingLetterRange(password []rune, rangeLength int) bool {
	for i := rangeLength; i < len(password); i++ {
		l := 1
		c := i - rangeLength
		for r := 1; r < rangeLength; r++ {
			if password[c] != (password[c+r] - rune(r)) {
				break
			}
			l++
		}
		if l == rangeLength {
			return true
		}
	}
	return false
}

func passwordHasValidPairs(password []rune) bool {
	n := len(password) / 2
	for p := 0; p < n; p++ {
		h := password[p*2]
		l := password[p*2+1]
		if passwordPairIsValid(h, l) == false {
			return false
		}
	}
	return true
}

func isValidPassword(password []rune) bool {
	if passwordHasValidPairs(password) == false {
		return false
	}
	if passwordHasLetterPairs(password, 2) == false {
		return false
	}
	if passwordHasIncreasingLetterRange(password, 3) == false {
		return false
	}
	return true
}

func incPasswordPair(password []rune, pair int) (wrapped bool) {
	h := password[pair*2]
	l := password[pair*2+1]

	wrapped = false
	for {
		if l == 'z' {
			l = 'a'
			if h == 'z' {
				h = 'a'
				wrapped = true
			} else {
				h++
			}
		} else {
			l++
		}
		if passwordPairIsValid(h, l) == true {
			break
		}
	}
	password[pair*2] = h
	password[pair*2+1] = l
	return wrapped
}

func incPassword(password []rune, pair int) {
	for incPasswordPair(password, pair) {
		if pair > 0 {
			pair--
		} else {
			pair = (len(password) / 2) - 1
		}
	}
}

func fromFile(filename string, linenumber int) (oldpw string, newpw string) {
	computator := func(text string, line int) {
		password := []rune(text)
		pair := (len(password) / 2) - 1
		oldpw = text
		for {
			incPassword(password, pair)
			if isValidPassword(password) {
				break
			}
			fmt.Printf("(pair:%d)%s\n", pair, string(password))
		}
		newpw = string(password)
	}
	iterateOverLinesInTextFile(filename, computator)

	return
}

func main() {
	oldpw, newpw := fromFile("input.text", 2)
	fmt.Printf("Santa's old password was : %s\n", oldpw)
	fmt.Printf("Santa's new password is  : %s\n", newpw)
}
