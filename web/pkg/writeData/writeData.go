package writeData

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"bitbucket.org/nmontes/WebSSHManageAlcatel/web/pkg/equipment"
	"bitbucket.org/nmontes/WebSSHManageAlcatel/web/pkg/file"
	S "bitbucket.org/nmontes/WebSSHManageAlcatel/web/pkg/structures"
)

// WriteStatus write card status
func WriteStatus(dataFile S.Data, id string) {

	fmt.Println("Write DSLAM Status")

	var oldDSLAM S.DSLAM
	var newDSLAM S.DSLAM

	dslamPos := equipment.GetDslamPosByID(dataFile, id)
	oldDSLAM = equipment.GetDslamByID(dataFile, id)

	req, _ := http.NewRequest("POST", "http://127.0.0.1:8080/API/session", bytes.NewBufferString("{\"user\": \""+oldDSLAM.User+"\", \"host\": \""+oldDSLAM.Address+"\", \"password\": \""+oldDSLAM.Password+"\"}"))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var response S.SessionResponse
	var err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("JSON Unmarshal body in writeStatus error:", err)
	}

	newDSLAM = oldDSLAM
	newDSLAM.Status = response.Status

	dataFile.DSLAM[dslamPos] = newDSLAM

	file.WriteFile(dataFile, "data.json")

	reqBis, _ := http.NewRequest("DELETE", "http://127.0.0.1:8080/API/session", bytes.NewBufferString("{\"ID\": \""+response.ID+"\"}"))
	reqBis.Header.Set("Content-Type", "application/json")
	clientBis := &http.Client{}
	respBis, _ := clientBis.Do(reqBis)
	defer respBis.Body.Close()

}

// WriteCard write all DSLAM Card
func WriteCard(dataFile S.Data, sessionID, dslamID string) {

	fmt.Println("Update Card")

	req, _ := http.NewRequest("POST", "http://127.0.0.1:8080/API/command", bytes.NewBufferString("{\"sessionID\": \""+sessionID+"\", \"command\": \"show equipment slot xml\"}"))
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
	bodyS = strings.Replace(bodyS, "\x00", "", -1)
	body = []byte(bodyS)

	var response S.CommandOut
	var err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("JSON Unmarshal body in writeCard error:", err)
		fmt.Println(bodyS)
	}

	var xmlBB bytes.Buffer
	reXML := regexp.MustCompile("^[\\ \\<]+.*")
	for _, line := range response.CommandOut {
		xmlLine := reXML.FindString(line)
		xmlBB.WriteString(xmlLine)
	}

	xmlB := []byte(xmlBB.String())
	var sec S.ShowEquipmentSlot

	err = xml.Unmarshal(xmlB, &sec)
	if err != nil {
		fmt.Println("XML Unmarshal xmlB in writeCard error:", err)
		fmt.Println(bodyS)
	}

	card := make([]S.Card, len(sec.Card))

	for i := 0; i < len(sec.Card); i++ {
		card[i].Name = sec.Card[i].Info[0][0].Value
		card[i].Slot = sec.Card[i].ResID
		card[i].OperStatus = sec.Card[i].Info[1][0].Value
		card[i].ErrorStatus = sec.Card[i].Info[2][0].Value
		card[i].Availability = sec.Card[i].Info[3][0].Value
	}

	dslam := equipment.GetDslamByID(dataFile, dslamID)
	dslam.Card = card

	dslamPos := equipment.GetDslamPosByID(dataFile, dslamID)
	dataFile.DSLAM[dslamPos] = dslam
	file.WriteFile(dataFile, "data.json")

	WritePort(dataFile, sessionID, dslamID)
}

// WritePort write all DSLAM port
func WritePort(dataFile S.Data, sessionID, dslamID string) {

	fmt.Println("Update Port")

	req, _ := http.NewRequest("POST", "http://127.0.0.1:8080/API/command", bytes.NewBufferString("{\"sessionID\": \""+sessionID+"\", \"command\": \"show xdsl operational-data line xml\"}"))
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
		fmt.Println("JSON Unmarshal body in writePort error:", err)
		fmt.Println(bodyS)
	}

	var xmlBB bytes.Buffer
	reXML := regexp.MustCompile("^[\\ \\<]+.*")
	for _, port := range response.CommandOut {
		xmlPort := reXML.FindString(port)
		xmlBB.WriteString(xmlPort)
	}

	xmlB := []byte(xmlBB.String())
	var sec S.ShowXdslOperationalDataLine

	err = xml.Unmarshal(xmlB, &sec)
	if err != nil {
		fmt.Println("XML Unmarshal xmlB in writePort error:", err)
		fmt.Println(bodyS)
	}

	port := make([]S.Port, len(sec.Port))

	for i := 0; i < len(sec.Port); i++ {
		port[i].ID = sec.Port[i].ResID
		port[i].AdmState = sec.Port[i].Info[0][0].Value
		port[i].OprStateTxRateDs = sec.Port[i].Info[1][0].Value
		port[i].CurOpMode = sec.Port[i].Info[2][0].Value
	}

	dslam := equipment.GetDslamByID(dataFile, dslamID)

	reCardID := regexp.MustCompile("1/1/[0-9]*")
	var cardID string

	for i := 0; i < len(dslam.Card); i++ {
		dslam.Card[i].Port = []S.Port{}
		cardID = reCardID.FindString(dslam.Card[i].Slot)
		reCardIDPort := regexp.MustCompile("^1/1/[0-9]*")
		var cardIDPort string
		for j := 0; j < len(port); j++ {
			cardIDPort = reCardIDPort.FindString(port[j].ID)
			if cardID == cardIDPort {
				dslam.Card[i].Port = append(dslam.Card[i].Port, port[j])
			}
		}
	}

	dslamPos := equipment.GetDslamPosByID(dataFile, dslamID)
	dataFile.DSLAM[dslamPos] = dslam
	file.WriteFile(dataFile, "data.json")

	WriteAllServices(dataFile, sessionID, dslamID)
}

