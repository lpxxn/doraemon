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
	host := "5.8.1.4:22"
	user := "ec2-user"
	privateKey := ""
	//termlog := "./test_termlog"
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	if envHost := os.Getenv("SSH_HOST"); len(envHost) > 0 {
		host = envHost
	}
	if privateKey = os.Getenv("SSH_PRIVATE_KEY"); len(privateKey) == 0 {
		panic("SSH_PRIVATE_KEY")
	}

	pemBytes, err := ioutil.ReadFile(privateKey)
	if err != nil {
		log.Fatal(err)
	}
	signer, err := ssh.ParsePrivateKey(pemBytes)
	if err != nil {
		log.Fatalf("parse key failed:%v", err)
	}

	client, err := utils.CreateSSHClient(&utils.SSHPrivateKeyConfig{
		URI:         host,
		User:        user,
		AuthMethods: []ssh.AuthMethod{ssh.PublicKeys(signer)},
	}, utils.ProxyConfig(proxyConf()))
	if err != nil {
		panic(err)
	}
	// Set terminal log
	//client.SetLog(termlog, false)

	// Create Session
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

func proxyConf() *utils.SSHPrivateKeyConfig {
	rev := &utils.SSHPrivateKeyConfig{}
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
