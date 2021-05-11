package config

import (
	"time"

	"github.com/BurntSushi/toml"
)

type LoginConfig struct {
	SSHInfo   []*sshInfo   `toml:"sshInfo"`
	LoginInfo []*loginInfo `toml:"loginInfo"`
}

var LoginConf *LoginConfig

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

func ParseConfig() error {
	f, err := GetConfig()
	if err != nil {
		return err
	}
	_, err = toml.DecodeReader(f, LoginConf)
	return err
}
