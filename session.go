package executor

import (
	"bytes"
	"fmt"
	"time"
	"unicode"

	//"os"
	"strings"

	//"github.com/spf13/cobra/cobra/cmd"
	"golang.org/x/crypto/ssh"
	"net"
	"os/exec"
)

type Session interface {
	Run(cmd string) error
	Close() error
	Output() string
	ErrMsg() string
}

type SessionConfig struct {
	User     string
	Password string
	Host     string
	Port     int8
	Timeout  int
}

type localSession struct {
	cmd    *exec.Cmd
	msgBuf *bytes.Buffer
	errBuf *bytes.Buffer
}

type remoteSession struct {
	session *ssh.Session
	msgBuf  *bytes.Buffer
	errBuf  *bytes.Buffer
}

func newRemoteSession(config SessionConfig) func() (Session, error) {
	return func() (Session, error) {
		r := remoteSession{
			msgBuf: &bytes.Buffer{},
			errBuf: &bytes.Buffer{},
		}
		session, err := sshConnect(config.Host, config.Port, config.User, config.Password, config.Timeout)
		if err != nil {
			return nil, err
		}
		r.session = session
		r.session.Stdout = r.msgBuf
		r.session.Stderr = r.errBuf
		return r, nil
	}
}

func newLocalSession() func() (Session, error) {
	return func() (Session, error) {
		r := localSession{
			cmd:    &exec.Cmd{},
			msgBuf: &bytes.Buffer{},
			errBuf: &bytes.Buffer{},
		}
		r.cmd.Stdout = r.msgBuf
		r.cmd.Stderr = r.errBuf
		return r, nil
	}
}

func (r localSession) Run(command string) error {
	args := strings.FieldsFunc(command, unicode.IsSpace)
	r.cmd.Path = args[0]
	r.cmd.Args = args
	if err := r.cmd.Run(); err != nil {
		r.errBuf.WriteString(err.Error())
		return err
	}
	return nil
}
func (r localSession) Close() error {
	return nil
}
func (r localSession) Output() string {
	return r.msgBuf.String()
}
func (r localSession) ErrMsg() string {
	return r.errBuf.String()
}

func (r remoteSession) Run(command string) error {
	err := r.session.Run(command)
	if err != nil {
		r.errBuf.WriteString(err.Error())
	}
	return err
}

func (r remoteSession) Close() error {
	return r.session.Close()
}
func (r remoteSession) Output() string {
	return r.msgBuf.String()
}
func (r remoteSession) ErrMsg() string {
	return r.errBuf.String()
}

func sshConnect(host string, port int8, user, password string, timeout int) (*ssh.Session, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		session      *ssh.Session
		err          error
	)
	if timeout == 0 {
		timeout = 20000
	}
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	hostKeyCb := func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		return nil
	}

	clientConfig = &ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		Timeout:         time.Duration(timeout) * time.Second,
		HostKeyCallback: hostKeyCb,
	}

	addr = fmt.Sprintf("%s:%d", host, port)

	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}
	// create session
	if session, err = client.NewSession(); err != nil {
		return nil, err
	}
	return session, nil
}
