package utils

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestMsg1(t *testing.T) {
	SendMsg(true, "hello", "li", Green, false)
	SendMsg(false, "xxx", "yyyy", Red, true)
}

func TestMsg2(t *testing.T) {
	SendMsg(false, "hello", "li", Green, false)
	SendMsg(false, "xxx", "yyyy", Red, true)
}

//  go test -v -gcflags all="-N -l"  -run TestDotBar
func TestDotBar(t *testing.T) {
	for i := 10; i > 0; i-- {
		fmt.Printf("\x1B[2K%s", "\r.")
		//fmt.Printf("\x1B[A%s", "\r.")
		time.Sleep(time.Second / 2)
		fmt.Print("\r..")
		time.Sleep(time.Second / 2)
		fmt.Print("\r...")
		time.Sleep(time.Second / 2)
	}
}

/*
\33[2K erases the entire line your cursor is currently on
\033[A moves your cursor up one line, but in the same column i.e. not to the start of the line
\r brings your cursor to the beginning of the line (r is for carriage return N.B. carriage returns do not include a newline so cursor remains on the same line) but does not erase anything
如\33或\033，八进制33，十进制27，ASCII码表上是ESC，其实就是向左箭头
\x1b，十六进制1b，十进制27，也是左箭头
*/

func TestDotBar1(t *testing.T) {
	for i := 5; i > 0; i-- {
		//fmt.Printf("%s", "\r.  ")
		fmt.Printf("\x1B[?25l%s", "\r.  ")
		time.Sleep(time.Second / 4)
		fmt.Print("\r.. ")
		time.Sleep(time.Second / 4)
		fmt.Print("\r...")
		time.Sleep(time.Second / 4)
	}
}

//  go test -v -gcflags all="-N -l"  -run TestConsole
func TestConsole(t *testing.T) {
	s := NewSpinner("working...")
	isTTY := isTTY()
	for i := 0; i < 100; i++ {
		if isTTY {
			s.Tick()
		}
		time.Sleep(100 * time.Millisecond)
	}
}

//var spinChars = `|/-\`
var spinChars = `+-*/`

type Spinner struct {
	message string
	i       int
}

func NewSpinner(message string) *Spinner {
	return &Spinner{message: message}
}

func (s *Spinner) Tick() {
	fmt.Printf("%s %s %c \r", Hide(), s.message, spinChars[s.i])
	//fmt.Printf("%s %c \r", s.message, spinChars[s.i])
	s.i = (s.i + 1) % len(spinChars)
}

var Esc = "\x1b"

func escape(format string, args ...interface{}) string {
	return fmt.Sprintf("%s%s", Esc, fmt.Sprintf(format, args...))
}

// Show returns ANSI escape sequence to show the cursor
func Show() string {
	return escape("[?25h")
}

// Hide returns ANSI escape sequence to hide the cursor
func Hide() string {
	return escape("[?25l")
}

func isTTY() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}
