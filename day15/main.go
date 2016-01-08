package main

/*
--- Day 15: Science for Hungry People ---

Today, you set out on the task of perfecting your milk-dunking cookie recipe.
All you have to do is find the right balance of ingredients.

Your recipe leaves room for exactly 100 teaspoons of ingredients. You make a
list of the remaining ingredients you could use to finish the recipe (your
puzzle input) and their properties per teaspoon:

capacity (how well it helps the cookie absorb milk)
durability (how well it keeps the cookie intact when full of milk)
flavor (how tasty it makes the cookie)
texture (how it improves the feel of the cookie)
calories (how many calories it adds to the cookie)
You can only measure ingredients in whole-teaspoon amounts accurately, and
you have to be accurate so you can reproduce your results in the future.
The total score of a cookie can be found by adding up each of the properties
(negative totals become 0) and then multiplying together everything except
calories.

For instance, suppose you have these two ingredients:

Butterscotch: capacity -1, durability -2, flavor 6, texture 3, calories 8
Cinnamon: capacity 2, durability 3, flavor -2, texture -1, calories 3
Then, choosing to use 44 teaspoons of butterscotch and 56 teaspoons of cinnamon
 (because the amounts of each ingredient must add up to 100) would result in a
 cookie with the following properties:

A capacity of 44*-1 + 56*2 = 68
A durability of 44*-2 + 56*3 = 80
A flavor of 44*6 + 56*-2 = 152
A texture of 44*3 + 56*-1 = 76
Multiplying these together (68 * 80 * 152 * 76, ignoring calories for now)
results in a total score of 62842880, which happens to be the best score
possible given these ingredients. If any properties had produced a negative
total, it would have instead become zero, causing the whole score to multiply
to zero.

Given the ingredients in your kitchen and their properties, what is the total
score of the highest-scoring cookie you can make?
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

type ingredient struct {
	name      string
	propValue []int
}

var sPropertyNames = [...]string{"capacity", "durability", "flavor", "texture", "calories"}

func propCount() int {
	return len(sPropertyNames)
}

func propName(i int) string {
	return sPropertyNames[i]
}

var sPropertyUsedInScore = [...]bool{true, true, true, true, false}

func propUsedInScore(i int) bool {
	return sPropertyUsedInScore[i]
}

// Part 2 of this problem requires a certain property (Calories) to have a
// certain value in a final score.
//var sPropertyCompare = [...]bool{false, false, false, false, false}
var sPropertyCompare = [...]bool{false, false, false, false, true}
var sPropertyCompareValue = [...]int{0, 0, 0, 0, 500}

func propCompare(i int, v int) bool {
	if sPropertyCompare[i] {
		return sPropertyCompareValue[i] == v
	}
	return true
}

func (r *ingredient) print() {
	fmt.Printf("%s {", r.name)
	for i := range r.propValue {
		fmt.Printf(" %s:%v", propName(i), r.propValue[i])
	}
	fmt.Print(" }\n")
}

func (r *ingredient) deserialize(str string) {
	// Example line:
	//      Sprinkles: capacity 2, durability 0, flavor -2, texture 0, calories 3
	r.propValue = make([]int, 5, 5)
	fmt.Sscanf(str, "%s capacity %d, durability %d, flavor %d, texture %d, calories %d", &r.name, &r.propValue[0], &r.propValue[1], &r.propValue[2], &r.propValue[3], &r.propValue[4])
	return
}

func readInputFromFile(filename string) (ingredients []*ingredient) {
	ingredients = make([]*ingredient, 0)

	computator := func(text string, line int) {
		r := new(ingredient)
		r.deserialize(text)
		ingredients = append(ingredients, r)
	}
	iterateOverLinesInTextFile(filename, computator)

	return
}

type cookie struct {
	spoons int
	column int
	distr  []int
}

func (c *cookie) init(spoons int, N int) {
	c.spoons = spoons

	// Initialize an array to keep track of ingredient distribution over N spoons
	c.distr = make([]int, N)
	for i := 0; i < N; i++ {
		c.distr[i] = 0
	}
	c.column = N - 1
	c.distr[c.column] = c.spoons
}

func (c *cookie) printSpoon(eol string) {
	for i := 0; i < len(c.distr); i++ {
		fmt.Printf("[%d]", c.distr[i])
	}
	fmt.Printf(eol)
}

func (c *cookie) sum(column int) (sum int) {
	sum = 0
	for i := column; i >= 0; i-- {
		sum += c.distr[i]
	}
	return
}

func (c *cookie) increment(column int, umax int) (ok bool) {
	for column >= 0 {
		// The maximum of the current column is the current maximum
		// minus the sum of the upper columns
		cmax := umax - c.sum(column-1)
		if c.distr[column] == cmax {
			c.distr[column] = 0
			column--
		} else {
			c.distr[column]++
			return true
		}
	}
	return false
}

func (c *cookie) nextSpoonCombination() bool {
	// End condition
	if c.distr[0] == c.spoons {
		return false
	}

	max := c.sum(c.column-1) + 1
	if max > c.spoons {
		max = c.spoons
	}

	if c.increment(c.column-1, max) {
		c.distr[c.column] = c.spoons - c.sum(c.column-1)
	} else {
		c.column--
	}
	return true
}

func (c *cookie) computeBestCookie(ingredients []*ingredient) (bestScore int) {
	cookieScore := new(scoring)
	cookieScore.init(propCount())

	bestScore = 0

	c.init(c.spoons, len(ingredients))
	for r := 0; ; r++ {
		cookieScore.reset()
		for i, s := range c.distr {
			cookieScore.addIngredient(s, ingredients[i])
		}

		if cookieScore.checkPropScores() {
			score := cookieScore.finalizeScore()
			if score > bestScore {
				bestScore = score
				c.printSpoon(" | ")
				cookieScore.printScore(" | ")
				fmt.Printf("-----> Best Score: %v", bestScore)
				fmt.Print("\n")
			} else if r%100 == 0 {
				c.printSpoon("\n")
			}
		}

		if c.nextSpoonCombination() == false {
			break
		}
	}
	return
}

func (c *cookie) iterateAllSpoons(N int) {
	c.init(c.spoons, N)
	for r := 0; ; r++ {
		c.printSpoon("\n")
		if c.nextSpoonCombination() == false {
			break
		}
	}
}

type scoring struct {
	props []int
}

func (s *scoring) init(size int) {
	s.props = make([]int, size, size)
}

func (s *scoring) reset() {
	for i := 0; i < len(s.props); i++ {
		s.props[i] = 0
	}
}

func (s *scoring) addIngredient(spoons int, g *ingredient) {
	if spoons > 0 {
		for p, v := range g.propValue {
			s.props[p] = s.props[p] + spoons*v
		}
	}
}

func (s *scoring) printScore(eol string) {
	score := s.finalizeScore()
	fmt.Printf("Score[%v] : ", score)
	for i := 0; i < len(s.props); i++ {
		fmt.Printf("[%d]", s.props[i])
	}
	fmt.Printf(eol)
}

func (s *scoring) checkPropScores() bool {
	for p := range sPropertyCompare {
		if propCompare(p, s.props[p]) == false {
			return false
		}
	}
	return true
}

func (s *scoring) finalizeScore() (score int) {
	score = 1
	negative := false
	for i, v := range s.props {
		if propUsedInScore(i) {
			negative = negative || v < 0
			score = score * v
		}
	}
	if negative {
		score = 0
	}
	return
}

func printAll(ingredients []*ingredient) {
	for _, r := range ingredients {
		r.print()
	}
}

func main() {
	ingredients := readInputFromFile("input.text")
	printAll(ingredients)

	var bestCookieScore int
	bestCookie := &cookie{spoons: 100}
	bestCookieScore = bestCookie.computeBestCookie(ingredients)
	//bestCookie.iterateAllSpoons(3)
	fmt.Printf("Best cookie has score: %v", bestCookieScore)
}
