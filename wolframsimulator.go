package main

import (
	"github.com/nsf/termbox-go"
	"errors"
	"flag"
	"log"
	"os"
)

type row []byte

var displayHeight, displayWidth, ruleId int
var logging bool

const (
	FgColor = termbox.ColorYellow
	BgColor = termbox.ColorBlue
)

func handleError(err error, logThis bool) {
	if logThis {
		log.Println(err.Error())
	}
	panic(err)
}

// sets middle element to 1, all others 0
func getDeterminantRow() row {
	rowWidth := 2*(displayHeight-1) + displayWidth
	rowMiddle := rowWidth / 2

	r := make(row, rowWidth, rowWidth)
	r[rowMiddle] = 1

	return r
}

func getDisplayRune(value byte) (displayString rune) {
	if value == 0 {
		displayString = ' '
	} else {
		displayString = '#'
	}

	return
}

// centers simulation
func getDisplayBounds() (upper, lower, left, right int) {
	tWidth, tHeight := termbox.Size()

	upper = (tHeight / 2) - (displayHeight / 2)
	lower = upper + displayHeight

	left = (tWidth / 2) - (displayWidth / 2)
	right = left + displayWidth

	return
}

func drawRow(y int, offset int, r row) error {
	width, height := termbox.Size()

	if y < 0 || y >= height {
		return errors.New("Row out of range")
	} else if len(r) > width {
		return errors.New("Row too wide")
	}

	for i, val := range r {
		termbox.SetCell(i+offset, y, getDisplayRune(val), FgColor, BgColor)
	}

	termbox.Flush()

	return nil
}

// validates ruleId, displayHeight, displayWidth
func validateParameters(ignoreRule bool) error {
	tWidth, tHeight := termbox.Size()

	const (
		invRuleId        = "Invalid ruleId: 0 <= ruleId < 256"
		invDisplayWidth  = "Invalid displayWidth: 0 < displayWidth <= tWidth"
		invDisplayHeight = "Invalid displayHeight: 0 < displayHeight <= tHeight"
	)

	switch {
	case displayWidth <= 0 || displayWidth > tWidth:
		return errors.New(invDisplayWidth)
	case displayHeight <= 0 || displayHeight > tHeight:
		return errors.New(invDisplayHeight)
	case !ignoreRule && (ruleId < 0 || ruleId > 256):
		return errors.New(invRuleId)
	}

	return nil
}

func applyRule(input []byte) byte {
	output := byte(ruleId)

	for i, n := range input {
		output >>= (n & 1) << uint(i)
	}

	return output & 1
}

func getNextRow(previous row) row {
	rowWidth := len(previous) - 2
	next := make(row, rowWidth, rowWidth)

	for i := 0; i < rowWidth; i++ {
		next[i] = applyRule(previous[i : i+3])
	}

	return next
}

func getCenter(r row) row {
	left := (len(r) / 2) - (displayWidth / 2)

	return r[left : left+displayWidth]
}

func runSimulation() {
	currentRow := getDeterminantRow()
	upper, lower, left, _ := getDisplayBounds()

	for y := upper; y < lower; y++ {
		drawRow(y, left, getCenter(currentRow))
		currentRow = getNextRow(currentRow)
	}
}

func init() {
	if err := termbox.Init(); err != nil {
		log.Fatal(err)
	}

	const (
		usgDisplayHeight = "height of simulation"
		usgDisplayWidth  = "width of simulation"
		usgRuleId        = "rule to be simulated"
		usgSimAllRules   = "display all rules, one every second"
		usgLogging       = "create log file"
	)

	defaultWidth, defaultHeight := termbox.Size()
	defaultRule := 30

	// -height, -h ==> displayHeight
	flag.IntVar(&displayHeight, "height", defaultHeight, usgDisplayHeight)
	flag.IntVar(&displayHeight, "h", defaultHeight, usgDisplayHeight+" (abbr. of -height)")

	// -width, -w ==> displayWidth
	flag.IntVar(&displayWidth, "width", defaultWidth, usgDisplayWidth)
	flag.IntVar(&displayWidth, "w", defaultWidth, usgDisplayWidth+" (abbr. of -width)")

	// -rule, -r ==> ruleId
	flag.IntVar(&ruleId, "rule", defaultRule, usgRuleId)
	flag.IntVar(&ruleId, "r", defaultRule, usgRuleId+" (abbr. of -rule)")

	// -all, -a ==> display all rules, one every second
	var simAllRules bool
	flag.BoolVar(&simAllRules, "all", false, usgSimAllRules)
	flag.BoolVar(&simAllRules, "a", false, usgSimAllRules+" (abbr. of -all)")

	// -log, -l ==> create log file
	flag.BoolVar(&logging, "log", false, usgLogging)
	flag.BoolVar(&logging, "l", false, usgLogging+" (abbr of -log)")

	flag.Parse()

	if err := validateParameters(simAllRules); err != nil {
		termbox.Close()
		handleError(err, false)
	}
}

func main() {
	defer termbox.Close()

	if logging {
		logFile, err := os.OpenFile("wsim_log.txt", os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			handleError(errors.New("Log file initialization error"), false)
		}
		log.SetOutput(logFile)
		log.Println("[Wolfram Simulator Log File]")

		defer logFile.Close()
	}

	runSimulation()

	//wait for any keypress to continue
	for {
		if currentEvent := termbox.PollEvent(); currentEvent.Type == termbox.EventError {
			handleError(currentEvent.Err, logging)
		} else if currentEvent.Type == termbox.EventKey {
			break
		}
	}
}
