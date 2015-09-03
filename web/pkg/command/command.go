package command

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	S "bitbucket.org/nmontes/WebSSHManageAlcatel/web/pkg/structures"
)

// GetOut return command result
func GetOut(sessionID, command string) []string {
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8080/API/command", bytes.NewBufferString("{\"sessionID\": \""+sessionID+"\", \"command\": \""+command+"\"}"))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	bodyS := string(body)
	bodyS = strings.Replace(bodyS, "\a", "", -1)
	bodyS = strings.Replace(bodyS, "\x1b", "", -1)
	bodyS = strings.Replace(bodyS, "-[1D", "", -1)
	bodyS = strings.Replace(bodyS, "[1D", "", -1)
	bodyS = strings.Replace(bodyS, "\\ ", " ", -1)
	body = []byte(bodyS)

	var response S.CommandOut
	var err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("JSON Unmarshal body in getCommandOut error:", err)
	}

	return response.CommandOut
}

// Set send command
func Set(sessionID, command string) {
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8080/API/command", bytes.NewBufferString("{\"sessionID\": \""+sessionID+"\", \"command\": \""+command+"\"}"))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
}
