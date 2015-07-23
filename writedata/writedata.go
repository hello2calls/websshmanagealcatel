package writedata

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

func writeStatus(dataFile Data, id string) {
	var oldDSLAM DSLAM
	var newDSLAM DSLAM

	dslamPos := getDslamPosByID(dataFile, id)
	oldDSLAM = getDslamByID(dataFile, id)

	req, _ := http.NewRequest("POST", "http://127.0.0.1:8080/API/session", bytes.NewBufferString("{\"user\": \""+oldDSLAM.User+"\", \"host\": \""+oldDSLAM.Address+"\", \"password\": \""+oldDSLAM.Password+"\"}"))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var response SessionResponse
	var err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("JSON Unmarshal body in writeStatus error:", err)
	}

	newDSLAM = oldDSLAM
	newDSLAM.Status = response.Status

	dataFile.DSLAM[dslamPos] = newDSLAM
	jsonIndent, _ := json.MarshalIndent(dataFile, "", "\t")
	ioutil.WriteFile("data.json", jsonIndent, 0777)
}

func writeCard(dataFile Data, sessionID, dslamID string) {
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

	var response CommandOut
	var err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("JSON Unmarshal body in writeCard error:", err)
	}

	var xmlBB bytes.Buffer
	reXML := regexp.MustCompile("^[\\ \\<]+.*")
	for _, line := range response.CommandOut {
		xmlLine := reXML.FindString(line)
		xmlBB.WriteString(xmlLine)
	}

	xmlB := []byte(xmlBB.String())
	var sec ShowEquipmentSlot

	err = xml.Unmarshal(xmlB, &sec)
	if err != nil {
		fmt.Println("XML Unmarshal xmlB in writeCard error:", err)
	}

	card := make([]Card, len(sec.Card))

	for i := 0; i < len(sec.Card); i++ {
		card[i].Name = sec.Card[i].Info[0][0].Value
		card[i].Slot = sec.Card[i].ResID
		card[i].OperStatus = sec.Card[i].Info[1][0].Value
		card[i].ErrorStatus = sec.Card[i].Info[2][0].Value
		card[i].Availability = sec.Card[i].Info[3][0].Value
	}

	dslam := getDslamByID(dataFile, dslamID)
	dslam.Card = card

	dslamPos := getDslamPosByID(dataFile, dslamID)
	dataFile.DSLAM[dslamPos] = dslam
	jsonIndent, _ := json.MarshalIndent(dataFile, "", "\t")
	ioutil.WriteFile("data.json", jsonIndent, 0777)

	writePort(dataFile, sessionID, dslamID)
}

func writePort(dataFile Data, sessionID, dslamID string) {
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

	var response CommandOut
	var err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("JSON Unmarshal body in writePort error:", err)
	}

	var xmlBB bytes.Buffer
	reXML := regexp.MustCompile("^[\\ \\<]+.*")
	for _, port := range response.CommandOut {
		xmlPort := reXML.FindString(port)
		xmlBB.WriteString(xmlPort)
	}

	xmlB := []byte(xmlBB.String())
	var sec ShowXdslOperationalDataLine

	err = xml.Unmarshal(xmlB, &sec)
	if err != nil {
		fmt.Println("XML Unmarshal xmlB in writeCard error:", err)
	}

	port := make([]Port, len(sec.Port))

	for i := 0; i < len(sec.Port); i++ {
		port[i].ID = sec.Port[i].ResID
		port[i].AdmState = sec.Port[i].Info[0][0].Value
		port[i].OprStateTxRateDs = sec.Port[i].Info[1][0].Value
		port[i].CurOpMode = sec.Port[i].Info[2][0].Value
	}

	dslam := getDslamByID(dataFile, dslamID)

	reCardID := regexp.MustCompile("1/1/[0-9]*")
	var cardID string

	for i := 0; i < len(dslam.Card); i++ {
		dslam.Card[i].Port = []Port{}
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

	dslamPos := getDslamPosByID(dataFile, dslamID)
	dataFile.DSLAM[dslamPos] = dslam
	jsonIndent, _ := json.MarshalIndent(dataFile, "", "\t")
	ioutil.WriteFile("data.json", jsonIndent, 0777)

	go writeService(dataFile, sessionID, dslamID)
}

func writeService(dataFile Data, sessionID, dslamID string) {
	fmt.Println("Update Service")

	dslam := getDslamByID(dataFile, dslamID)

	for i := 0; i < len(dslam.Card); i++ {
		for j := 0; j < len(dslam.Card[i].Port); j++ {
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

			var response CommandOut
			var err = json.Unmarshal(body, &response)
			if err != nil {
				fmt.Println("JSON Unmarshal body in writeService error:", err)
			}

			var xmlBB bytes.Buffer
			reXML := regexp.MustCompile("^[\\ \\<]+.*")
			for _, service := range response.CommandOut {
				xmlService := reXML.FindString(service)
				xmlBB.WriteString(xmlService)
			}

			xmlB := []byte(xmlBB.String())
			var sec ShowXdslOperDataPortIndexBridgePort

			err = xml.Unmarshal(xmlB, &sec)
			if err != nil {
				fmt.Println("XML Unmarshal xmlB in writeService error:", err)
			}

			service := make([]Service, len(sec.Service))

			for k := 0; k < len(sec.Service); k++ {
				service[k].ID = sec.Service[k].ResID
				service[k].Vlan = sec.Service[k].Parameter[0][0].Value
			}

			for l := 0; l < len(service); l++ {
				dslam.Card[i].Port[j].Service = append(dslam.Card[i].Port[j].Service, service[l])
			}

			dslamPos := getDslamPosByID(dataFile, dslamID)
			dataFile.DSLAM[dslamPos] = dslam
			jsonIndent, _ := json.MarshalIndent(dataFile, "", "\t")
			ioutil.WriteFile("data.json", jsonIndent, 0777)
		}
	}
}

func writeServiceOnePort(dataFile Data, sessionID, dslamID, index string) {
	dslam := getDslamByID(dataFile, dslamID)

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

	var response CommandOut
	var err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("JSON Unmarshal body in writeService error:", err)
	}

	var xmlBB bytes.Buffer
	reXML := regexp.MustCompile("^[\\ \\<]+.*")
	for _, service := range response.CommandOut {
		xmlService := reXML.FindString(service)
		xmlBB.WriteString(xmlService)
	}

	xmlB := []byte(xmlBB.String())
	var sec ShowXdslOperDataPortIndexBridgePort

	err = xml.Unmarshal(xmlB, &sec)
	if err != nil {
		fmt.Println("XML Unmarshal xmlB in writeService error:", err)
	}

	service := make([]Service, len(sec.Service))

	for k := 0; k < len(sec.Service); k++ {
		service[k].ID = sec.Service[k].ResID
		service[k].Vlan = sec.Service[k].Parameter[0][0].Value
	}

	for i := 0; i < len(dslam.Card); i++ {
		for j := 0; j < len(dslam.Card[i].Port); j++ {
			if dslam.Card[i].Port[j].ID == index {
				for l := 0; l < len(service); l++ {
					dslam.Card[i].Port[j].Service = append(dslam.Card[i].Port[j].Service, service[l])
				}
			}
		}
	}

	dslamPos := getDslamPosByID(dataFile, dslamID)
	dataFile.DSLAM[dslamPos] = dslam
	jsonIndent, _ := json.MarshalIndent(dataFile, "", "\t")
	ioutil.WriteFile("data.json", jsonIndent, 0777)

	fmt.Println("Data File Updated")
}
