package web

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"text/template"

	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"

	"code.google.com/p/go-uuid/uuid"
	"github.com/GeertJohan/go.rice"
)

// ---------------
// JSON Structures
// ---------------

// Service set triple play services
type Service struct {
	ID   string `json:"Id"`
	Vlan string `json:"Vlan"`
}

// Port define DSLAM ports
type Port struct {
	ID               string    `json:"Index"`
	Name             string    `json:"Name"`
	AdmState         string    `json:"Adm-State"`
	OprStateTxRateDs string    `json:"Opr-State/Tx-Rate-Ds"`
	CurOpMode        string    `json:"Cur-Op-Mode"`
	Service          []Service `json:"Service"`
}

// Card define DSLAM card
type Card struct {
	Name         string `json:"Name"`
	Slot         string `json:"Slot"`
	OperStatus   string `json:"Opers-Status"`
	ErrorStatus  string `json:"Error-Status"`
	Availability string `json:"Availability"`
	Port         []Port `json:"Port"`
}

// DSLAM define DSLAM
type DSLAM struct {
	ID       string `json:"Id"`
	Name     string `json:"Name"`
	Status   string `json:"Status"`
	Address  string `json:"Address"`
	User     string `json:"User"`
	Password string `json:"Password"`
	Card     []Card `json:"Card"`
}

// Data define
type Data struct {
	DSLAM []DSLAM `json:"DSLAM"`
}

// SessionResponse define session response structure when we establish DSLAM connection
type SessionResponse struct {
	ID     string `json:"ID"`
	Status string `json:"Status"`
}

// CommandOut define DSLAM return
type CommandOut struct {
	CommandOut []string `json:"CommandOut"`
	ReturnCode string   `json:"ReturnCode"`
}

// --------------
// XML Structures
// --------------

// InfoAttr define XML info of commandOut
type InfoAttr struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

// Parameter define XML parameter of commandOut
type Parameter struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

// Card

// XMLCard define XML DSLAM card struct
type XMLCard struct {
	XMLName xml.Name     `xml:"instance"`
	ResID   string       `xml:"res-id"`
	Info    [][]InfoAttr `xml:"info"`
}

// ShowEquipmentSlot Define DSLAM return of command show equipment slot
type ShowEquipmentSlot struct {
	XMLName xml.Name  `xml:"runtime-data"`
	Card    []XMLCard `xml:"hierarchy>hierarchy>hierarchy>instance"`
}

// Port

// XMLPort define XML DSLAM port struct
type XMLPort struct {
	XMLName xml.Name     `xml:"instance"`
	ResID   string       `xml:"res-id"`
	Info    [][]InfoAttr `xml:"info"`
}

// ShowXdslOperationalDataLine Define DSLAM return of command show xdsl operational-data line
type ShowXdslOperationalDataLine struct {
	XMLName xml.Name  `xml:"runtime-data"`
	Port    []XMLPort `xml:"hierarchy>hierarchy>hierarchy>hierarchy>instance"`
}

// Service

// XMLService define XML DSLAM Service struct
type XMLService struct {
	XMLName   xml.Name      `xml:"instance"`
	ResID     string        `xml:"res-id"`
	Parameter [][]Parameter `xml:"parameter"`
}

// ShowXdslOperDataPortIndexBridgePort Define DSLAM return of command show xdsl oper-data-port *Index* bridge-port xml
type ShowXdslOperDataPortIndexBridgePort struct {
	XMLName xml.Name     `xml:"runtime-data"`
	Service []XMLService `xml:"hierarchy>hierarchy>hierarchy>hierarchy>instance"`
}

// ----
// Code
// ----

// Run lunch Website Interface
func Run() {

	fmt.Println("WebSite Lunched")

	// When Client ask /
	http.HandleFunc("/", indexHandler)

}

