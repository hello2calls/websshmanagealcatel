package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

type Connection struct {
	User		string
	Host		string
	Password	string
}

// Declare Slices
var client []*ssh.CLient
var session []*ssh.Session

func main() {

	http.HandleFunc("/", callFunc)
	http.ListenAndServe(":8080", nil)

	//client, session := connectToHost("test", "192.168.200.10:22", "test")

	//out := sendCommand(session, "lsqdf")
	//fmt.Println(string(out))

	//closeSession(client)

}


func callFunc(w http.ResponseWriter, r *http.Request) {

	// If GET / (Website)
	if r.Method == "GET" && r.URL.Path == "/" {
		fmt.Println("GET /")
		w.Write([]byte("GET /"))
	}

	// If POST /API/session (API)
	if r.Method == "POST" && r.URL.Path == "/API/session" {
		fmt.Println("POST /API/session")
		body, _ := ioutil.ReadAll(r.Body)
		var c Connection
		err := json.Unmarshal(body, &c)
		if err != nil {
			panic(err)
		}
		client, session := connectToHost(c.User, c.Host, c.Password)
		w.Write([]byte("client"))
	}

	// If DELETE /API/session (API)
	if r.Method == "DELETE" && r.URL.Path == "/API/session" {
		fmt.Println("POST /API/session")
		body, _ := ioutil.ReadAll(r.Body)

		//closeSession(body)
		w.Write([]byte("DELETE /API/session"))
	}

	// If POST /API/command (API)
	if r.Method == "POST" && r.URL.Path == "/API/command" {
		fmt.Println("POST /API/command")
		w.Write([]byte("POST /API/command"))
	}

}


// Connection to Host
func connectToHost(user, host, password string) (*ssh.Client, *ssh.Session) {

	// Create sshConfig variable
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password(password)},
	}

	// Create Client
	client, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		panic(err)
	}

	// Create Session
	session, err := client.NewSession()
	if err != nil {
		client.Close()
		panic(err)
	}

	// Return client and session
	return client, session

}

// SendCommand to Host
func sendCommand(session *ssh.Session, command string) (string) {

	// Send command
	out, err := session.CombinedOutput(command)
	if err != nil {}

	// Return command result
	return string(out)

}


// Close session
func closeSession(client *ssh.Client) {

	client.Close()

}
