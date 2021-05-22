package config

import (
	"io"
	"testing"

	"github.com/BurntSushi/toml"
)

func TestGetConfig(t *testing.T) {
	f, err := GetConfig()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(io.ReadAll(f))
}

func TestWriteConf(t *testing.T) {
	conf := &AppConfig{LoginInfo: []*loginInfo{
		&loginInfo{
			URL:          "urltest",
			ClientID:     "ac",
			ClientSecret: "dafd",
			Name:         "haha",
			Pwd:          "asdfasdf",
			PwdUseMin:    false,
		},
	}, SSHInfo: []*sshInfo{&sshInfo{
		Name:          "test1",
		AuthMethod:    "publickey",
		URI:           "127.0.0.1:22",
		User:          "li",
		PublicKeyPath: "a/b/c",
		Timout:        0,
		ProxySSHName:  "proxy",
		Desc:          "test client",
	}, &sshInfo{
		Name:          "proxy",
		AuthMethod:    "publickey",
		URI:           "127.0.0.1:22",
		User:          "li",
		PublicKeyPath: "a/b/c",
		Timout:        0,
		ProxySSHName:  "",
		Desc:          "proxy",
	}}}
	err := WritToConfig(conf)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDecode(t *testing.T) {
	src := `[[sshInfo]]
  authMethod = "publickey"
  uri = "127.0.0.1:22"
  user = "li"
  publicKeyPath = "a/b/c"
  timout = 0

[[loginInfo]]
  url = "urltest"
  clientID = "ac"
  clientSecret = "dafd"
  pwd = """asdfas"df"
   dddd
"""
  pwdUseMin = false
  name = """h \
a"ha hi"""
`
	conf := &AppConfig{}
	if _, err := toml.Decode(src, conf); err != nil {
		t.Fatal(err)
	}
	t.Log(*conf.LoginInfo[0])
	/*
	有 \ 就没有\n
	  pwd = """asdfas"df" \
	   dddd
	"""
	// asdfas"df" dddd

	没有 \ 就会有\n 换行
	asdfas"df"
	   dddd

	*/
}

func TestOpenConfDir(t *testing.T) {
	if err := OpenConfDir(); err != nil {
		t.Fatal(err)
	}
}
