package internal

import (
	"fmt"
	"os"

	"github.com/skratchdot/open-golang/open"
)

type Color string

const (
	Red     Color = "red"
	Blue          = "blue"
	Green         = "green"
	Cyan          = "cyan"
	Yellow        = "yellow"
	NoColor       = ""
)

// BeautifyText function for send (colored or common) message to output.
func BeautifyText(text string, color Color) string {
	// Define variables.
	var (
		red       string = "\033[0;31m"
		green     string = "\033[0;32m"
		cyan      string = "\033[0;36m"
		yellow    string = "\033[1;33m"
		blue      string = "\033[0;34m"
		noColor   string = "\033[0m"
		textColor string
	)

	// Switch color.
	switch color {
	case NoColor:
		textColor = noColor
	case Blue:
		textColor = blue
	case Green:
		textColor = green
	case Yellow:
		textColor = yellow
	case Red:
		textColor = red
	case Cyan:
		textColor = cyan
	}

	// Send common or colored text.
	return textColor + text + noColor
}

// SendMsg function for send message to output.
func SendMsg(startWithNewLine bool, caption, text string, color Color, endWithNewLine bool) {
	// Define variables.
	var startNewLine, endNewLine string

	if startWithNewLine {
		startNewLine = "\n" // set new line
	}

	if endWithNewLine {
		endNewLine = "\n" // set new line
	}

	if caption == "" {
		fmt.Println(startNewLine + text + endNewLine) // common text
	} else {
		fmt.Println(startNewLine + BeautifyText(caption, color) + " " + text + endNewLine) // colorized text
	}
}

func MakeFolder(folderName string, chmod os.FileMode) error {
	if _, err := os.Stat(folderName); os.IsNotExist(err) {
		return os.MkdirAll(folderName, chmod)
	} else {
		return err
	}
}

func OpenFolder(folderName string) error {
	return open.Run(folderName)
}
