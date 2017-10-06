package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"errors"
	"log"
)


// Because each bit depends on the above bit and those adjacent, each
// row of width n is determined by a row of width n+2. Given
// nRows, arrWidth holds size of initial row needed to completely
// determine nRows of width dispWidth. arrMiddle holds middle index
func getSizeInfo(nRows, dispWidth int) (arrWidth, arrMiddle int) {
        // offset arrWidth by dispWidth so the last row is
        // dispWidth wide rather than 0 wide
        //
        // dispWidth guaranteed to be odd, so arrWidth will
        // always be odd i.e. have a middle element
	arrWidth = 2 * (nRows - 1) + dispWidth
	arrMiddle = (arrWidth - 1) / 2
	return
}

//Return printable character to display
func getDispString(value byte) (dispString string) {
	if value == 0 {
		dispString = " "
	} else {
		dispString = "#"
	}

	return
}

// Rule encoded as single byte where each bit
// represents output for an input of a bit's index (0-7)
// when expressed as 3-bit binary number (000-111)
func applyRule(input, rule byte) byte {
	return (rule >> (input & 7)) & 1
}

// initial row has all 0s except a single 1 in the
// middle to start
func initialRow(nRows, dispWidth int) []byte {
	width, middle := getSizeInfo(nRows, dispWidth)
	firstRow := make([]byte, width)
	firstRow[middle] = 1

	return firstRow
}

// returns indices start, end such that row[start:end] contains
// middle dispWidth elements of entire row
//
// rows are guaranteed to have odd length, so there is
// always a middle subset of elements
func getDispBounds(arrWidth, dispWidth int) (start, end int) {
	start = (arrWidth - dispWidth) / 2
	end = start + dispWidth

	return
}

// given a slice of a row, insert into a byte
// to be processed by applyRule
func parseInput(previousInput []byte) byte {
	length := len(previousInput)

	if length != 3 {
		log.Fatal(errors.New("parseInput given slice of invalid length"))
	}

	var input byte

	for i := 0; i < length; i++ {
		input += previousInput[i]
		input <<= 1
	}

	return input >> 1
}

// generates new row based on previous row and rule
func generateRow(previousRow []byte, rule byte) []byte {
	arrWidth := len(previousRow) - 2
	currentRow := make([]byte, arrWidth)

        // row sizes are such that expressions with i will not be out of bounds
	for i := 0; i < arrWidth; i++ {
		currentRow[i] = applyRule(parseInput(previousRow[i:i+3]), rule)
	}

	return currentRow
}

// prints row
func printRow(row []byte, dispWidth int) {
	length := len(row)
	i, end := getDispBounds(length, dispWidth)

	for ; i < end; i++ {
		fmt.Printf("%v", getDispString(row[i]))
	}

	fmt.Println()
}

// print simple label for a given rule
func printRuleLabel(ruleId byte) {
	fmt.Printf("\n====================\n")
	fmt.Printf("||    Rule %3d    ||\n", ruleId)
	fmt.Printf("====================\n\n")
}

// simulates rule ruleId for nRows rows of width dispWidth
func dispRule(ruleId, nRows, dispWidth int) {
	rule := byte(ruleId)
	printRuleLabel(rule)

	currentRow := initialRow(nRows, dispWidth)
	printRow(currentRow, dispWidth)

	for i := 1; i < nRows; i++ {
		currentRow = generateRow(currentRow, rule)
		printRow(currentRow, dispWidth)
	}
}

func getTerminalDimensions() (int, int, error) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, cmdErr := cmd.Output()

	if cmdErr != nil {
		return 0, 0, cmdErr
	}

	dimensions := strings.Split(strings.TrimSpace(string(out)), " ")

	height, atoiErr1 := strconv.Atoi(dimensions[0])
	width, atoiErr2 := strconv.Atoi(dimensions[1])

	if atoiErr1 != nil {
		return 0, 0, atoiErr1
	}

	if atoiErr2 != nil {
		return 0, 0, atoiErr2
	}

	return width, height, nil
}

func main() {
	args := os.Args
	numArgs := len(args)

	width, height, err := getTerminalDimensions()

	if err != nil {
		log.Fatal(err)
	}

	if height < 6 {
		log.Fatal(errors.New(fmt.Sprintf("Terminal window height (%d) must be at least 6.\n", height)))
	}

	// initialize to invalid values
        // if unchanged, will be invalid
	nRows, ruleId, dispWidth := -1, -1, -1

	switch numArgs {
		// numArgs == 1 (i.e. only program name)
		// print every rule, one per second
		case 1:
			for i := 0; i < 256; i++ {
				dispRule(i, height - 5, width)
				time.Sleep(time.Second)
			}

			// successful exit
			os.Exit(0)

		// numArgs > 4
		default:
			log.Fatal(errors.New("Too many arguments."))

                // will extract only arguments that exist using fallthroughs

		// numArgs == 4 --> extract all except program name from os.args
		case 4:
			dispWidth, err = strconv.Atoi(os.Args[3])

			if err != nil {
				log.Fatal(err)
			}

			fallthrough

		// numArgs == 3 --> extract all except program name, dispWidth from os.args
		case 3:
			nRows, err = strconv.Atoi(os.Args[2])

			if err != nil {
				log.Fatal(err)
			}

			fallthrough

		// numArgs == 2 --> extract all except program name, dispWidth, nRows from os.args
		case 2:
			ruleId, err = strconv.Atoi(os.Args[1])

			if err != nil {
				log.Fatal(err)
			}
	}

	// check if within range (also implicitly whether they
	// were assigned values to begin with)
	//
	// if invalid(if unassigned, guaranteed to be invalied),
        // use known valid values
	if dispWidth > 0 {
            if dispWidth % 2 == 0 {
		dispRule(ruleId, nRows, dispWidth)
            } else {
                dispRule(ruleId, nRows, dispWidth - 1)
            }
	} else if nRows > 0 {
		dispRule(ruleId, nRows, width)
	} else if ruleId >= 0 {
		dispRule(ruleId, height - 6, width)
	}

}
