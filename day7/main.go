package main

/*
--- Day 7: Some Assembly Required ---

This year, Santa brought little Bobby Tables a set of wires and bitwise logic gates! Unfortunately, little Bobby is a little under the recommended age range, and he needs help assembling the circuit.

Each wire has an identifier (some lowercase letters) and can carry a 16-bit signal (a number from 0 to 65535). A signal is provided to each wire by a gate, another wire, or some specific value. Each wire can only get a signal from one source, but can provide its signal to multiple destinations. A gate provides no signal until all of its inputs have a signal.

The included instructions booklet describes how to connect the parts together: x AND y -> z means to connect wires x and y to an AND gate, and then connect its output to wire z.

For example:

123 -> x means that the signal 123 is provided to wire x.
x AND y -> z means that the bitwise AND of wire x and wire y is provided to wire z.
p LSHIFT 2 -> q means that the value from wire p is left-shifted by 2 and then provided to wire q.
NOT e -> f means that the bitwise complement of the value from wire e is provided to wire f.
Other possible gates include OR (bitwise OR) and RSHIFT (right-shift). If, for some reason, you'd like to emulate the circuit instead, almost all programming languages (for example, C, JavaScript, or Python) provide operators for these gates.

For example, here is a simple circuit:

123 -> x
456 -> y
x AND y -> d
x OR y -> e
x LSHIFT 2 -> f
y RSHIFT 2 -> g
NOT x -> h
NOT y -> i
After it is run, these are the signals on the wires:

d: 72
e: 507
f: 492
g: 114
h: 65412
i: 65079
x: 123
y: 456
In little Bobby's kit's instructions booklet (provided as your puzzle input), what signal is ultimately provided to wire a?
*/

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

type eOpcode int

const (
	cOPASSIGN  eOpcode = 0
	cOPAND     eOpcode = 1
	cOPOR      eOpcode = 2
	cOPXOR     eOpcode = 3
	cOPNOT     eOpcode = 4
	cOPLSHIFT  eOpcode = 5
	cOPRSHIFT  eOpcode = 6
	cOPILLEGAL eOpcode = 255
)

type iInstruction struct {
	lineNumber  int
	opcode      eOpcode
	lhsOperand  string
	lhsRegister string
	lhsValue    uint16
	rhsOperand  string
	rhsRegister string
	rhsValue    uint16

	resRegister string
}

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

func decodeOpcode(opcodestr string) (opcode eOpcode, isopcode bool) {
	isopcode = true
	if opcodestr == "AND" {
		opcode = cOPAND
	} else if opcodestr == "OR" {
		opcode = cOPOR
	} else if opcodestr == "XOR" {
		opcode = cOPXOR
	} else if opcodestr == "LSHIFT" {
		opcode = cOPLSHIFT
	} else if opcodestr == "RSHIFT" {
		opcode = cOPRSHIFT
	} else if opcodestr == "NOT" {
		opcode = cOPNOT
	} else {
		opcode = cOPILLEGAL
		isopcode = false
	}
	return
}

func decodeNumber(numberstr string) (number uint16, isnumber bool) {
	isnumber = true
	number = 0
	for _, char := range numberstr {
		if unicode.IsNumber(char) == false {
			isnumber = false
			break
		}
		number = (number*uint16(10) + uint16(char-'0'))
	}
	return
}

func isRegister(str string) (isreg bool) {
	isreg = true
	for _, char := range str {
		if unicode.IsLetter(char) == false {
			isreg = false
			return
		}
	}
	return
}

func decodeInstruction(instructionstr string, lineNumber int) (valid bool, instruction iInstruction) {
	// Default instruction
	instruction.lineNumber = lineNumber
	instruction.opcode = cOPILLEGAL
	instruction.lhsOperand = string("illegal")
	instruction.lhsValue = 0
	instruction.lhsRegister = string("?")
	instruction.rhsOperand = string("illegal")
	instruction.rhsValue = 0
	instruction.rhsRegister = string("?")
	instruction.resRegister = string("UNKNOWN")

	valid = false
	fields := strings.Fields(instructionstr)
	operand := "lhs"
	for _, field := range fields {
		if field == "" || field == " " {
			continue
		}

		if field == "->" {
			operand = string("res")
		} else {
			opcode, isopcode := decodeOpcode(field)
			if isopcode {
				//fmt.Printf(" OPERAND(%s)", field)
				operand = string("rhs")
				instruction.opcode = opcode
			} else {
				number, isnumber := decodeNumber(field)
				if isnumber {
					if operand == "lhs" {
						//fmt.Printf(" LHS(%s,%d)", field, number)
						instruction.lhsOperand = "value"
						instruction.lhsValue = number
						instruction.lhsRegister = "NA"
						instruction.opcode = cOPASSIGN
						operand = "opcode"
					} else if operand == "rhs" {
						//fmt.Printf(" RHS(%s,%d)", field, number)
						instruction.rhsOperand = "value"
						instruction.rhsValue = number
						instruction.rhsRegister = "NA"
						operand = "->"
					}
				} else if isRegister(field) {
					if operand == "lhs" {
						//fmt.Printf(" LHS-REG(%s,%x)", field, register)
						instruction.lhsOperand = "reg"
						instruction.lhsValue = 0
						instruction.lhsRegister = field
						instruction.opcode = cOPASSIGN
						operand = "opcode"
					} else if operand == "rhs" {
						//fmt.Printf(" RHS-REG(%s,%x)", field, register)
						instruction.rhsOperand = "reg"
						instruction.rhsValue = 0
						instruction.rhsRegister = field
						operand = "->"
					} else if operand == "res" {
						//fmt.Printf(" RES-REG(%s)", field)
						instruction.resRegister = field
					}
				}
			}
		}
	}
	//fmt.Print("\n")
	valid = instruction.opcode != cOPILLEGAL && instruction.resRegister != ""
	return
}

