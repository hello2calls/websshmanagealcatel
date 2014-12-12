package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	)

	func main() {

		client, session, err := connectToHost()
		if err != nil {
			panic(err)
		}
		out, err := session.CombinedOutput("ls")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(out))
		client.Close()
	}


	func connectToHost() (*ssh.Client, *ssh.Session, error) {

		sshConfig := &ssh.ClientConfig{
			User: "test",
			Auth: []ssh.AuthMethod{ssh.Password("test")},
		}

		client, err := ssh.Dial("tcp", "192.168.200.10:22", sshConfig)
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
