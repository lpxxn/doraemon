package ssh_utils

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

type sshClient struct {
	Client  *ssh.Client
	Session *ssh.Session
	Stdin   io.Reader
	Stdout  io.Writer
	Stderr  io.Writer

	logging      bool
	logTimestamp bool
	logFile      string
}

func CreateSSHClient(host, port, user string, authMethods []ssh.AuthMethod) (*sshClient, error) {
	uri := net.JoinHostPort(host, port)
	config := &ssh.ClientConfig{
		User:            user,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}
	conn, err := ssh.Dial("tcp", uri, config)

	if err != nil {
		return nil, err
	}
	return &sshClient{Client: conn}, nil
}

func (s *sshClient) SetLog(path string, timestamp bool) {
	s.logging = true
	s.logFile = path
	s.logTimestamp = timestamp
}

func (s *sshClient) CreateSession() (*ssh.Session, error) {
	return s.Client.NewSession()
}
func (s *sshClient) Shell(session *ssh.Session) (err error) {
	// Input terminal Make raw
	fd := int(os.Stdin.Fd())
	state, err := terminal.MakeRaw(fd)
	if err != nil {
		return
	}
	defer terminal.Restore(fd, state)

	// setup
	err = s.setupShell(session)
	if err != nil {
		return
	}

	// Start shell
	err = session.Shell()
	if err != nil {
		return
	}

	// keep alive packet
	// go s.SendKeepAlive(session)

	err = session.Wait()
	if err != nil {
		return
	}

	return
}

func (s *sshClient) setupShell(session *ssh.Session) error {
	// set FD
	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	// Logging
	if s.logging {
		err := s.logger(session)
		if err != nil {
			log.Println(err)
		}
	}
	// Request tty
	return RequestTty(session)
}
func RequestTty(session *ssh.Session) error {
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	// Get terminal window size
	fd := int(os.Stdin.Fd())
	width, height, err := terminal.GetSize(fd)
	if err != nil {
		return err
	}

	term := os.Getenv("TERM")
	if err = session.RequestPty(term, height, width, modes); err != nil {
		session.Close()
		return err
	}

	// Terminal resize goroutine.
	winch := syscall.Signal(0x1c)
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, winch)
	go func() {
		for {
			s := <-signalCh
			switch s {
			case winch:
				fd := int(os.Stdout.Fd())
				width, height, _ = terminal.GetSize(fd)
				session.WindowChange(height, width)
			}
		}
	}()

	return nil
}

func (s *sshClient) logger(session *ssh.Session) error {
	logfile, err := os.OpenFile(s.logFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	if s.logTimestamp {
		buf := new(bytes.Buffer)
		session.Stdout = io.MultiWriter(session.Stdout, buf)
		session.Stderr = io.MultiWriter(session.Stderr, buf)

		go func() {
			preLine := []byte{}
			for {
				if buf.Len() > 0 {
					line, err := buf.ReadBytes('\n')

					if err == io.EOF {
						preLine = append(preLine, line...)
						continue
					} else {
						timestamp := time.Now().Format("2006/01/02 15:04:05 ") // yyyy/mm/dd HH:MM:SS
						fmt.Fprintf(logfile, timestamp+string(append(preLine, line...)))
						preLine = []byte{}
					}
				} else {
					time.Sleep(10 * time.Millisecond)
				}
			}
		}()

	} else {
		session.Stdout = io.MultiWriter(session.Stdout, logfile)
		session.Stderr = io.MultiWriter(session.Stderr, logfile)
	}

	return nil
}
