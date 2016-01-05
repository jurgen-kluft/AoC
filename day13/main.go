package main

/*
--- Day 13: Knights of the Dinner Table ---

In years past, the holiday feast with your family hasn't gone so well.
Not everyone gets along! This year, you resolve, will be different.
You're going to find the optimal seating arrangement and avoid all those
awkward conversations.

You start by writing up a list of everyone invited and the amount their
happiness would increase or decrease if they were to find themselves sitting
next to each other person. You have a circular table that will be just big
enough to fit everyone comfortably, and so each person will have exactly
two neighbors.

For example, suppose you have only four attendees planned, and you calculate
their potential happiness as follows:

Alice would gain 54 happiness units by sitting next to Bob.
Alice would lose 79 happiness units by sitting next to Carol.
Alice would lose 2 happiness units by sitting next to David.
Bob would gain 83 happiness units by sitting next to Alice.
Bob would lose 7 happiness units by sitting next to Carol.
Bob would lose 63 happiness units by sitting next to David.
Carol would lose 62 happiness units by sitting next to Alice.
Carol would gain 60 happiness units by sitting next to Bob.
Carol would gain 55 happiness units by sitting next to David.
David would gain 46 happiness units by sitting next to Alice.
David would lose 7 happiness units by sitting next to Bob.
David would gain 41 happiness units by sitting next to Carol.
Then, if you seat Alice next to David, Alice would lose 2 happiness units
(because David talks so much), but David would gain 46 happiness units
(because Alice is such a good listener), for a total change of 44.

If you continue around the table, you could then seat Bob next to Alice
(Bob gains 83, Alice gains 54). Finally, seat Carol, who sits next to Bob
(Carol gains 60, Bob loses 7) and David (Carol gains 55, David gains 41).
The arrangement looks like this:

     +41 +46
+55   David    -2
Carol       Alice
+60    Bob    +54
     -7  +83
After trying every other seating arrangement in this hypothetical scenario,
you find that this one is the most optimal, with a total change in happiness
of 330.

What is the total change in happiness for the optimal seating arrangement
of the actual guest list?

*/

