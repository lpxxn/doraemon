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

type AppConfig struct {
	SSHInfo    sshInfoList `toml:"sshInfo"`
	CmdInfo    cmdInfoList `toml:"CmdInfo"`
	sshMapInfo map[string]*sshInfo
	cmdMapInfo map[string]*cmdInfo
}

var LoginConf *AppConfig

func (a *AppConfig) ConfigByName(name string) (*sshInfo, error) {
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
		return nil, configNotExist(sshName)
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

func CustomConfigByName(name string) (*cmdInfo, error) {
	item, ok := LoginConf.cmdMapInfo[name]
	if !ok {
		return nil, configNotExist(name)
	}
	return item, nil
}

type Info interface {
	GetName() string
	GetDesc() string
}

type InfoIterator interface {
	HasNext() bool
	Next() Info
}
type InfoCollection interface {
	GetIterator() InfoIterator
}

type sshInfoList []*sshInfo
type sshIterator struct {
	data  []*sshInfo
	index int
}

func (s sshIterator) HasNext() bool {
	return len(s.data) > s.index
}

func (s *sshIterator) Next() Info {
	if s.HasNext() {
		v := s.data[s.index]
		s.index++
		return v
	}
	return nil
}

func (s sshInfoList) GetIterator() InfoIterator {
	return &sshIterator{
		data:  s,
		index: 0,
	}
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

func (s *sshInfo) GetName() string {
	return s.Name
}

func (s *sshInfo) GetDesc() string {
	return s.Desc
}

type cmdInfo struct {
	Name string `toml:"name"`
	Cmd  string `toml:"cmd"`
	Desc string `toml:"desc"`
}

func (c cmdInfo) GetName() string {
	return c.Name
}

func (c cmdInfo) GetDesc() string {
	return c.Desc
}

type cmdInfoList []*cmdInfo

func (c cmdInfoList) GetIterator() InfoIterator {
	return &cmdIterator{
		data:  c,
		index: 0,
	}
}

type cmdIterator struct {
	data  cmdInfoList
	index int
}

func (c cmdIterator) HasNext() bool {
	return len(c.data) > c.index
}

func (c *cmdIterator) Next() Info {
	if c.HasNext() {
		v := c.data[c.index]
		c.index++
		return v
	}
	return nil
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
	} else if authMethod == utils.Password {
		sshConf := &utils.SSHPasswordConfig{SSHBaseConfig: &utils.SSHBaseConfig{
			MethodName:   authMethod,
			URI:          s.URI,
			User:         s.User,
			AuthMethods:  []ssh.AuthMethod{ssh.Password(s.Passphrase)},
			Timout:       s.Timout,
			Passphrase:   s.Passphrase,
			StartCommand: s.StartCommand,
		},
		}
		return sshConf, nil
	}
	return nil, errors.New("ToSSHConfig error invalid authMethod: " + string(authMethod))
}

func (s *sshInfo) HaveProxy() bool {
	return len(s.ProxySSHName) > 0
}

func ParseConfig() (*AppConfig, error) {
	f, err := GetConfig()
	if err != nil {
		return nil, err
	}
	if LoginConf == nil {
		LoginConf = &AppConfig{sshMapInfo: map[string]*sshInfo{}, cmdMapInfo: map[string]*cmdInfo{}}
	}
	if _, err = toml.DecodeReader(f, LoginConf); err != nil {
		return nil, err
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
			return nil, configNotExist(item)
		}
	}
	for _, item := range LoginConf.CmdInfo {
		LoginConf.cmdMapInfo[item.Name] = item
	}
	return LoginConf, nil
}

func configNotExist(name string) error {
	return errors.New(fmt.Sprintf("config [%s] info not in config", name))
}