// indexHandler define WebServer
func indexHandler(w http.ResponseWriter, r *http.Request) {

	// Define Folders
	var templateBox, _ = rice.FindBox("templates")
	var viewBox, _ = rice.FindBox("views")
	var cssBox, _ = rice.FindBox("static-files/css")
	var fontBox, _ = rice.FindBox("static-files/fonts")
	var imageBox, _ = rice.FindBox("static-files/images")
	var jsBox, _ = rice.FindBox("static-files/js")
	var headerTemplate, _ = templateBox.String("header.tmpl")
	var footerTemplate, _ = templateBox.String("footer.tmpl")
	var indexView, _ = viewBox.String("index.tmpl")
	var optionView, _ = viewBox.String("options.tmpl")

	// Read File
	var dataFile Data
	dataFile = readFile("data.json")

	switch r.URL.Path {
	// Return Files when the user ask the root path (index)
	case "/":
		index := headerTemplate
		index += indexView
		index += footerTemplate
		response, _ := template.New("index").Parse(index)
		DSLAMList := ""
		for i := 0; i < len(dataFile.DSLAM); i++ {
			if dataFile.DSLAM[i].Status != "OK" {
				DSLAMList += "<button class=\"list-button pure-button pure-u-1\" disabled id=\"" + dataFile.DSLAM[i].ID + "\"><span class=\"fa-stack fa-lg\"><i class=\"fa fa-cube fa-stack-2x\"></i><i class=\"fa fa-close fa-stack-2x\"></i></span><span class=\"textButtonDSLAM\">" + dataFile.DSLAM[i].Name + "</span></a>"
			} else {
				DSLAMList += "<button onclick=\"getDslam('" + dataFile.DSLAM[i].ID + "')\" class=\"list-button pure-button pure-u-1\"><span class=\"fa-stack fa-lg\"><i class=\"fa fa-cube fa-stack-2x\"></i></span><span class=\"textButtonDSLAM\">" + dataFile.DSLAM[i].Name + "</span></a>"
			}
		}
		CardList := ""
		PortList := ""
		response.Execute(w, map[string]string{"DSLAMList": DSLAMList, "CardList": CardList, "PortList": PortList})
	// Manage DSLAM
	case "/DSLAM":
		switch r.Method {
		// GET DSLAM Informations when we are in option page
		case "GET":
			option := headerTemplate
			option += optionView
			option += footerTemplate
			response, _ := template.New("option").Parse(option)
			re := regexp.MustCompile("[a-z0-9\\-]*$")
			DSLAMList := ""
			OptionList := ""
			if r.URL.RawQuery != "" {
				dslamID := re.FindString(r.URL.RawQuery)
				dslamPos := getDslamPosByID(dataFile, dslamID)
				for i := 0; i < len(dataFile.DSLAM); i++ {
					if dataFile.DSLAM[i].Status != "OK" {
						DSLAMList += "<a class=\"list-button pure-button pure-u-1\" href=\"/DSLAM?id=" + dataFile.DSLAM[i].ID + "\"><span class=\"fa-stack fa-lg\"><i class=\"fa fa-cube fa-stack-2x\"></i><i class=\"fa fa-close fa-stack-2x\"></i></span><span class=\"textButtonDSLAM\">" + dataFile.DSLAM[i].Name + "</span></a>"
					} else {
						DSLAMList += "<a class=\"list-button pure-button pure-u-1\" href=\"/DSLAM?id=" + dataFile.DSLAM[i].ID + "\"><span class=\"fa-stack fa-lg\"><i class=\"fa fa-cube fa-stack-2x\"></i></span><span class=\"textButtonDSLAM\">" + dataFile.DSLAM[i].Name + "</span></a>"
					}
				}
				OptionList = "<form style=\"margin-top:30px\" class=\"pure-form pure-form-aligned\" action=\"/DSLAM?\" method=\"POST\"><fieldset>"
				OptionList += "<div class=\"pure-control-group\"><label for=\"name\">Nom</label><input id=\"name\" name=\"name\" type=\"text\" value=" + dataFile.DSLAM[dslamPos].Name + "></div>"
				OptionList += "<div class=\"pure-control-group\"><label for=\"address\">Adresse</label><input id=\"address\" name=\"address\" type=\"text\" value=" + dataFile.DSLAM[dslamPos].Address + "></div>"
				OptionList += "<div class=\"pure-control-group\"><label for=\"user\">Utilisateur</label><input id=\"user\" name=\"user\" type=\"text\" value=" + dataFile.DSLAM[dslamPos].User + "></div>"
				OptionList += "<div class=\"pure-control-group\"><label for=\"password\">Mot de Passe</label><input id=\"password\" name=\"password\" type=\"password\" value=" + dataFile.DSLAM[dslamPos].Password + "></div>"
				OptionList += "<input type=\"hidden\" name=\"id\" id=\"id\" value=" + dataFile.DSLAM[dslamPos].ID + ">"
				OptionList += "<button type=\"submit\" style=\"margin-left:180px\" class=\"pure-button-primary pure-button\">Envoyer</button>"
				OptionList += "</fieldset></form>"
				OptionList += "<button onclick=\"sendDelete()\" class=\"button-error pure-button\" style=\"margin-left:180px\">Supprimer</button>"
			}
			response.Execute(w, map[string]string{"DSLAMList": DSLAMList, "Options": OptionList})
		// POST new DSLAM when we are in option page
		case "POST":
			option := headerTemplate
			option += optionView
			option += footerTemplate
			response, _ := template.New("option").Parse(option)
			DSLAMList := ""
			r.ParseForm()
			var newDSLAM DSLAM
			newDSLAM.Name = r.Form.Get("name")
			newDSLAM.Address = r.Form.Get("address")
			newDSLAM.User = r.Form.Get("user")
			newDSLAM.Password = r.Form.Get("password")
			if r.Form.Get("id") != "" {
				newDSLAM.ID = r.Form.Get("id")
				dslamPos := getDslamPosByID(dataFile, r.Form.Get("id"))
				var oldDSLAM DSLAM
				oldDSLAM = getDslamByID(dataFile, newDSLAM.ID)
				newDSLAM.ID = r.Form.Get("id")
				newDSLAM.Card = oldDSLAM.Card
				//replace DSLAM
				dataFile.DSLAM[dslamPos] = newDSLAM
			} else {
				// Add DSLAM
				newDSLAM.ID = uuid.New()
				dataFile.DSLAM = append(dataFile.DSLAM, newDSLAM)
			}
			writeFile(dataFile, "data.json")

			go writeStatus(dataFile, newDSLAM.ID)
			// Create Response Body
			for i := 0; i < len(dataFile.DSLAM); i++ {
				if dataFile.DSLAM[i].Status != "OK" {
					DSLAMList += "<a class=\"list-button pure-button pure-u-1\" href=\"/DSLAM?id=" + dataFile.DSLAM[i].ID + "\"><span class=\"fa-stack fa-lg\"><i class=\"fa fa-cube fa-stack-2x\"></i><i class=\"fa fa-close fa-stack-2x\"></i></span><span class=\"textButtonDSLAM\">" + dataFile.DSLAM[i].Name + "</span></a>"
				} else {
					DSLAMList += "<a class=\"list-button pure-button pure-u-1\" href=\"/DSLAM?id=" + dataFile.DSLAM[i].ID + "\"><span class=\"fa-stack fa-lg\"><i class=\"fa fa-cube fa-stack-2x\"></i></span><span class=\"textButtonDSLAM\">" + dataFile.DSLAM[i].Name + "</span></a>"
				}
			}
			response.Execute(w, map[string]string{"DSLAMList": DSLAMList, "Options": ""})
		// DELETE DSLAM when we are in option page
		case "DELETE":
			if r.URL.RawQuery != "" {
				// Extract ID
				re := regexp.MustCompile("[a-z0-9\\-]*$")
				// Run Regex on String
				dslamID := re.FindString(r.URL.RawQuery)
				// Get Pos
				dslamPos := getDslamPosByID(dataFile, dslamID)
				// Remove DSLAM
				dataFile.DSLAM = append(dataFile.DSLAM[:dslamPos], dataFile.DSLAM[dslamPos+1:]...)
				// Write JSON in File
				writeFile(dataFile, "data.json")
				del := ""
				response, _ := template.New("index").Parse(del)
				response.Execute(w, map[string]string{"Delete": "OK"})
			} else {
				del := ""
				response, _ := template.New("index").Parse(del)
				response.Execute(w, map[string]string{"Delete": "Need DSLAM ID"})
			}
		}
	// Return option page with DSLAM list
	case "/option":
		option := headerTemplate
		option += optionView
		option += footerTemplate
		response, _ := template.New("option").Parse(option)
		DSLAMList := ""
		for i := 0; i < len(dataFile.DSLAM); i++ {
			if dataFile.DSLAM[i].Status != "OK" {
				DSLAMList += "<a class=\"list-button pure-button pure-u-1\" href=\"/DSLAM?id=" + dataFile.DSLAM[i].ID + "\"><span class=\"fa-stack fa-lg\"><i class=\"fa fa-cube fa-stack-2x\"></i><i class=\"fa fa-close fa-stack-2x\"></i></span><span class=\"textButtonDSLAM\">" + dataFile.DSLAM[i].Name + "</span></a>"
			} else {
				DSLAMList += "<a class=\"list-button pure-button pure-u-1\" href=\"/DSLAM?id=" + dataFile.DSLAM[i].ID + "\"><span class=\"fa-stack fa-lg\"><i class=\"fa fa-cube fa-stack-2x\"></i></span><span class=\"textButtonDSLAM\">" + dataFile.DSLAM[i].Name + "</span></a>"
			}
		}
		Options := ""
		// When we want to add a new DSLAM ()/option?add)
		switch r.URL.RawQuery {
		case "add":
			Options = "<form style=\"margin-top:30px\" class=\"pure-form pure-form-aligned\" action=\"/DSLAM?\" method=\"POST\"><fieldset>"
			Options += "<div class=\"pure-control-group\"><label for=\"name\">Nom</label><input id=\"name\" name=\"name\" type=\"text\" placeholder=\"Name\"></div>"
			Options += "<div class=\"pure-control-group\"><label for=\"address\">Adresse</label><input id=\"address\" name=\"address\" type=\"text\" placeholder=\"Adresse\"></div>"
			Options += "<div class=\"pure-control-group\"><label for=\"user\">Utilisateur</label><input id=\"user\" name=\"user\" type=\"text\" placeholder=\"User\"></div>"
			Options += "<div class=\"pure-control-group\"><label for=\"password\">Mot de Passe</label><input id=\"password\" name=\"password\" type=\"password\" placeholder=\"Password\"></div>"
			Options += "<button type=\"submit\" style=\"margin-left:180px\" class=\"pure-button-primary pure-button\">Envoyer</button>"
			Options += "</fieldset></form>"
		default:
			Options = ""
		}
		response.Execute(w, map[string]string{"DSLAMList": DSLAMList, "Options": Options})
	// Web API
	// GET all datas
	case "/SITEAPI/all":
		w.Header().Set("Content-Type", "application/json")
		data, _ := json.MarshalIndent(dataFile, "", "\t")
		w.Write(data)
	// Update data file on server (with real information -> SSH)
	case "/SITEAPI/update":
		w.Header().Set("Content-Type", "application/json")
		i := 0
		for i = 0; i < len(dataFile.DSLAM); i++ {
			var dslamID = dataFile.DSLAM[i].ID
			fmt.Println("Update DSLAM " + dslamID)
			writeStatus(dataFile, dslamID)
			if dataFile.DSLAM[i].Status == "OK" {
				sessionID := getSSHSession(dataFile, dslamID)
				if sessionID != "SSH_KO" {
					writeCard(dataFile, sessionID, dslamID)
				}
			}
		}
		fmt.Println("Data File Updated")
		w.Write([]byte("OK"))
	// Update services on DSLAM with post informations (interface name, internet, voip, iptv)
	case "/SITEAPI/services":
		if r.URL.RawQuery != "" {
			params, _ := url.ParseQuery(r.URL.RawQuery)
			name := params.Get("portName")
			internet := params.Get("internetSwitch")
			voip := params.Get("voipSwitch")
			iptv := params.Get("iptvSwitch")
			dslamID := params.Get("dslamID")
			slot := params.Get("slot")
			portIndex := params.Get("portIndex")
			sessionID := getSSHSession(dataFile, dslamID)
			dslam := getDslamByID(dataFile, dslamID)
			var oldService []Service
			for i := 0; i < len(dslam.Card); i++ {
				if dslam.Card[i].Slot == slot {
					for j := 0; j < len(dslam.Card[i].Port); i++ {
						if dslam.Card[i].Port[j].ID == portIndex {
							oldService = dslam.Card[i].Port[j].Service
						}
					}
				}
			}
			oldInternet, oldVoip, oldIptv := false, false, false
			for k := 0; k < len(oldService); k++ {
				if oldService[k].Vlan == "10" {
					oldInternet = true
				} else if oldService[k].Vlan == "20" {
					oldVoip = true
				} else if oldService[k].Vlan == "30" {
					oldIptv = true
				}
			}
			// Update internet service
			if internet == "true" && oldInternet == false {
				_ = getCommandOut(sessionID, "configure bridge port "+portIndex+":8:35")
				_ = getCommandOut(sessionID, "configure bridge port "+portIndex+":8:35 vlan-id 10")
				_ = getCommandOut(sessionID, "configure bridge port "+portIndex+":8:35 pvid 10")
				writeServiceOnePort(dataFile, sessionID, dslamID, portIndex)
			} else if oldInternet == false {
			} else if internet == "true" && oldInternet == true {
			} else if oldInternet == true {
				_ = getCommandOut(sessionID, "configure bridge no port "+portIndex+":8:35")
				writeServiceOnePort(dataFile, sessionID, dslamID, portIndex)
			}
			// Update voip service
			if voip == "true" && oldVoip == false {
				_ = getCommandOut(sessionID, "configure bridge port "+portIndex+":8:36")
				_ = getCommandOut(sessionID, "configure bridge port "+portIndex+":8:36 vlan-id 20")
				_ = getCommandOut(sessionID, "configure bridge port "+portIndex+":8:36 pvid 20")
				writeServiceOnePort(dataFile, sessionID, dslamID, portIndex)
			} else if oldVoip == false {
			} else if voip == "true" && oldVoip == true {
			} else if oldVoip == true {
				_ = getCommandOut(sessionID, "configure bridge no port "+portIndex+":8:36")
				writeServiceOnePort(dataFile, sessionID, dslamID, portIndex)
			}
			// Update iptv service
			if iptv == "true" && oldIptv == false {
				_ = getCommandOut(sessionID, "configure bridge port "+portIndex+":8:37")
				_ = getCommandOut(sessionID, "configure bridge port "+portIndex+":8:37 vlan-id 30")
				_ = getCommandOut(sessionID, "configure bridge port "+portIndex+":8:37 pvid 30")
				writeServiceOnePort(dataFile, sessionID, dslamID, portIndex)
			} else if oldIptv == false {
			} else if iptv == "true" && oldIptv == true {
			} else if oldIptv == true {
				_ = getCommandOut(sessionID, "configure bridge no port "+portIndex+":8:37")
				writeServiceOnePort(dataFile, sessionID, dslamID, portIndex)
			}

			for i := 0; i < len(dslam.Card); i++ {
				for j := 0; j < len(dslam.Card[i].Port); j++ {
					if dslam.Card[i].Port[j].ID == portIndex {
						dslam.Card[i].Port[j].Name = name
					}
				}
			}
			w.Header().Set("Content-Type", "application/json")
			data, _ := json.MarshalIndent(dataFile, "", "\t")
			w.Write(data)
		} else {
			w.Write([]byte("Need Data"))
		}

	// Send commande to DSLAM
	case "/SITEAPI/command":
		if r.URL.RawQuery != "" {
			reSession := regexp.MustCompile("sessionID=[a-z0-9\\-]*")
			session := reSession.FindString(r.URL.RawQuery)
			resessionID := regexp.MustCompile("[a-z0-9\\-]*$")
			sessionID := resessionID.FindString(session)
			reCommand := regexp.MustCompile("command=[a-zA-Z0-9\\/\\-%]*$")
			command := reCommand.FindString(r.URL.RawQuery)
			reCommandRaw := regexp.MustCompile("[a-zA-Z0-9\\/\\-%]*$")
			commandRaw := reCommandRaw.FindString(command)
			commandRawReplace := strings.Replace(commandRaw, "%20", " ", -1)
			commandOut := getCommandOut(sessionID, commandRawReplace)
			out := "["
			for i := 0; i < len(commandOut); i++ {
				out += "\"" + commandOut[i] + "\""
				if i != len(commandOut)-1 {
					out += ","
				}
			}
			out += "]"
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("{\"commandOut\": " + out + "}"))
		} else {
			w.Write([]byte("Need Session ID and Command"))
		}

	// Static Files
	case "/css/pure-min.css":
		response, _ := cssBox.String("pure-min.css")
		w.Header().Set("Content-Type", "text/css")
		w.Write([]byte(response))
	case "/css/grids-responsive-min.css":
		response, _ := cssBox.String("grids-responsive-min.css")
		w.Header().Set("Content-Type", "text/css")
		w.Write([]byte(response))
	case "/css/grids-responsive-old-ie-min.css":
		response, _ := cssBox.String("grids-responsive-old-ie-min.css")
		w.Header().Set("Content-Type", "text/css")
		w.Write([]byte(response))
	case "/css/WebManageAlcatel.css":
		response, _ := cssBox.String("WebManageAlcatel.css")
		w.Header().Set("Content-Type", "text/css")
		w.Write([]byte(response))
	case "/css/WebManageAlcatel-old-ie.css":
		response, _ := cssBox.String("WebManageAlcatel-old-ie.css")
		w.Header().Set("Content-Type", "text/css")
		w.Write([]byte(response))
	case "/css/font-awesome.min.css":
		response, _ := cssBox.String("font-awesome.min.css")
		w.Header().Set("Content-Type", "text/css")
		w.Write([]byte(response))
	case "/fonts/fontawesome-webfont.woff":
		response, _ := fontBox.String("fontawesome-webfont.woff")
		w.Header().Set("Content-Type", "*/*")
		w.Write([]byte(response))
	case "/fonts/fontawesome-webfont.ttf":
		response, _ := fontBox.String("fontawesome-webfont.ttf")
		w.Header().Set("Content-Type", "*/*")
		w.Write([]byte(response))
	case "/fonts/fontawesome-webfont.eot":
		response, _ := fontBox.String("fontawesome-webfont.eot")
		w.Header().Set("Content-Type", "*/*")
		w.Write([]byte(response))
	case "/js/WebManageAlcatel.js":
		response, _ := jsBox.String("WebManageAlcatel.js")
		w.Header().Set("Content-Type", "text/javascript")
		w.Write([]byte(response))
	case "/favicon.png":
		response, _ := imageBox.String("favicon.png")
		w.Header().Set("Content-Type", "image/png")
		w.Write([]byte(response))
	case "/background.jpg":
		response, _ := imageBox.String("3587308442_ce1c047329_o.jpg")
		w.Header().Set("Content-Type", "image/jpg")
		w.Write([]byte(response))
	case "/mitel.png":
		response, _ := imageBox.String("mitel.png")
		w.Header().Set("Content-Type", "image/png")
		w.Write([]byte(response))
	default:
	}

}

