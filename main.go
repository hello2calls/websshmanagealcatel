package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"time"
	"strconv"
)

type Connection struct {
	User		string
	Host		string
	Password	string
}

type SessionID struct {
	ID	int64
}

type Command struct {
	SessionID int64
	Command	string
}

// Declare Maps
var clientList = make(map[int64]*ssh.Client)
var sessionList =  make(map[int64]*ssh.Session)

func main() {

	fmt.Println("WebServer Listen on all interfaces, port 8080")

	http.HandleFunc("/", callFunc)
	http.ListenAndServe(":8080", nil)
	
}


func callFunc(w http.ResponseWriter, r *http.Request) {

	// If GET / (Website)
	if r.Method == "GET" && r.URL.Path == "/" {
		fmt.Println("GET /")
		w.Write([]byte("GET /"))
	}

	// If POST /API/session (API)
	if r.Method == "POST" && r.URL.Path == "/API/session" {
		body, _ := ioutil.ReadAll(r.Body)
		var c Connection
		json.Unmarshal(body, &c)
		client, session := connectToHost(c.User, c.Host, c.Password)
		timeNow := time.Now().UnixNano()
		clientList[timeNow] = client
		sessionList[timeNow] = session
		timeString := strconv.FormatInt(timeNow, 10)
		w.Write([]byte("{ID : " + timeString + "}"))
	}

	// If DELETE /API/session (API)
	if r.Method == "DELETE" && r.URL.Path == "/API/session" {
		body, _ := ioutil.ReadAll(r.Body)
		var s SessionID
		json.Unmarshal(body, &s)
		//closeSession
		closeSession(clientList[s.ID])
		delete(clientList, s.ID)
		w.Write([]byte("{sessionRemoved : true}"))
	}

	// If POST /API/command (API)
	if r.Method == "POST" && r.URL.Path == "/API/command" {
		body, _ := ioutil.ReadAll(r.Body)
		var c Command
		json.Unmarshal(body, &c)
		session := sessionList[c.SessionID]
		out := sendCommand(session, c.Command)
		w.Write([]byte("{commandOut : " + out + "}"))
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