func executeInstruction(i iInstruction, regmap map[string]uint16) (executed bool) {
	// Excute the opcode and write to result register
	executed = false

	lhsDescr := "value"
	lhs := i.lhsValue
	if i.lhsOperand == "reg" {
		if _, ok := regmap[i.lhsRegister]; ok == false {
			return
		}
		lhsDescr = i.lhsRegister
		lhs = regmap[i.lhsRegister]
	}

	rhsDescr := "value"
	rhs := i.rhsValue
	if i.rhsOperand == "reg" {
		if _, ok := regmap[i.rhsRegister]; ok == false {
			return
		}
		rhsDescr = i.rhsRegister
		rhs = regmap[i.rhsRegister]
	}

	fmt.Printf("LINE(%d): ", i.lineNumber)
	res := uint16(0)
	switch i.opcode {
	case cOPASSIGN:
		res = lhs
		fmt.Printf("RES(%s=%d) = %d(%s)\n", i.resRegister, res, lhs, lhsDescr)
	case cOPOR:
		res = lhs | rhs
		fmt.Printf("RES(%s=%d) = %d(%s) OR %d(%s))\n", i.resRegister, res, lhs, lhsDescr, rhs, rhsDescr)
	case cOPAND:
		res = lhs & rhs
		fmt.Printf("RES(%s=%d) = %d(%s) AND %d(%s))\n", i.resRegister, res, lhs, lhsDescr, rhs, rhsDescr)
	case cOPXOR:
		res = lhs ^ rhs
		fmt.Printf("RES(%s=%d) = %d(%s) XOR %d(%s))\n", i.resRegister, res, lhs, lhsDescr, rhs, rhsDescr)
	case cOPLSHIFT:
		res = lhs << rhs
		fmt.Printf("RES(%s=%d) = %d(%s) LSHIFT %d(%s))\n", i.resRegister, res, lhs, lhsDescr, rhs, rhsDescr)
	case cOPRSHIFT:
		res = lhs >> rhs
		fmt.Printf("RES(%s=%d) = %d(%s) RSHIFT %d(%s))\n", i.resRegister, res, lhs, lhsDescr, rhs, rhsDescr)
	case cOPNOT:
		res = rhs ^ uint16(0xFFFF)
		fmt.Printf("RES(%s=%d) = NOT %d (%s))\n", i.resRegister, res, rhs, rhsDescr)
	case cOPILLEGAL:
		res = 0
		fmt.Print("ILLEGAL OPCODE\n")
	}

	regmap[i.resRegister] = res
	executed = true
	return
}

func applySignalsFromFile(filename string) (wire uint16) {
	wire = 0

	var codemap = make([]iInstruction, 0)

	computator := func(text string, line int) {
		valid, instruction := decodeInstruction(text, line)
		if valid {
			codemap = append(codemap, instruction)
		} else {
			fmt.Printf("INVALID instruction at line %d", line)
		}
	}
	iterateOverLinesInTextFile(filename, computator)

	// Execute code until all instructions can be evaluated
	regmap := make(map[string]uint16)
	exemap := make(map[int]bool)

	fullyEvaluated := false
	for !fullyEvaluated {
		fullyEvaluated = true
		for _, i := range codemap {
			if _, executed := exemap[i.lineNumber]; executed == false {
				if executeInstruction(i, regmap) {
					exemap[i.lineNumber] = true
				} else {
					fullyEvaluated = false
				}
			}
		}
	}

	wire, _ = regmap["a"]
	return
}

func main() {
	var valueOnWire = applySignalsFromFile("input.text")
	fmt.Printf("The signal on the wire is: %v", valueOnWire)
}
