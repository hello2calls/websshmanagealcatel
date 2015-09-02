package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"text/template"

	"code.google.com/p/go-uuid/uuid"
	"github.com/GeertJohan/go.rice"

	"bitbucket.org/nmontes/WebSSHManageAlcatel/web/pkg/equipment"
	"bitbucket.org/nmontes/WebSSHManageAlcatel/web/pkg/file"
	S "bitbucket.org/nmontes/WebSSHManageAlcatel/web/pkg/structures"
	WD "bitbucket.org/nmontes/WebSSHManageAlcatel/web/pkg/writeData"
)

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
	var dataFile S.Data
	dataFile = file.ReadFile("data.json")

	switch r.URL.Path {
	// Return Files when the user ask the root path (index)
	case "/":
		//Check if data file contain DSLAM and redirect to option if not
		if len(dataFile.DSLAM) == 0 {
			http.Redirect(w, r, "/option", http.StatusFound)
		} else {
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
		}
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
				dslamPos := equipment.GetDslamPosByID(dataFile, dslamID)
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
			var newDSLAM S.DSLAM
			newDSLAM.Name = r.Form.Get("name")
			newDSLAM.Address = r.Form.Get("address")
			newDSLAM.User = r.Form.Get("user")
			newDSLAM.Password = r.Form.Get("password")
			if r.Form.Get("id") != "" {
				newDSLAM.ID = r.Form.Get("id")
				dslamPos := equipment.GetDslamPosByID(dataFile, r.Form.Get("id"))
				var oldDSLAM S.DSLAM
				oldDSLAM = equipment.GetDslamByID(dataFile, newDSLAM.ID)
				newDSLAM.ID = r.Form.Get("id")
				newDSLAM.Card = oldDSLAM.Card
				//replace DSLAM
				dataFile.DSLAM[dslamPos] = newDSLAM
			} else {
				// Add DSLAM
				newDSLAM.ID = uuid.New()
				dataFile.DSLAM = append(dataFile.DSLAM, newDSLAM)
			}
			file.WriteFile(dataFile, "data.json")

			go WD.WriteStatus(dataFile, newDSLAM.ID)
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
				dslamPos := equipment.GetDslamPosByID(dataFile, dslamID)
				// Remove DSLAM
				dataFile.DSLAM = append(dataFile.DSLAM[:dslamPos], dataFile.DSLAM[dslamPos+1:]...)
				// Write JSON in File
				file.WriteFile(dataFile, "data.json")
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
			WD.WriteStatus(dataFile, dslamID)
			if dataFile.DSLAM[i].Status == "OK" {
				sessionID := getSSHSession(dataFile, dslamID)
				if sessionID != "SSH_KO" {
					WD.WriteCard(dataFile, sessionID, dslamID)
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
			dslam := equipment.GetDslamByID(dataFile, dslamID)
			var oldService []S.Service
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
				WD.WriteServiceOnePort(dataFile, sessionID, dslamID, portIndex)
			} else if oldInternet == false {
			} else if internet == "true" && oldInternet == true {
			} else if oldInternet == true {
				_ = getCommandOut(sessionID, "configure bridge no port "+portIndex+":8:35")
				WD.WriteServiceOnePort(dataFile, sessionID, dslamID, portIndex)
			}
			// Update voip service
			if voip == "true" && oldVoip == false {
				_ = getCommandOut(sessionID, "configure bridge port "+portIndex+":8:36")
				_ = getCommandOut(sessionID, "configure bridge port "+portIndex+":8:36 vlan-id 20")
				_ = getCommandOut(sessionID, "configure bridge port "+portIndex+":8:36 pvid 20")
				WD.WriteServiceOnePort(dataFile, sessionID, dslamID, portIndex)
			} else if oldVoip == false {
			} else if voip == "true" && oldVoip == true {
			} else if oldVoip == true {
				_ = getCommandOut(sessionID, "configure bridge no port "+portIndex+":8:36")
				WD.WriteServiceOnePort(dataFile, sessionID, dslamID, portIndex)
			}
			// Update iptv service
			if iptv == "true" && oldIptv == false {
				_ = getCommandOut(sessionID, "configure bridge port "+portIndex+":8:37")
				_ = getCommandOut(sessionID, "configure bridge port "+portIndex+":8:37 vlan-id 30")
				_ = getCommandOut(sessionID, "configure bridge port "+portIndex+":8:37 pvid 30")
				WD.WriteServiceOnePort(dataFile, sessionID, dslamID, portIndex)
			} else if oldIptv == false {
			} else if iptv == "true" && oldIptv == true {
			} else if oldIptv == true {
				_ = getCommandOut(sessionID, "configure bridge no port "+portIndex+":8:37")
				WD.WriteServiceOnePort(dataFile, sessionID, dslamID, portIndex)
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

// GET SSH Session to connect to DSLAM in data file
func getSSHSession(dataFile S.Data, id string) string {
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

	var response S.CommandOut
	var err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("JSON Unmarshal body in getCommandOut error:", err)
	}

	return response.CommandOut
}
