package main

/*
--- Day 18: Like a GIF For Your Yard ---

After the million lights incident, the fire code has gotten stricter: now, at
most ten thousand lights are allowed. You arrange them in a 100x100 grid.

Never one to let you down, Santa again mails you instructions on the ideal
lighting configuration. With so few lights, he says, you'll have to resort to
animation.

Start by setting your lights to the included initial configuration (your puzzle
input). A # means "on", and a . means "off".

Then, animate your grid in steps, where each step decides the next configuration
based on the current one. Each light's next state (either on or off) depends on
its current state and the current states of the eight lights adjacent to it
(including diagonals). Lights on the edge of the grid might have fewer than
eight neighbors; the missing ones always count as "off".

The state a light should have next is based on its current state (on or off)
plus the number of neighbors that are on:

*) A light which is on stays on when 2 or 3 neighbors are on, and turns off otherwise.
*) A light which is off turns on if exactly 3 neighbors are on, and stays off otherwise.
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

func deserializeRow(str string) []byte {
	// Example line:
	//      #..####.##..#...#..#...#...###.#.#.#..#....#.##..#...##...#..#.....##..#####....#.##..##....##.#....

	row := make([]byte, len(str), len(str))
	for i, c := range str {
		if c == '#' {
			row[i] = 1
		} else if c == '@' {
			row[i] = 0x11
		} else if c == '.' {
			row[i] = 0
		} else if c == ',' {
			row[i] = 0x10
		}
	}
	return row
}

func readInputFromFile(filename string) (grid [][]byte) {
	grid = make([][]byte, 0, 100)

	computator := func(text string, line int) {
		row := deserializeRow(text)
		grid = append(grid, row)
	}
	iterateOverLinesInTextFile(filename, computator)

	return
}

func copyGrid(grid [][]byte) [][]byte {
	gridCopy := make([][]byte, len(grid), len(grid))
	for y := range grid {
		row := grid[y]
		rowCopy := make([]byte, len(row), len(row))
		gridCopy[y] = rowCopy
		for x := range row {
			rowCopy[x] = row[x]
		}
	}
	return gridCopy
}
func resetGrid(grid [][]byte) {
	for y := range grid {
		row := grid[y]
		for x := range row {
			row[x] = 0
		}
	}
}

func getState(x int, y int, grid [][]byte) (byte, bool) {
	height := len(grid)
	width := len(grid[0])
	if x < 0 || y < 0 {
		return 0, false
	} else if x >= width || y >= height {
		return 0, false
	}
	return grid[y][x], true
}

func setState(state byte, x int, y int, grid [][]byte) {
	grid[y][x] = state
}

func switchLight(x int, y int, grid [][]byte) byte {
	center, _ := getState(x, y, grid)
	if center&0x10 == 0x10 {
		return center
	}

	// check all 8 neighbors
	neighbors := 0
	neighborsON := 0
	for ny := y - 1; ny <= y+1; ny++ {
		for nx := x - 1; nx <= x+1; nx++ {
			if nx == x && ny == y {
				continue
			}
			neigbor, ok := getState(nx, ny, grid)
			if ok {
				neighbors++
				if neigbor&1 == 1 {
					neighborsON++
				}
			}
		}
	}

	// Flip the state of the light
	// state = ON, stays ON when 2 OR 3 neighbors are ON, otherwise --> OFF
	// state = OFF, 3 neighbors are ON --> ON
	if center&1 == 1 {
		//fmt.Printf("ON: neighbors ON = %v\n", neighborsON)
		if 3 == neighborsON || 2 == neighborsON {
			center = 1
		} else {
			center = 0
		}
	} else {
		//fmt.Printf("OFF: neighbors ON = %v\n", neighborsON)
		if 3 == neighborsON {
			center = 1
		}
	}
	return center
}

func printGrid(grid [][]byte) {
	for y := 0; y < len(grid); y++ {
		row := grid[y]
		for x := 0; x < len(row); x++ {
			state, _ := getState(x, y, grid)
			if state&1 == 1 {
				fmt.Print("#")
			} else if state&1 == 0 {
				fmt.Print(".")
			}
		}
		fmt.Print("\n")
	}
}

func animateFrame(current [][]byte, next [][]byte) int {
	ON := 0
	for y := 0; y < len(current); y++ {
		row := current[y]
		for x := 0; x < len(row); x++ {
			s := switchLight(x, y, current)
			setState(s, x, y, next)
			if s&1 == 1 {
				ON++
			}
		}
	}
	return ON
}

func animate(animations int, grid [][]byte) int {
	ON := 0
	frame := copyGrid(grid)
	resetGrid(frame)

	for a := 0; a < animations; a++ {
		fmt.Printf("Frame: %v\n", a)
		if (a & 1) == 0 {
			ON = animateFrame(grid, frame)
			printGrid(frame)
		} else {
			ON = animateFrame(frame, grid)
			printGrid(grid)
		}
	}
	return ON
}

func setStuckLights(grid [][]byte) {
	h := len(grid) - 1
	w := len(grid[0]) - 1
	grid[0][0] = 0x11
	grid[0][w] = 0x11
	grid[h][0] = 0x11
	grid[h][w] = 0x11
}

func main() {
	grid := readInputFromFile("input.text")
	setStuckLights(grid)
	numON := animate(100, grid)
	fmt.Printf("Total number of lights ON %v\n", numON)
}
