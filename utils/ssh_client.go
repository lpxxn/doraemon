package utils

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

type AuthMethod string

const (
	PublicKey AuthMethod = "publickey"
)

var AuthMethodList map[AuthMethod]struct{} = map[AuthMethod]struct{}{PublicKey: {}}

func VerifyAuthMethod(m AuthMethod) bool {
	_, ok := AuthMethodList[m]
	return ok
}

type SSHDial interface {
	Dial() (*ssh.Client, error)
	SSHConfig() *ssh.ClientConfig
	AuthMethodName() AuthMethod
}

type sshClient struct {
	Client       *ssh.Client
	proxyConf    *SSHPrivateKeyConfig
	logging      bool
	logTimestamp bool
	logFile      string
}

type SSHPrivateKeyConfig struct {
	MethodName  AuthMethod
	URI         string
	User        string
	AuthMethods []ssh.AuthMethod
	Timout      time.Duration
	Proxy       SSHDial
}

func (s *SSHPrivateKeyConfig) Dial() (*ssh.Client, error) {
	if s.Proxy == nil {
		return ssh.Dial("tcp", s.URI, s.SSHConfig())
	}
	proxyClient, err := s.Proxy.Dial()
	if err != nil {
		return nil, err
	}
	conn, err := proxyClient.Dial("tcp", s.URI)
	if err != nil {
		return nil, err
	}
	ncc, newCh, reqs, err := ssh.NewClientConn(conn, s.URI, s.SSHConfig())
	if err != nil {
		return nil, err
	}

	return ssh.NewClient(ncc, newCh, reqs), nil
}
func (s *SSHPrivateKeyConfig) SSHConfig() *ssh.ClientConfig {
	timeout := defaultTimeout
	if s.Timout > 0 {
		timeout = s.Timout
	}
	return &ssh.ClientConfig{
		User:            s.User,
		Auth:            s.AuthMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         timeout,
	}
}
func (s *SSHPrivateKeyConfig) AuthMethodName() AuthMethod {
	return s.MethodName
}

var (
	defaultTimeout = time.Second * 10
)

type SSHClientOption func(client *sshClient)

func ProxyConfig(sshOpts *SSHPrivateKeyConfig) SSHClientOption {
	return func(client *sshClient) {
		client.proxyConf = sshOpts
	}
}

func CreateSSHClient(conf *SSHPrivateKeyConfig, opts ...SSHClientOption) (*sshClient, error) {
	//uri := net.JoinHostPort(host, port)
	c := &sshClient{}
	timeout := defaultTimeout
	if conf.Timout > 0 {
		timeout = conf.Timout
	}
	targetSSHConfig := &ssh.ClientConfig{
		User:            conf.User,
		Auth:            conf.AuthMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         timeout,
	}
	for _, o := range opts {
		o(c)
	}
	var connClient *ssh.Client
	var err error
	if c.proxyConf == nil {
		if connClient, err = ssh.Dial("tcp", conf.URI, targetSSHConfig); err != nil {
			return nil, err
		}
	} else {
		timeout := defaultTimeout
		if conf.Timout > 0 {
			timeout = c.proxyConf.Timout
		}
		proxySSHConf := &ssh.ClientConfig{
			User:            c.proxyConf.User,
			Auth:            c.proxyConf.AuthMethods,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         timeout,
		}
		proxyClient, err := ssh.Dial("tcp", c.proxyConf.URI, proxySSHConf)
		if err != nil {
			return nil, err
		}
		conn, err := proxyClient.Dial("tcp", conf.URI)
		if err != nil {
			return nil, err
		}
		ncc, newCh, reqs, err := ssh.NewClientConn(conn, conf.URI, targetSSHConfig)
		if err != nil {
			return nil, err
		}

		connClient = ssh.NewClient(ncc, newCh, reqs)
	}

	c.Client = connClient
	return c, nil
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
