package main

/*
--- Day 6: Probably a Fire Hazard ---

Because your neighbors keep defeating you in the holiday house decorating contest year after year, you've decided to deploy one million lights in a 1000x1000 grid.

Furthermore, because you've been especially nice this year, Santa has mailed you instructions on how to display the ideal lighting configuration.

Lights in your grid are numbered from 0 to 999 in each direction; the lights at each corner are at 0,0, 0,999, 999,999, and 999,0. The instructions include whether to turn on, turn off, or toggle various inclusive ranges given as coordinate pairs. Each coordinate pair represents opposite corners of a rectangle, inclusive; a coordinate pair like 0,0 through 2,2 therefore refers to 9 lights in a 3x3 square. The lights all start turned off.

To defeat your neighbors this year, all you have to do is set up your lights by doing the instructions Santa sent you in order.

For example:

turn on 0,0 through 999,999 would turn on (or leave on) every light.
toggle 0,0 through 999,0 would toggle the first line of 1000 lights, turning off the ones that were on, and turning on the ones that were off.
turn off 499,499 through 500,500 would turn off (or leave off) the middle four lights.
After following the instructions, how many lights are lit?
*/

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type eLightAction int

const (
	cTurnOn  eLightAction = 1
	cTurnOff eLightAction = 2
	cToggle  eLightAction = 3
)

type loc2D struct {
	X int
	Y int
}

type range2D struct {
	from loc2D
	to   loc2D
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

func decodeLine(line string) (lightAction eLightAction, lightRange range2D) {
	lightAction = cTurnOn

	var av, as = "", ""
	var fx, fy, tx, ty = 0, 0, 0, 0

	if strings.HasPrefix(line, "turn") {
		fmt.Sscanf(line, "%s %s %d,%d through %d,%d", &av, &as, &fx, &fy, &tx, &ty)
	} else {
		fmt.Sscanf(line, "%s %d,%d through %d,%d", &as, &fx, &fy, &tx, &ty)
	}

	lightRange.from.X = fx
	lightRange.from.Y = fy
	lightRange.to.X = tx
	lightRange.to.Y = ty

	// fmt.Printf("Action = '%s', range(%d,%d - %d,%d)\n", as, fx, fy, tx, ty)

	if as == "on" {
		lightAction = cTurnOn
	} else if as == "off" {
		lightAction = cTurnOff
	} else if as == "toggle" {
		lightAction = cToggle
	}

	return
}

func determineLitCountByAction(lightAction eLightAction, current bool) (lit bool) {
	switch lightAction {
	case cTurnOn:
		lit = true
	case cTurnOff:
		lit = false
	case cToggle:
		lit = !current
	}
	return
}

func applyLightActionToLightRange(lightAction eLightAction, lightRange range2D, lightGrid map[uint64]bool) {
	for x := lightRange.from.X; x <= lightRange.to.X; x++ {
		for y := lightRange.from.Y; y <= lightRange.to.Y; y++ {
			loc := (uint64(x) << 32) | (uint64(y) & 0xFFFFFFFF)
			if lit, ok := lightGrid[loc]; ok == false {
				lit = false
				lightGrid[loc] = determineLitCountByAction(lightAction, lit)
			} else {
				lightGrid[loc] = determineLitCountByAction(lightAction, lit)
			}
		}
	}
}

func countNumberOfLitLights(lightGrid map[uint64]bool) (lit int) {
	lit = 0
	for _, v := range lightGrid {
		if v {
			lit++
		}
	}
	return
}

func howManyLightsAreLitFromFile(filename string) (lit int) {
	lit = 0

	lightGrid := make(map[uint64]bool)

	computator := func(line string) {
		lightAction, lightRange := decodeLine(line)
		applyLightActionToLightRange(lightAction, lightRange, lightGrid)
	}

	iterateOverLinesInTextFile(filename, computator)

	lit = countNumberOfLitLights(lightGrid)

	return
}

func main() {
	var howManyLightsAreLit = howManyLightsAreLitFromFile("input.text")
	fmt.Printf("There are %v lights lit", howManyLightsAreLit)
}
