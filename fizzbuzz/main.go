package main

import "fmt"

func printFizzBuzz(number int) (printed bool) {
	isMultipleOf3 := (number % 3) == 0
	isMultipleOf5 := (number % 5) == 0

	if isMultipleOf3 {
		fmt.Printf("Fizz")
	}
	if isMultipleOf5 {
		fmt.Printf("Buzz")
	}

	if isMultipleOf3 || isMultipleOf5 {
		fmt.Print("\n")
	}

	return isMultipleOf3 || isMultipleOf5
}

func main() {
	for i := 0; i < 100; i++ {
		if printFizzBuzz(i) == false {
			fmt.Printf("%v\n", i)
		}
	}
}