// GET DSLAM By ID in data file
func getDslamByID(data Data, id string) DSLAM {

	for i := 0; i < len(data.DSLAM); i++ {
		if data.DSLAM[i].ID == id {
			return data.DSLAM[i]
		}
	}
	var null DSLAM
	return null

}

// GET DSLAM Position By ID in data file
func getDslamPosByID(data Data, id string) int {

	for i := 0; i < len(data.DSLAM); i++ {
		if data.DSLAM[i].ID == id {
			return i
		}
	}
	return -1

}

// GET SSH Session to connect to DSLAM in data file
func getSSHSession(dataFile Data, id string) string {
	var dslam DSLAM
	dslam = getDslamByID(dataFile, id)
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8080/API/session", bytes.NewBufferString("{\"user\": \""+dslam.User+"\", \"host\": \""+dslam.Address+"\", \"password\": \""+dslam.Password+"\"}"))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var response SessionResponse
	var err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("JSON Unmarshal body in getSSHSession error:", err)
	}

	var status string
	if response.Status != "OK" {
		status = "SSH_KO"
	} else {
		status = response.ID
	}
	return status
}

// GET Commend Out with session ID and Command
func getCommandOut(sessionID, command string) []string {
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

	var response CommandOut
	var err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("JSON Unmarshal body in getCommandOut error:", err)
	}

	return response.CommandOut
}

