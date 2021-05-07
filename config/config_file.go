package config

import (
	"os"
	"path"
)

const (
	confDirName  = ".doraemon"
	confFileName = "doraemon.toml"
)

func GetConfig() (*os.File, error) {
	rootDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	dirPath := path.Join(rootDir, confDirName)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err = os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return nil, err
		}
	}
	confPath := path.Join(dirPath, confFileName)
	if _, err := os.Stat(confPath); os.IsNotExist(err) {
		return os.Create(confPath)
	}
	//return os.OpenFile(confPath, os.O_RDWR, 0666)
	return os.Open(confPath)
}
