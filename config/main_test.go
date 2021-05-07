package config

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	confFileName = "doraemon_test.toml"
	os.Exit(m.Run())
}
