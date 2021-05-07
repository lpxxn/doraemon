package config

import (
	"io"
	"testing"
)

func TestGetConfig(t *testing.T) {
	f, err := GetConfig()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(io.ReadAll(f))
}

func TestWriteConf(t *testing.T) {
	conf := &LoginConfig{LoginInfo: []*loginInfo{
		&loginInfo{
			URL:          "urltest",
			ClientID:     "ac",
			ClientSecret: "dafd",
			Name:         "haha",
			Pwd:          "asdfasdf",
			PwdUseMin:    false,
		},
	}}
	err := WritToConfig(conf)
	if err != nil {
		t.Fatal(err)
	}
}

func TestOpenConfDir(t *testing.T) {
	if err := OpenConfDir(); err != nil {
		t.Fatal(err)
	}
}
