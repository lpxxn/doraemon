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
