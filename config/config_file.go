package config

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/skratchdot/open-golang/open"
)

const (
	confDirName = ".doraemon"
)

var confFileName = "doraemon.toml"

//go:embed config_sample
var dummyConfData []byte

func GetConfig() (*os.File, error) {
	confPath := ConfFilePath()
	if _, err := os.Stat(confPath); os.IsNotExist(err) {
		// if not exit create and open the config dir
		defer func() {
			err := OpenConfDir()
			if err != nil {
				fmt.Println("open conf dir error")
			}
		}()
		return WritStringToConfig(string(dummyConfData))
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

func WritTomlToConfig(v interface{}) error {
	w, err := writeDataToConfig()
	if err != nil {
		return err
	}
	return toml.NewEncoder(w).Encode(v)
}

func WritStringToConfig(d string) (*os.File, error) {
	w, err := writeDataToConfig()
	if err != nil {
		return nil, err
	}
	_, err = io.WriteString(w, d)
	return w, nil
}

func writeDataToConfig() (*os.File, error) {
	confPath := ConfFilePath()
	var f *os.File
	if _, err := os.Stat(confPath); os.IsNotExist(err) {
		if f, err = os.Create(confPath); err != nil {
			return nil, err
		}
	} else if f, err = os.OpenFile(confPath, os.O_RDWR, 0666); err != nil {
		return nil, err
	}
	return f, nil
}
