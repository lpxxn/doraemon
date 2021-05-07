package config

import (
	"github.com/BurntSushi/toml"
)

type LoginConfig struct {
	LoginInfo []*loginInfo `toml:"loginInfo"`
}

var LoginConf *LoginConfig

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
