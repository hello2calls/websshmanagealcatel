package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
)

func main() {

	client, session := connectToHost("test", "192.168.200.10:22", "test")

	out := sendCommand(session, "lsqdf")
	fmt.Println(string(out))

	closeSession(client)

}


func connectToHost(user, host, password string) (*ssh.Client, *ssh.Session) {

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password(password)},
	}

	client, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		panic(err)
	}

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		panic(err)
	}

	return client, session

}


func sendCommand(session *ssh.Session, command string) (string) {

	out, err := session.CombinedOutput(command)
	if err != nil {}

	return string(out)

}

func closeSession(client *ssh.Client) {

	client.Close()

}
