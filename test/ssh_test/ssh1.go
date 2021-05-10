package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/lpxxn/doraemon/ssh_utils"
	"golang.org/x/crypto/ssh"
)

func main() {
	host := "5.8.1.4:22"
	user := "ec2-user"
	privateKey := "/Users/li/.ssh/my_test.pem"
	//termlog := "./test_termlog"
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	if envHost := os.Getenv("SSH_HOST"); len(envHost) > 0 {
		host = envHost
	}

	pemBytes, err := ioutil.ReadFile(privateKey)
	if err != nil {
		log.Fatal(err)
	}
	signer, err := ssh.ParsePrivateKey(pemBytes)
	if err != nil {
		log.Fatalf("parse key failed:%v", err)
	}

	client, err := ssh_utils.CreateSSHClient(host, user, []ssh.AuthMethod{ssh.PublicKeys(signer)})
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
