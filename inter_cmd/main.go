package main

import (
	"fmt"
	"os"
	"time"

	"github.com/lpxxn/doraemon/utils"
)

func main() {
	utils.SendMsg(true, "hello", "li", utils.Green, false)
	//for i := 1000; i > 0; i-- {
	//	fmt.Print("\r\\\r")
	//	time.Sleep(time.Second / 2)
	//	fmt.Print("\r|\r")
	//	time.Sleep(time.Second / 2)
	//	fmt.Print("\r/\r")
	//	time.Sleep(time.Second / 2)
	//}

	for i := 1000; i > 0; i-- {
		fmt.Printf("\x1B[?25l%s", "\r.  ")
		time.Sleep(time.Second / 2)
		fmt.Print("\r.. ")
		time.Sleep(time.Second / 2)
		fmt.Print("\r...")
		time.Sleep(time.Second / 2)
	}

}

var spinChars = `|/-\`

type Spinner struct {
	message string
	i       int
}

func NewSpinner(message string) *Spinner {
	return &Spinner{message: message}
}

func (s *Spinner) Tick() {
	fmt.Printf("%s %c \r", s.message, spinChars[s.i])
	s.i = (s.i + 1) % len(spinChars)
}

func isTTY() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}
