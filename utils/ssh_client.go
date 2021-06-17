package utils

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

type AuthMethod string

const (
	PublicKey AuthMethod = "publickey"
	Password  AuthMethod = "password"
)

var AuthMethodMap map[AuthMethod]struct{} = map[AuthMethod]struct{}{PublicKey: {}, Password: {}}

func VerifyAuthMethod(m AuthMethod) bool {
	_, ok := AuthMethodMap[m]
	return ok
}

type SSHConfig interface {
	SSHConfig() *ssh.ClientConfig
	AuthMethodName() AuthMethod
	GetProxy() SSHConfig
	SetProxy(proxyConfig SSHConfig)
	GetURI() string
	GetStartCommand() string
}

type SSHDial interface {
	NewSSHClient(conf SSHConfig) (*sshClient, error)
}

type sshClient struct {
	Client       *ssh.Client
	config       SSHConfig
	logging      bool
	logTimestamp bool
	logFile      string
}

type SSHBaseConfig struct {
	MethodName   AuthMethod
	URI          string
	User         string
	Passphrase   string
	AuthMethods  []ssh.AuthMethod
	Timout       time.Duration
	Proxy        SSHConfig
	StartCommand string
}

type SSHPrivateKeyConfig struct {
	*SSHBaseConfig
}

type SSHPasswordConfig struct {
	*SSHBaseConfig
}

func NewSSHClient(c SSHConfig) (*sshClient, error) {
	client, err := newSSHClient(c)
	if err != nil {
		return nil, err
	}
	return &sshClient{Client: client, config: c}, nil
}

func newSSHClient(c SSHConfig) (*ssh.Client, error) {
	if c.GetProxy() == nil {
		return ssh.Dial("tcp", c.GetURI(), c.SSHConfig())
	}
	proxyClient, err := newSSHClient(c.GetProxy())
	if err != nil {
		return nil, err
	}
	conn, err := proxyClient.Dial("tcp", c.GetURI())
	if err != nil {
		return nil, err
	}
	ncc, newCh, reqs, err := ssh.NewClientConn(conn, c.GetURI(), c.SSHConfig())
	if err != nil {
		return nil, err
	}

	return ssh.NewClient(ncc, newCh, reqs), nil
}

func (c *SSHBaseConfig) SSHConfig() *ssh.ClientConfig {
	timeout := defaultTimeout
	if c.Timout > 0 {
		timeout = c.Timout
	}
	return &ssh.ClientConfig{
		User:            c.User,
		Auth:            c.AuthMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         timeout,
	}
}
func (c *SSHBaseConfig) AuthMethodName() AuthMethod {
	return c.MethodName
}
func (c *SSHBaseConfig) GetProxy() SSHConfig {
	return c.Proxy
}

func (c *SSHBaseConfig) SetProxy(proxyConfig SSHConfig) {
	c.Proxy = proxyConfig
}

func (c *SSHBaseConfig) GetURI() string {
	return c.URI
}
func (c *SSHBaseConfig) GetStartCommand() string {
	return c.StartCommand
}

var (
	defaultTimeout = time.Second * 10
)

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

	if len(s.config.GetStartCommand()) > 0 {
		var buf bytes.Buffer
		buf.WriteString(s.config.GetStartCommand() + "\n")
		session.Stdin = io.MultiReader(&buf, session.Stdin)
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
	return s.RequestTty(session)
}
func (s *sshClient) RequestTty(session *ssh.Session) error {
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

func SetSttySane() {
	// https://github.com/c-bata/go-prompt/issues/233
	//rawModeOff := exec.Command("/bin/stty", "-raw", "echo")
	rawModeOff := exec.Command("/bin/stty", "sane")
	rawModeOff.Stdin = os.Stdin
	_ = rawModeOff.Run()
	rawModeOff.Wait()
}
