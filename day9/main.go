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

func back(self traversal) bool {
	if self.length == 0 {
		return false
	}

	// travel backwards, remove the current node from 'visited'
	// reduce the travelled length
	current := self.route[self.length]
	delete(self.visited, current.name)
	self.length--

	// re-initialize the direction of the now current node
	current = self.route[self.length]
	self.dir = self.visited[current.name]

	return self.dir < len(current.roads)
}

func empty(self traversal) bool {
	return self.length == 0
}

func complete(self traversal) bool {
	return self.length == self.maximum
}

func distance(self traversal) int {
	distance := 0
	if self.length > 0 {
		current := self.route[0]
		for p, n := range self.route {
			if p > 0 {
				distance += current.distance[n.name]
				current = n
			}
		}
	}
	return distance
}

func forward(self traversal) bool {
	current := self.route[self.length]
	for self.dir < len(current.roads) {
		next := current.roads[self.dir]
		if _, visited := self.visited[next.name]; !visited {
			self.visited[current.name] = self.dir + 1
			self.visited[next.name] = 0
			self.dir = 0
			self.length++
			self.route[self.length] = next
			return true
		}
	}

	return false
}

func addRoute(from string, to string, distance int, graph map[string]*node) {
	var nfrom = new(node)
	var nto = new(node)
	var exists = false
	if nfrom, exists = graph[from]; !exists {
		nfrom = &node{name: from, roads: make([]*node, 0), distance: make(map[string]int)}
		graph[from] = nfrom
	}
	if nto, exists = graph[to]; !exists {
		nto = &node{name: to, roads: make([]*node, 0), distance: make(map[string]int)}
		graph[to] = nto
	}
	nfrom.roads = append(nfrom.roads, nto)
	nfrom.distance[to] = distance
}

func computeShortestRoute(graph map[string]*node) (distanceOfShortestRoute int) {
	for originName, originNode := range graph {
		traverse := new(traversal)
		for traverse.complete == false && traverse.
	}
	return
}

func findShortestRouteFromFile(filename string) (distanceOfShortestRoute int) {
	graph := make(map[string]*node)

	computator := func(text string, line int) {
		from, to, distance := decodeRoute(text)
		fmt.Printf("From %s to %s, distance %v\n", from, to, distance)
		addRoute(from, to, distance, graph)
	}
	iterateOverLinesInTextFile(filename, computator)

	distanceOfShortestRoute = computeShortestRoute(graph)
	return
}

func main() {
	distanceOfShortestRoute := findShortestRouteFromFile("input.text")
	fmt.Printf("The shortest route is: %v", distanceOfShortestRoute)
}