// WriteAllServices write all port services
func WriteAllServices(dataFile S.Data, sessionID, dslamID string) {
	fmt.Println("Update Service")

	dslam := equipment.GetDslamByID(dataFile, dslamID)

	for i := 0; i < len(dslam.Card); i++ {
		if dslam.Card[i].ErrorStatus == "no-error" {
			for j := 0; j < len(dslam.Card[i].Port); j++ {
				if dslam.Card[i].Port[j].AdmState == "up" {
					index := dslam.Card[i].Port[j].ID

					req, _ := http.NewRequest("POST", "http://127.0.0.1:8080/API/command", bytes.NewBufferString("{\"sessionID\": \""+sessionID+"\", \"command\": \"show xdsl oper-data-port "+index+" bridge-port xml\"}"))
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
						fmt.Println("JSON Unmarshal body in writeService error:", err)
						fmt.Println(bodyS)
					}

					var xmlBB bytes.Buffer
					reXML := regexp.MustCompile("^[\\ \\<]+.*")
					for _, service := range response.CommandOut {
						xmlService := reXML.FindString(service)
						xmlBB.WriteString(xmlService)
					}

					xmlB := []byte(xmlBB.String())
					var sec S.ShowXdslOperDataPortIndexBridgePort

					err = xml.Unmarshal(xmlB, &sec)
					if err != nil {
						fmt.Println("XML Unmarshal xmlB in writeService error:", err)
						fmt.Println(bodyS)
					}

					service := make([]S.Service, len(sec.Service))

					for k := 0; k < len(sec.Service); k++ {
						service[k].ID = sec.Service[k].ResID
						service[k].Vlan = sec.Service[k].Parameter[0][0].Value
					}

					for l := 0; l < len(service); l++ {
						dslam.Card[i].Port[j].Service = append(dslam.Card[i].Port[j].Service, service[l])
					}

					dslamPos := equipment.GetDslamPosByID(dataFile, dslamID)
					dataFile.DSLAM[dslamPos] = dslam
					file.WriteFile(dataFile, "data.json")
				}
			}
		}
	}
}

// WriteServiceOnePort Write Service of One Port in data file when we update service
func WriteServiceOnePort(dataFile S.Data, sessionID, dslamID, index string) {
	dslam := equipment.GetDslamByID(dataFile, dslamID)

	req, _ := http.NewRequest("POST", "http://127.0.0.1:8080/API/command", bytes.NewBufferString("{\"sessionID\": \""+sessionID+"\", \"command\": \"show xdsl oper-data-port "+index+" bridge-port xml\"}"))
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
		fmt.Println("JSON Unmarshal body in writeService error:", err)
		fmt.Println(bodyS)
	}

	var xmlBB bytes.Buffer
	reXML := regexp.MustCompile("^[\\ \\<]+.*")
	for _, service := range response.CommandOut {
		xmlService := reXML.FindString(service)
		xmlBB.WriteString(xmlService)
	}

	xmlB := []byte(xmlBB.String())
	var sec S.ShowXdslOperDataPortIndexBridgePort

	err = xml.Unmarshal(xmlB, &sec)
	if err != nil {
<<<<<<< HEAD
		fmt.Println("XML Unmarshal xmlB in writeServiceOnePort error:", err)
=======
		fmt.Println("XML Unmarshal xmlB in writeService error:", err)
>>>>>>> 4a63e4a2ae8545624741a7207a685e797ff93859
		fmt.Println(bodyS)
	}

	service := make([]S.Service, len(sec.Service))

	for k := 0; k < len(sec.Service); k++ {
		service[k].ID = sec.Service[k].ResID
		service[k].Vlan = sec.Service[k].Parameter[0][0].Value
	}

	for i := 0; i < len(dslam.Card); i++ {
		for j := 0; j < len(dslam.Card[i].Port); j++ {
			if dslam.Card[i].Port[j].ID == index {
				dslam.Card[i].Port[j].Service = service
			}
		}
	}

	dslamPos := equipment.GetDslamPosByID(dataFile, dslamID)
	dataFile.DSLAM[dslamPos] = dslam
	file.WriteFile(dataFile, "data.json")

	fmt.Println("Data File Updated")
}
