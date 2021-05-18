package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/lpxxn/doraemon/utils"
	"golang.org/x/crypto/ssh"
)

var ()

type appConfig struct {
	SSHInfo    []*sshInfo   `toml:"sshInfo"`
	LoginInfo  []*loginInfo `toml:"loginInfo"`
	sshMapInfo map[string]*sshInfo
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

func SSHConfigByName(sshName string) (utils.SSHConfig, error) {
	item, ok := LoginConf.sshMapInfo[sshName]
	if !ok {
		return nil, sshConfigNotExist(sshName)
	}
	sshConfig, err := item.ToSSHConfig()
	if err != nil {
		return nil, err
	}
	if !item.HaveProxy() {
		return sshConfig, nil
	}
	proxyConfig, err := LoginConf.sshMapInfo[item.ProxySSHName].ToSSHConfig()
	if err != nil {
		return nil, err
	}
	sshConfig.SetProxy(proxyConfig)
	return sshConfig, nil
}

type sshInfo struct {
	Name          string        `toml:"name"`
	AuthMethod    string        `toml:"authMethod"`
	URI           string        `toml:"uri"`
	User          string        `toml:"user"`
	PublicKeyPath string        `toml:"publicKeyPath"`
	Passphrase    string        `toml:"passphrase"`
	Timout        time.Duration `toml:"timout"`
	ProxySSHName  string        `toml:"proxySSHName"`
	Desc          string        `toml:"desc"`
	StartCommand  string        `toml:"startCommand"`
}

type loginInfo struct {
	URL          string `toml:"url"`
	ClientID     string `toml:"clientID"`
	ClientSecret string `toml:"clientSecret"`
	Name         string `toml:"name"`
	Pwd          string `toml:"pwd"`
	PwdUseMin    bool   `toml:"pwdUseMin"`
}

func (s *sshInfo) ToSSHConfig() (utils.SSHConfig, error) {
	authMethod := utils.AuthMethod(s.AuthMethod)
	if authMethod == utils.PublicKey {
		sshConf := &utils.SSHPrivateKeyConfig{SSHBaseConfig: &utils.SSHBaseConfig{
			MethodName:   authMethod,
			URI:          s.URI,
			User:         s.User,
			AuthMethods:  nil,
			Timout:       s.Timout,
			Passphrase:   s.Passphrase,
			StartCommand: s.StartCommand,
		},
		}
		pemBytes, err := ioutil.ReadFile(s.PublicKeyPath)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		var signer ssh.Signer
		if len(s.Passphrase) > 0 {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(pemBytes, []byte(s.Passphrase))
		} else {
			signer, err = ssh.ParsePrivateKey(pemBytes)
		}
		if err != nil {
			log.Printf("parse key failed:%v", err)
			return nil, err
		}
		sshConf.AuthMethods = []ssh.AuthMethod{ssh.PublicKeys(signer)}
		return sshConf, nil
	}
	return nil, errors.New("ToSSHConfig error invalid authMethod")
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
		LoginConf = &appConfig{sshMapInfo: map[string]*sshInfo{}}
	}
	if _, err = toml.DecodeReader(f, LoginConf); err != nil {
		return err
	}
	// verify
	var proxyName []string
	for _, item := range LoginConf.SSHInfo {
		if _, ok := LoginConf.sshMapInfo[item.Name]; ok {
			continue
		}
		LoginConf.sshMapInfo[item.Name] = item
		if item.HaveProxy() {
			proxyName = append(proxyName, item.Name)
		}
	}
	for _, item := range proxyName {
		if _, ok := LoginConf.sshMapInfo[item]; !ok {
			return sshConfigNotExist(item)
		}
	}
	return nil
}

func sshConfigNotExist(name string) error {
	return errors.New(fmt.Sprintf("ssh [%s] info not in config", name))
}
