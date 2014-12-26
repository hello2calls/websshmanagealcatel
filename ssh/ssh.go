package sshConnect

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"time"
	"strconv"
	"io"
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
var sessionIn =  make(map[int64]io.WriteCloser)
var sessionOut =  make(map[int64]io.Reader)
var sessionErr =  make(map[int64]io.Reader)


func SessionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

		case "GET":
		fmt.Println("GET /")
		w.Write([]byte("GET /"))
		case "POST":
		body, _ := ioutil.ReadAll(r.Body)
		var c Connection
		json.Unmarshal(body, &c)
		timeNow := time.Now().UnixNano()
		client, session := connectToHost(c.User, c.Host, c.Password, timeNow)
		clientList[timeNow] = client
		sessionList[timeNow] = session
		timeString := strconv.FormatInt(timeNow, 10)
		w.Write([]byte("{ID : " + timeString + "}"))
		case "PUT":
		w.Write([]byte("Not Implemented"))
		case "DELETE":
		body, _ := ioutil.ReadAll(r.Body)
		var s SessionID
		json.Unmarshal(body, &s)
		//closeSession
		closeSession(clientList[s.ID])
		delete(clientList, s.ID)
		w.Write([]byte("{sessionRemoved : true}"))
	}
}


func CommandHandler(w http.ResponseWriter, r *http.Request) {
	// If POST /API/command (API)
	if r.Method == "POST" && r.URL.Path == "/API/command" {
		body, _ := ioutil.ReadAll(r.Body)
		var c Command
		json.Unmarshal(body, &c)
		session := sessionList[c.SessionID]
		out := sendCommand(session, c.Command, c.SessionID)
		w.Write([]byte("{commandOut : " + out + "}"))
	}
}


// Connection to Host
func connectToHost(user, host, password string, timeNow int64) (*ssh.Client, *ssh.Session) {

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

	sessionOut[timeNow], _ = session.StdoutPipe()
	sessionIn[timeNow], _ = session.StdinPipe()
	sessionErr[timeNow], _ = session.StderrPipe()
	session.Shell()
	buf := make([]byte, 10000)
	sessionOut[timeNow].Read(buf)

	// Return client and session
	return client, session

}


// SendCommand to Host
func sendCommand(session *ssh.Session, command string, sessionID int64) (string) {

	buf := make([]byte, 10000)

	// Send command
	switch command {
		case "":
		return string("Command is Empty")
		case "\n":
		return string("Command is not defined")
		default :
		sessionIn[sessionID].Write([]byte(command + "\n"))
		n, _ := sessionOut[sessionID].Read(buf)
		loadStr := string(buf[:n])
		// Return command result
		return string(loadStr)
	}

}


// Close session
func closeSession(client *ssh.Client) {

	client.Close()

}
