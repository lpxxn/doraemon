package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/lpxxn/doraemon/utils"
	"golang.org/x/crypto/ssh"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	conf := mConf()
	client, err := utils.NewSSHClient(conf)
	if err != nil {
		panic(err)
	}
	session, err := client.CreateSession()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Start ssh shell
	if err := client.Shell(session); err != nil {
		panic(err)
	}
}

func mConf() *utils.SSHPrivateKeyConfig {
	rev := &utils.SSHPrivateKeyConfig{SSHBaseConfig: &utils.SSHBaseConfig{MethodName: utils.PublicKey, StartCommand: "ls -l; whoami"}}
	if envHost := os.Getenv("SSH_HOST"); len(envHost) > 0 {
		rev.URI = envHost
	}
	if envUser := os.Getenv("SSH_USER"); len(envUser) > 0 {
		rev.User = envUser
	}
	if privateKey := os.Getenv("SSH_PRIVATE_KEY"); len(privateKey) > 0 {
		pemBytes, err := ioutil.ReadFile(privateKey)
		if err != nil {
			log.Fatal(err)
		}
		signer, err := ssh.ParsePrivateKey(pemBytes)
		if err != nil {
			log.Fatalf("parse key failed:%v", err)
		}
		rev.AuthMethods = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	}
	rev.Proxy = proxyConf()
	return rev
}

func proxyConf() *utils.SSHPrivateKeyConfig {
	rev := &utils.SSHPrivateKeyConfig{SSHBaseConfig: &utils.SSHBaseConfig{MethodName: utils.PublicKey}}
	if envProxyHost := os.Getenv("SSH_PROXY_HOST"); len(envProxyHost) > 0 {
		rev.URI = envProxyHost
	}
	if envProxyUser := os.Getenv("SSH_PROXY_USER"); len(envProxyUser) > 0 {
		rev.User = envProxyUser
	}
	if privateKey := os.Getenv("SSH_PROXY_PRIVATE_KEY"); len(privateKey) > 0 {
		pemBytes, err := ioutil.ReadFile(privateKey)
		if err != nil {
			log.Fatal(err)
		}
		signer, err := ssh.ParsePrivateKey(pemBytes)
		if err != nil {
			log.Fatalf("parse key failed:%v", err)
		}
		rev.AuthMethods = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	}
	return rev
}
