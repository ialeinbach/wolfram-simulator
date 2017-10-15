package main

import (
	"fmt"
	"time"
	"errors"
	"log"
	"flag"
	"github.com/nsf/termbox-go"
)

// minimum height of terminal required
const MIN_HEIGHT = 4

// receive values from flags
var nRows, displayWidth, ruleId int

// hold terminal dimensions
var terminalWidth, terminalHeight int

// Because each bit depends on the above bit and those adjacent, a
// row of width n is determined by a row of width n+2. Given
// nRows, arrWidth holds size of initial row needed to completely
// determine nRows of width displayWidth. arrMiddle holds middle index
// of arrWidth
func getSizeInfo() (arrWidth, arrMiddle int) {
	// offset arrWidth by displayWidth so
	// the last row is displayWidth wide
	//
	// displayWidth guaranteed to be odd, so arrWidth will
	// always be odd i.e. have a middle element
	arrWidth = 2 * (nRows - 1) + displayWidth
	arrMiddle = (arrWidth - 1) / 2
	return
}

// return printable character to represent a bit value
func getDisplayString(value byte) (displayString string) {
	if value == 0 {
		displayString = " "
	} else {
		displayString = "#"
	}

	return
}

// rules encoded as single byte where each bit
// represents output resulting from input of a
// bit's index (0-7) when expressed as 3-bit
// binary number (000-111)
func applyRule(input, rule byte) byte {
	// apply bit mask over bottom three bits of input
	// to prevent over-shifting
	return (rule >> (input & 7)) & 1
}

// initial row has all 0s except a
// single 1 in the middle
func generateInitialRow() (initialRow []byte) {
	width, middle := getSizeInfo()

	initialRow = make([]byte, width)
	initialRow[middle] = 1

	return
}

// returns indices start, end such that row[start:end] contains
// middle displayWidth elements of entire row
//
// rows are guaranteed to have odd length, so there is
// always a middle subset of elements
func getDisplayBounds(arrWidth int) (start, end int) {
	start = (arrWidth - displayWidth) / 2
	end = start + displayWidth

	return
}

// given a slice of a row, insert into a byte
// to be processed by applyRule()
func parseInput(previousInput []byte) byte {
	length := len(previousInput)

	if length != 3 {
		log.Fatal(errors.New("parseInput given slice with invalid length"))
	}

	// initialized to 0
	var input byte

	for i := 0; i < length; i++ {
		input += previousInput[i]
		input <<= 1
	}

	return input >> 1
}

// generates new row based on previous row and rule
func generateNextRow(previousRow []byte, rule byte) (nextRow []byte) {
	arrWidth := len(previousRow) - 2
	nextRow = make([]byte, arrWidth)

	// iterate across previousRow
	// generate nextRow according to rule
	for i := 0; i < arrWidth; i++ {
		nextRow[i] = applyRule(parseInput(previousRow[i:i+3]), rule)
	}

	return
}

// print at most the central displayWidth elements of a row
func printRow(row []byte) {
	arrWidth := len(row)
	i, end := getDisplayBounds(arrWidth)

	for ; i < end; i++ {
		fmt.Printf("%s", getDisplayString(row[i]))
	}

	fmt.Println()
}

// print simple label for a rule
func printRuleLabel() {
	fmt.Printf(".------------------.\n")
	fmt.Printf("|     Rule %3d     |\n", ruleId)
	fmt.Printf("`------------------'\n")
}

// validates ruleId, nRows, and displayWidth
func validateParameters() error {
	//ruleId = 256 is a special value, see main() for behavior
	if ruleId < 0 || ruleId > 256 {
		return errors.New("Invalid ruleId: 0 <= ruleId < 256 or rule = 256 to print every rule, one every second")
	} else if displayWidth > terminalWidth {
		return errors.New(fmt.Sprintf("Invalid displayWidth: displayWidth <= terminalWidth(%d)", terminalWidth))
	} else if displayWidth < 0 {
		return errors.New("Invalid displayWidth: displayWidth > 0")
	} else if nRows <= 0 {
		return errors.New("Invalid nRows: nRows > 0")
	}

	return nil
}

// simulates rule ruleId for nRows rows with a width of displayWidth
// note: assumes validated parameters (ruleId, nRows, displayWidth)
//   \__
//      `--> validateParameters() called in init()
func displayRule() {
	rule := byte(ruleId)
	printRuleLabel()

	currentRow := generateInitialRow()
	printRow(currentRow)

	for i := 1; i < nRows; i++ {
		currentRow = generateNextRow(currentRow, rule)
		printRow(currentRow)
	}
}

// clear terminal and detect terminal dimensions
func initializeTerminal() error {
	// termbox.Init() also clears terminal
	if err := termbox.Init(); err != nil {
		return err
	}

	// only need to close if *succesfully* initialized
	defer termbox.Close()

	terminalWidth, terminalHeight = termbox.Size()

	if terminalHeight < MIN_HEIGHT {
		return errors.New(fmt.Sprintf("Terminal window height must be at least %d.", MIN_HEIGHT))
	}

	return nil
}

func init() {
	// also detects terminal dimensions
	if err := initializeTerminal(); err != nil {
		log.Fatal(err)
	}

	// usage messages for flags
	const (
		usgNRows = "number of rows for which to simulate rule(s)"
		usgDisplayWidth = "width of rows"
		usgRuleId = "valid rule to be simulated or 256 to display all rules, one every second"
	)

	// define flags
	flag.IntVar(&nRows, "rows", terminalHeight - MIN_HEIGHT, usgNRows)
	flag.IntVar(&displayWidth, "width", terminalWidth, usgDisplayWidth)
	flag.IntVar(&ruleId, "rule", 256, usgRuleId)

	flag.Parse()

	// ensure valid parameters for displayRule()
	if err := validateParameters(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	// special value used to print all rules, one every second
	if ruleId == 256 {
		for ruleId = 0; ruleId < 256; ruleId++ {
			displayRule()
			time.Sleep(time.Second)
		}
	} else {
		displayRule()
	}
}
