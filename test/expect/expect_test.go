package expect

import (
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/Netflix/go-expect"
	"github.com/stretchr/testify/assert"
)

func TestExpect1(t *testing.T) {
	c, err := expect.NewConsole(expect.WithStdout(os.Stdout))
	assert.Nil(t, err)
	cmd := exec.Command("vi", "a.txt")
	cmd.Stdin = c.Tty()
	cmd.Stdout = c.Tty()
	cmd.Stderr = c.Tty()
	go func() {
		rev, err := c.ExpectString("hello world")
		assert.Nil(t, err)
		t.Log(rev)
	}()
	err = cmd.Start()
	assert.Nil(t, err)
	time.Sleep(time.Second)
	c.Send("ihello world")
	time.Sleep(time.Second)
	c.Send("dd")
	time.Sleep(time.Second)
	c.SendLine(":q!")

	err = cmd.Wait()
	assert.Nil(t, err)
	time.Sleep(time.Second * 10)
}
