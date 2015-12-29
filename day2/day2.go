package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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

func lineToLWH(line string) (L int, W int, H int) {
	result := strings.Split(line, "x")
	L, _ = strconv.Atoi(result[0])
	W, _ = strconv.Atoi(result[1])
	H, _ = strconv.Atoi(result[2])
	return
}

func computeSlackArea(L int, W int, H int) (slack int) {
	slack = L * W
	if slack > (L * H) {
		slack = L * H
	}
	if slack > (W * H) {
		slack = W * H
	}
	return
}

func howMuchPaperToOrder(filename string) (totalSurfaceArea int, ok bool) {
	totalSurfaceArea = 0
	computator := func(line string) {
		var L, W, H = lineToLWH(line)
		// Surface-Area = 2LW + 2LH + 2WH
		totalSurfaceArea += (2 * L * W) + (2 * L * H) + (2 * W * H)
		// Slack = L*W
		totalSurfaceArea += computeSlackArea(L, W, H)
	}

	iterateOverLinesInTextFile(filename, computator)
	ok = true
	return
}

func main() {
	var totalSurfaceArea, ok = howMuchPaperToOrder("input.text")
	if ok {
		fmt.Printf("The elves need to order %v square-feet of wrapping paper", totalSurfaceArea)
	} else {
		fmt.Printf("Could not process the input")
	}
}
