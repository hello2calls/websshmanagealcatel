package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	)

	func main() {

		client, session, err := connectToHost("test", "192.168.200.10:22", "test")
		if err != nil {
			panic(err)
		}

		out, err := sendCommand(session, "ls")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(out))

		client.Close()
	}


	func connectToHost(user, host, password string) (*ssh.Client, *ssh.Session, error) {

		sshConfig := &ssh.ClientConfig{
			User: user,
			Auth: []ssh.AuthMethod{ssh.Password(password)},
		}

		client, err := ssh.Dial("tcp", host, sshConfig)
		if err != nil {
			return nil, nil, err
		}

		session, err := client.NewSession()
		if err != nil {
			client.Close()
			return nil, nil, err
		}

		return client, session, nil
	}

	func sendCommand(session *ssh.Session, command string) (string, error) {

		out, err := session.CombinedOutput(command)
		if err != nil {
			return "", err
		}

		return string(out), nil
	}
