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

func minInt(a int, b int) (min int) {
	min = a
	if b < min {
		min = b
	}
	return
}

func minimumSidesOfBox(L int, W int, H int) (min1 int, min2 int) {
	min1 = minInt(L, W)
	if H < min1 {
		min1 = H
		min2 = minInt(L, W)
	} else if min1 == L {
		min2 = minInt(W, H)
	} else if min1 == W {
		min2 = minInt(L, H)
	}
	return
}

func howMuchPaperToOrder(filename string) (totalSurfaceAreaOfWrappingPaper int, totalFeetOfRibbon int, ok bool) {
	totalSurfaceAreaOfWrappingPaper = 0
	totalFeetOfRibbon = 0

	computator := func(line string) {
		var L, W, H = lineToLWH(line)

		// Compute length of ribbon = volume of gift + sum of 2 smallest sides
		volumeOfPresentInFeet := L * W * H
		minSide1, minSide2 := minimumSidesOfBox(L, W, H)
		totalFeetOfRibbon += 2*minSide1 + 2*minSide2
		totalFeetOfRibbon += volumeOfPresentInFeet

		// Surface-Area = 2LW + 2LH + 2WH
		totalSurfaceAreaOfWrappingPaper += (2 * L * W) + (2 * L * H) + (2 * W * H)
		// Slack = L*W
		totalSurfaceAreaOfWrappingPaper += computeSlackArea(L, W, H)
	}

	iterateOverLinesInTextFile(filename, computator)
	ok = true
	return
}

func main() {
	var totalSurfaceAreaOfWrappingPaper, totalFeetOfRibbon, ok = howMuchPaperToOrder("input.text")
	if ok {
		fmt.Printf("The elves need to order %v square-feet of wrapping paper\n", totalSurfaceAreaOfWrappingPaper)
		fmt.Printf("and they need to order %v feet of ribbon\n", totalFeetOfRibbon)
	} else {
		fmt.Printf("Could not process the input\n")
	}
}
