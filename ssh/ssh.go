package sshConnect

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"io"
	"code.google.com/p/go-uuid/uuid"
)

// Define Connection
type Connection struct {
	User		string
	Host		string
	Password	string
}

// Define SessionID
type SessionID struct {
	ID	string
}

// Define Command
type Command struct {
	SessionID string
	Command	string
}

// Declare Maps
var clientList = make(map[string]*ssh.Client)
var sessionList = make(map[string]*ssh.Session)
var sessionIn = make(map[string]io.WriteCloser)
var sessionOut = make(map[string]io.Reader)
var sessionErr = make(map[string]io.Reader)


func SessionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case "GET":
			fmt.Println("GET /")
			w.Write([]byte("GET /"))
		case "POST":
			body, _ := ioutil.ReadAll(r.Body)
			var c Connection
			json.Unmarshal(body, &c)
			uuid := uuid.New()
			client, session, errorSsh := connectToHost(c.User, c.Host, c.Password, uuid)
			if errorSsh != "OK" {
				w.Write([]byte("{status : " + errorSsh + "}"))
			} else {
				clientList[uuid] = client
				sessionList[uuid] = session
				w.Write([]byte("{ID : " + uuid + ", status : OK}"))
			}
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
func connectToHost(user, host, password, uuid string) (client *ssh.Client, session *ssh.Session, errorSsh string) {

	// Create sshConfig variable
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password(password)},
	}

	// Create Client
	client, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		errorSsh = "SSH_KO"
		// Return client, session and error
		return client, session, errorSsh
	} else {
		// Create Session
		session, err = client.NewSession()
		if err != nil {
			client.Close()
			errorSsh = "Session_KO"
			// Return client, session and error
			return client, session, errorSsh
		} else {
			// Save Session
			sessionOut[uuid], err = session.StdoutPipe()
			sessionIn[uuid], _ = session.StdinPipe()
			sessionErr[uuid], _ = session.StderrPipe()
			session.Shell()
			var buf []byte
			sessionOut[uuid].Read(buf)
			// Return client, session and error
			errorSsh = "OK"
			return client, session, errorSsh
		}
	}

}


// SendCommand to Host
func sendCommand(session *ssh.Session, command, sessionID string) (string) {

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
