package sshsession

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"bitbucket.org/nmontes/WebSSHManageAlcatel/web/pkg/equipment"
	"bitbucket.org/nmontes/WebSSHManageAlcatel/web/pkg/logger"
	S "bitbucket.org/nmontes/WebSSHManageAlcatel/web/pkg/structures"
)

// Get to connect to DSLAM
func Get(dataFile S.Data, id string) string {
	var dslam S.DSLAM
	dslam = equipment.GetDslamByID(dataFile, id)
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8080/API/session", bytes.NewBufferString("{\"user\": \""+dslam.User+"\", \"host\": \""+dslam.Address+"\", \"password\": \""+dslam.Password+"\"}"))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var response S.SessionResponse
	var err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("JSON Unmarshal body in getSSHSession error:", err)
		logger.Print("JSON Unmarshal body in getSSHSession error:", err)
	}

	var status string
	if response.Status != "OK" {
		status = "SSH_KO"
	} else {
		status = response.ID
	}
	return status
}

// Delete to unconnect to DSLAM
func Delete(id string) string {
	req, _ := http.NewRequest("DELETE", "http://127.0.0.1:8080/API/session", bytes.NewBufferString("{\"ID\": \""+id+"\"}"))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var response S.SessionResponse
	var err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("JSON Unmarshal body in getSSHSession error:", err)
		logger.Print("JSON Unmarshal body in getSSHSession error:", err)
	}

	return "Removed"
}
