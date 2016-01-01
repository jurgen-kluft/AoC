package main

/*
--- Day 9: All in a Single Night ---

Every year, Santa manages to deliver all of his presents in a single night.

This year, however, he has some new locations to visit; his elves have provided him the distances between every pair of locations. He can start and end at any two (different) locations he wants, but he must visit each location exactly once. What is the shortest distance he can travel to achieve this?

For example, given the following distances:

London to Dublin = 464
London to Belfast = 518
Dublin to Belfast = 141
The possible routes are therefore:

Dublin -> London -> Belfast = 982
London -> Dublin -> Belfast = 605
London -> Belfast -> Dublin = 659
Dublin -> Belfast -> London = 659
Belfast -> Dublin -> London = 605
Belfast -> London -> Dublin = 982
The shortest of these is London -> Dublin -> Belfast = 605, and so the answer is 605 in this example.

What is the distance of the shortest route?
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

type city struct {
	name  string
	roads []string
}

func decodeRoute(str string) (from string, to string, distance int) {
	fmt.Sscanf(str, "%s to %s = %d", &from, &to, &distance)
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
	for i := 1; i < t.length; i++ {
		from := t.route[i-1]
		to := t.route[i]
		d := from.distance[to.name]
		fmt.Printf(" -> %d,%s", d, to.name)
		distance += d
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

	nto.roads = append(nto.roads, nfrom)
	nto.distance[from] = distance
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
	}

	distanceOfShortestRoute, distanceOfLongestRoute = computeShortestRoute(graph)
	return
}

func main() {
	distanceOfShortestRoute, distanceOfLongestRoute := findShortestRouteFromFile("input.text")
	fmt.Printf("The shortest route is: %v\n", distanceOfShortestRoute)
	fmt.Printf("The longest route is: %v\n", distanceOfLongestRoute)
}
