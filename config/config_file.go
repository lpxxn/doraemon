package config

import (
	"os"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/skratchdot/open-golang/open"
)

const (
	confDirName = ".doraemon"
	//	confFileName = "doraemon.toml"
)

var confFileName = "doraemon.toml"

func GetConfig() (*os.File, error) {
	confPath := ConfFilePath()
	if _, err := os.Stat(confPath); os.IsNotExist(err) {
		return os.Create(confPath)
	}
	//return os.OpenFile(confPath, os.O_RDWR, 0666)
	return os.Open(confPath)
}

func ConfFilePath() string {
	return path.Join(ConfDir(), confFileName)
}
func ConfDir() string {
	rootDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	dirPath := path.Join(rootDir, confDirName)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err = os.MkdirAll(dirPath, os.ModePerm); err != nil {
			panic(err)
		}
	}
	return dirPath
}
func OpenConfDir() error {
	return open.Run(ConfDir())
}

func WritToConfig(v interface{}) error {
	confPath := ConfFilePath()
	var f *os.File
	if _, err := os.Stat(confPath); os.IsNotExist(err) {
		if f, err = os.Create(confPath); err != nil {
			return err
		}
	} else if f, err = os.OpenFile(confPath, os.O_RDWR, 0666); err != nil {
		return err
	}
	return toml.NewEncoder(f).Encode(v)
}