/*
Note: This is known as the 'Travelling Salesman Problem' (TSP)
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

type city struct {
	name  string
	roads []string
}

func decodeRoute(str string) (from string, to string, distance int) {
	// Alice would lose 2 happiness units by sitting next to Bob.
	impact := ""
	fmt.Sscanf(str, "%s would %s %d happiness units by sitting next to %s", &from, &impact, &distance, &to)
	to = strings.TrimRight(to, ".")
	if impact == "lose" {
		distance = -distance
	}
	return
}

type node struct {
	name     string
	roads    []*node
	distance map[string]int
}

type traversal struct {
	visited map[string]int
	route   []*node
	dir     int
	length  int
	maximum int
}

func (t *traversal) back() bool {
	if t.length == 0 {
		return false
	}

	// travel backwards, remove the current node from the route
	// reduce the route length and remove the node from 'visited'
	t.length--
	current := t.route[t.length]
	delete(t.visited, current.name)

	if t.length == 0 {
		return false
	}

	// re-initialize the direction of the now current node
	current = t.route[t.length-1]
	t.dir = t.visited[current.name]

	return t.dir < len(current.roads)
}

func (t *traversal) forward() bool {
	// No need to travel when the route is complete or empty
	if t.complete() || t.empty() {
		return false
	}
	current := t.route[t.length-1]
	for t.dir < len(current.roads) {
		next := current.roads[t.dir]
		if _, visited := t.visited[next.name]; !visited {
			// Remember the next direction we should travel on current node
			t.visited[current.name] = t.dir + 1
			// Initialize the direction of the new node to the first direction
			t.dir = 0
			t.visited[next.name] = 0
			t.route[t.length] = next
			t.length++
			return true
		}
		// Try next dir to see if we can travel into that direction
		t.dir++
	}

	return false
}

func (t *traversal) start(origin *node, maximum int) {
	t.visited = make(map[string]int, maximum)
	t.route = make([]*node, maximum, maximum)
	t.length = 1
	t.maximum = maximum
	t.route[0] = origin
	t.dir = 0
}

func (t *traversal) empty() bool {
	return t.length == 0
}

func (t *traversal) complete() bool {
	return t.length == t.maximum
}

func (t *traversal) distance() int {
	distance := 0
	fmt.Printf("%s", t.route[0].name)
	for i := 0; i < t.length; i++ {
		from := t.route[i]
		to := t.route[(i+1)%t.length]
		d1 := from.distance[to.name]
		d2 := to.distance[from.name]
		fmt.Printf(" -> %d,%d,%s", d1, d2, to.name)
		distance += (d1 + d2)
	}
	fmt.Printf(" = %v\n", distance)
	return distance
}

func computeShortestRoute(graph map[string]*node) (distanceOfShortestRoute int, distanceOfLongestRoute int) {

	// Set shortest distance to a very high number
	distanceOfShortestRoute = 1 * 1024 * 1024 * 1024
	distanceOfLongestRoute = 0

	totalNumberOfRoutes := 0

	traverse := new(traversal)
	for _, originNode := range graph {
		traverse.start(originNode, len(graph))
		// Try all possible routes from this origin
		for traverse.empty() == false {

			// Travel forward until we have a complete route
			for traverse.complete() == false {
				if traverse.forward() == false {
					break
				}
			}

			// If traveling forward got us a complete route
			// then compute the distance and see if this is
			// now the shortest

			if traverse.complete() {
				distance := traverse.distance()
				totalNumberOfRoutes++
				if distance < distanceOfShortestRoute {
					distanceOfShortestRoute = distance
				}
				if distance > distanceOfLongestRoute {
					distanceOfLongestRoute = distance
				}
			}

			// Travel backwards on the route until we find a node that still
			// has directions that we can travel on
			for traverse.back() == false {
				if traverse.empty() {
					break
				}
			}
		}
		break
	}
	fmt.Printf("Traveled %v number of routes\n", totalNumberOfRoutes)
	return
}

func addRoute(from string, to string, distance int, graph map[string]*node) {
	if nfrom, exists := graph[from]; !exists {
		nfrom = &node{name: from, roads: make([]*node, 0, 8), distance: make(map[string]int)}
		graph[from] = nfrom
	}
	if nto, exists := graph[to]; !exists {
		nto = &node{name: to, roads: make([]*node, 0, 8), distance: make(map[string]int)}
		graph[to] = nto
	}

	nfrom := graph[from]
	nto := graph[to]

	nfrom.roads = append(nfrom.roads, nto)
	nfrom.distance[to] = distance
}

func findShortestRouteFromFile(filename string) (distanceOfShortestRoute int, distanceOfLongestRoute int) {
	graph := make(map[string]*node, 64)

	computator := func(text string, line int) {
		from, to, distance := decodeRoute(text)
		//fmt.Printf("From %s to %s, distance %v\n", from, to, distance)
		addRoute(from, to, distance, graph)
	}
	iterateOverLinesInTextFile(filename, computator)

	printgraph := false
	if printgraph {
		for nodeName, nodeObject := range graph {
			fmt.Printf("origin: %s, roads", nodeName)
			for _, nodeTo := range nodeObject.roads {
				distance := nodeObject.distance[nodeTo.name]
				fmt.Printf(": %s(%d)", nodeTo.name, distance)
			}
			fmt.Printf("\n")
		}
	} else {
		distanceOfShortestRoute, distanceOfLongestRoute = computeShortestRoute(graph)
	}
	return
}

func main() {
	distanceOfShortestRoute, distanceOfLongestRoute := findShortestRouteFromFile("input2.text")
	fmt.Printf("The shortest route is: %v\n", distanceOfShortestRoute)
	fmt.Printf("The longest route is: %v\n", distanceOfLongestRoute)
}
