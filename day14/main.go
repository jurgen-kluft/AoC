package main

/*
--- Day 14: Reindeer Olympics ---

This year is the Reindeer Olympics! Reindeer can fly at high speeds, but must
rest occasionally to recover their energy. Santa would like to know which of
his is fastest, and so he has them race.

Reindeer can only either be flying (always at their top speed) or resting
(not moving at all), and always spend whole seconds in either state.

For example, suppose you have the following Reindeer:

Comet can fly 14 km/s for 10 seconds, but then must rest for 127 seconds.
Dancer can fly 16 km/s for 11 seconds, but then must rest for 162 seconds.
After one second, Comet has gone 14 km, while Dancer has gone 16 km.
After ten seconds, Comet has gone 140 km, while Dancer has gone 160 km.
On the eleventh second, Comet begins resting (staying at 140 km), and Dancer
continues on for a total distance of 176 km. On the 12th second, both reindeer
are resting. They continue to rest until the 138th second, when Comet flies
for another ten seconds. On the 174th second, Dancer flies for another 11
seconds.

In this example, after the 1000th second, both reindeer are resting, and Comet
is in the lead at 1120 km (poor Dancer has only gotten 1056 km by that point).
So, in this situation, Comet would win (if the race ended at 1000 seconds).

Given the descriptions of each reindeer (in your puzzle input), after exactly
2503 seconds, what distance has the winning reindeer traveled?

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

type reindeer struct {
	name         string
	flySpeed     int
	flyDuration  int
	restDuration int
}

func (r *reindeer) print() {
	fmt.Printf("%s can fly for %v s with a speed of %v km/s but then has to rest for %v s.\n", r.name, r.flySpeed, r.flyDuration, r.restDuration)
}

func deserializeReindeer(str string) (r *reindeer) {
	r = new(reindeer)

	// Example line:
	//      Dancer can fly 37 km/s for 1 seconds, but then must rest for 36 seconds.
	var name string
	var flySpeed int
	var flyDuration int
	var restDuration int

	fmt.Sscanf(str, "%s can fly %v km/s for %v seconds, but then must rest for %v seconds.", &name, &flySpeed, &flyDuration, &restDuration)
	r.name = name
	r.flySpeed = flySpeed
	r.flyDuration = flyDuration
	r.restDuration = restDuration

	return
}

func readInputFromFile(filename string) (reindeers []*reindeer) {

	reindeers = make([]*reindeer, 0)

	computator := func(text string, line int) {
		r := deserializeReindeer(text)
		reindeers = append(reindeers, r)
	}
	iterateOverLinesInTextFile(filename, computator)

	return
}

func printAllReindeerObjects(reindeers []*reindeer) {
	for _, r := range reindeers {
		r.print()
	}
}

func (r *reindeer) computeDistanceOfReindeer(time int) (distance int) {
	numFullCycles := time / (r.flyDuration + r.restDuration)
	d := numFullCycles * (r.flySpeed * r.flyDuration)

	lastPartialCycle := time - (numFullCycles * (r.flyDuration + r.restDuration))
	if lastPartialCycle > r.flyDuration {
		lastPartialCycle = r.flyDuration
	}

	distance = d + lastPartialCycle*r.flySpeed
	return
}

func computeWinningDistance(reindeers []*reindeer, time int) (winningDistance int) {
	for _, r := range reindeers {
		d := r.computeDistanceOfReindeer(time)
		if d > winningDistance {
			winningDistance = d
		}
	}
	return
}

func computeWinningPoints(reindeers []*reindeer, time int) (winningPoints int) {
	points := make([]int, len(reindeers))
	dists := make([]int, len(reindeers))
	for i := 0; i < len(reindeers); i++ {
		points[i] = 0
		dists[i] = 0
	}

	for t := 1; t < time; t++ {
		headdist := 0
		for i, r := range reindeers {
			d := r.computeDistanceOfReindeer(t)
			dists[i] = d
			if d > headdist {
				headdist = d
			}
		}
		/// Handout points to the reindeer(s) that are the head of the race
		for i, d := range dists {
			if d >= headdist {
				points[i]++
			}
		}
	}

	// Figure out the winning points
	for _, p := range points {
		if p > winningPoints {
			winningPoints = p
		}
	}

	return
}

func main() {
	reindeers := readInputFromFile("input.text")
	printAllReindeerObjects(reindeers)
	winningDistance := computeWinningDistance(reindeers, 2503)
	fmt.Printf("The winning distance = %v\n", winningDistance)
	winningPoints := computeWinningPoints(reindeers, 2503)
	fmt.Printf("The winning points   = %v\n", winningPoints)
}
