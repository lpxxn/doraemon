package config

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/lpxxn/doraemon/internal"
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
	if err := internal.MakeFolder(dirPath, os.ModePerm); err != nil {
		panic(err)
	}
	return dirPath
}
func OpenConfDir() error {
	return internal.OpenFolder(ConfDir())
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
	if _, err := os.Stat(confPath); os.IsNotExist(err) {
		return os.Create(confPath)
	}
	return os.OpenFile(confPath, os.O_RDWR, 0666)
}
