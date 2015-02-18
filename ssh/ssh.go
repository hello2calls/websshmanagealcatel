package sshConnect

import (
	"fmt"
	//"golang.org/x/crypto/ssh"
	"code.google.com/marksheahan-sshblock/ssh"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"io"
	"code.google.com/p/go-uuid/uuid"
	"time"
	"bytes"
	"encoding/binary"
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
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte("{\"Status\" : \"" + errorSsh + "\"}"))
			} else {
				clientList[uuid] = client
				sessionList[uuid] = session
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte("{\"ID\" : \"" + uuid + "\", \"Status\" : \"OK\"}"))
			}
		case "PUT":
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("Not Implemented"))
		case "DELETE":
			body, _ := ioutil.ReadAll(r.Body)
			var s SessionID
			json.Unmarshal(body, &s)
			//closeSession
			closeSession(clientList[s.ID])
			delete(clientList, s.ID)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("{\"sessionRemoved\" : \"true\"}"))
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

		str := make([]string, len(bytes.Split([]byte(out), []byte{'\n'})))
		for i, line := range bytes.Split([]byte(out), []byte{'\n'}) {
			fmt.Println("line -"+string(line)+"-")
			fmt.Println(line)
			fmt.Println(line[len(line)-1])
			//fmt.Println(bytes.Contains([]byte(strconv.Itoa(13)), line))
			bs := make([]byte, 4)
			binary.LittleEndian.PutUint32(bs, 13)
			fmt.Println(bs[0])
			if bs[0] == line[len(line)-1] {
				if len(line)-1 != 0 {
					str[i] = string(line[:len(line)-1])
				}
				fmt.Println("contain 13")
			} else {
				str[i] = string(line)
			}
		}

		var commandOut string
		commandOut = "["
		for i, line := range str {
			fmt.Println(commandOut)
			commandOut += "'" + line + "'"
			if i != len(str)-1 {
				commandOut += ","
			}
		}
		commandOut += "]"
		fmt.Println("COMMAND : ")
		fmt.Println(commandOut)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{'CommandOut' : " + commandOut + "}"))
	}
}


// Connection to Host
func connectToHost(user, host, password, uuid string) (client *ssh.Client, session *ssh.Session, errorSsh string) {

	// Create sshConfig variable
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password(password)},
		Config: ssh.Config{
			Ciphers: ssh.AllSupportedCiphers(),
			},
	}

	// Create Client
	client, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		errorSsh = err.Error()
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

	buf := make([]byte, 1000)

	var loadStr bytes.Buffer

	// Send command
	switch command {
		case "":
			loadStr.WriteString("Command is Empty")
			return loadStr.String()
		case "\n":
			loadStr.WriteString("Command is not defined")
			return loadStr.String()
		default :
			sessionIn[sessionID].Write([]byte(command + "\n"))
			time.Sleep(1000 * time.Millisecond)

			n, _ := sessionOut[sessionID].Read(buf);

			for n > 0 {
				if n < 1000 {
					break
				}
				loadStr.WriteString(string(buf[:n]))
				time.Sleep(1000 * time.Millisecond)
				n, _ = sessionOut[sessionID].Read(buf)
			}

			loadStr.WriteString(string(buf[:n]))

			fmt.Println(loadStr.String())

			return loadStr.String()
	}

}


// Close session
func closeSession(client *ssh.Client) {

	client.Close()

}