// Write DSLAM Status in data File
func writeStatus(dataFile Data, id string) {

	fmt.Println("Write DSLAM Status")

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

	writeFile(dataFile, "data.json")

}

// Write Card list of DSLAM in data file
func writeCard(dataFile Data, sessionID, dslamID string) {

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
	writeFile(dataFile, "data.json")

	writePort(dataFile, sessionID, dslamID)
}

// Write Port List of DSLAM in data file
func writePort(dataFile Data, sessionID, dslamID string) {

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
	writeFile(dataFile, "data.json")

	writeService(dataFile, sessionID, dslamID)
}

// Write Service List of DSLAM in data file
func writeService(dataFile Data, sessionID, dslamID string) {
	fmt.Println("Update Service")

	dslam := getDslamByID(dataFile, dslamID)

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
					writeFile(dataFile, "data.json")
				}
			}
		}
	}
}

// Write Service of One Port in data file when we update service
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
				dslam.Card[i].Port[j].Service = service
			}
		}
	}

	dslamPos := getDslamPosByID(dataFile, dslamID)
	dataFile.DSLAM[dslamPos] = dslam
	writeFile(dataFile, "data.json")

	fmt.Println("Data File Updated")
}

// Write File
func writeFile(dataFile Data, file string) {
	jsonIndent, _ := json.MarshalIndent(dataFile, "", "\t")
	key := []byte("nit5i8drod0it1of0en8hylg9ov0out4") // 32 bytes!
	plaintext := []byte(jsonIndent)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	ioutil.WriteFile(file, ciphertext, 0777)
	//readFile(file)
}

func readFile(file string) Data {
	var ciphertext, _ = ioutil.ReadFile(file)

	key := []byte("nit5i8drod0it1of0en8hylg9ov0out4") // 32 bytes!

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(ciphertext, ciphertext)

	var dataFile Data
	err = json.Unmarshal(ciphertext, &dataFile)
	if err != nil {
		fmt.Println("JSON Unmarshal file in indexHandler error:", err)
	}

	return dataFile

}
