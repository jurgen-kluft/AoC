package main

/*
--- Day 12: JSAbacusFramework.io ---

Santa's Accounting-Elves need help balancing the books after a recent order.
Unfortunately, their accounting software uses a peculiar storage format.
That's where you come in.

They have a JSON document which contains a variety of things: arrays ([1,2,3]),
objects ({"a":1, "b":2}), numbers, and strings. Your first job is to simply
find all of the numbers throughout the document and add them together.

For example:

[1,2,3] and {"a":2,"b":4} both have a sum of 6.
[[[3]]] and {"a":{"b":4},"c":-1} both have a sum of 3.
{"a":[-1,1]} and [-1,{"a":1}] both have a sum of 0.
[] and {} both have a sum of 0.
You will not encounter any strings containing numbers.

What is the sum of all numbers in the document?

*/

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
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

func readJSON(filename string) (jsoncontent string) {
	computator := func(text string, line int) {
		jsoncontent += text
	}
	iterateOverLinesInTextFile(filename, computator)
	return
}

func jsonProcessMap(amap map[string]interface{}, action func(interface{})) {
	for _, v := range amap {
		switch reflect.TypeOf(v).Kind() {
		case reflect.String:
			if v == "red" {
				return
			}
		}
	}

	for _, v := range amap {
		switch reflect.TypeOf(v).Kind() {
		case reflect.Map:
			jsonProcessMap(v.(map[string]interface{}), action)
		case reflect.Slice:
			jsonProcessList(v.([]interface{}), action)
		default:
			action(v)
		}
	}
}

func jsonProcessList(alist []interface{}, action func(interface{})) {
	for _, v := range alist {
		switch reflect.TypeOf(v).Kind() {
		case reflect.Map:
			jsonProcessMap(v.(map[string]interface{}), action)
		case reflect.Slice:
			jsonProcessList(v.([]interface{}), action)
		default:
			//fmt.Printf("Pos:%v, Value:%v\n", p, v)
			action(v)
		}
	}
}

func main() {
	jsontext := readJSON("input.text")

	jsontree := make(map[string]interface{})
	json.Unmarshal([]byte(jsontext), &jsontree)

	allNumbersAdded := 0.0
	valueHandler := func(value interface{}) {
		switch reflect.TypeOf(value).Kind() {
		case reflect.Float64:
			f := value.(float64)
			allNumbersAdded += f
		case reflect.String:
		}
		//fmt.Printf("value: %v, type:%v\n", value, reflect.TypeOf(value).Kind())
	}

	for _, v := range jsontree {
		//fmt.Printf("index:%v kind:%s  type:%s\n", n, reflect.TypeOf(v).Kind(), reflect.TypeOf(v))
		alist := v.([]interface{})
		jsonProcessList(alist, valueHandler)
	}

	fmt.Printf("All numbers added up = %v", allNumbersAdded)
}
