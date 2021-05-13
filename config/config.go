package config

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/lpxxn/doraemon/utils"
	"golang.org/x/crypto/ssh"
)

type appConfig struct {
	SSHInfo   []*sshInfo   `toml:"sshInfo"`
	LoginInfo []*loginInfo `toml:"loginInfo"`
}

var LoginConf *appConfig

func (a *appConfig) ConfigByName(name string) (*sshInfo, error) {
	for _, item := range a.SSHInfo {
		if item.Name == name {
			return item, nil
		}
	}
	return nil, os.ErrNotExist
}

type sshInfo struct {
	Name          string        `toml:"name"`
	AuthMethod    string        `toml:"authMethod"`
	URI           string        `toml:"uri"`
	User          string        `toml:"user"`
	PublicKeyPath string        `toml:"publicKeyPath"`
	Timout        time.Duration `toml:"timout"`
	ProxySSHName  string        `toml:"proxySSHName"`
	Desc          string        `toml:"desc"`
}

type loginInfo struct {
	URL          string `toml:"url"`
	ClientID     string `toml:"clientID"`
	ClientSecret string `toml:"clientSecret"`
	Name         string `toml:"name"`
	Pwd          string `toml:"pwd"`
	PwdUseMin    bool   `toml:"pwdUseMin"`
}

func (s *sshInfo) ToSSHConfig() (*utils.SSHConfig, error) {
	sshConf := &utils.SSHConfig{
		AuthMethodName: utils.AuthMethod(s.AuthMethod),
		URI:            s.URI,
		User:           s.User,
		AuthMethods:    nil,
		Timout:         s.Timout,
	}
	if sshConf.AuthMethodName == utils.PublicKey {
		pemBytes, err := ioutil.ReadFile(s.PublicKeyPath)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		signer, err := ssh.ParsePrivateKey(pemBytes)
		if err != nil {
			log.Fatalf("parse key failed:%v", err)
			return nil, err
		}
		sshConf.AuthMethods = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	}
	return sshConf, nil
}

func (s *sshInfo) HaveProxy() bool {
	return len(s.ProxySSHName) > 0
}

func ParseConfig() error {
	f, err := GetConfig()
	if err != nil {
		return err
	}
	if LoginConf == nil {
		LoginConf = new(appConfig)
	}
	_, err = toml.DecodeReader(f, LoginConf)
	return err
}
