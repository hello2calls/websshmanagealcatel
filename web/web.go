package web

import (
	"fmt"
	"net/http"
	"text/template"
	"github.com/GeertJohan/go.rice"
	"code.google.com/p/go-uuid/uuid"
	"regexp"
	"encoding/json"
	"bytes"
	"io/ioutil"
	"strings"
)

//Service Structure
type Service struct {
	Id string `json:"Id"`
	Name string `json:"Name"`
	Status string `json:"Status"`
}

// Port Structure
type Port struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Status string `json:"status"`
	Service []Service `json:"Service"`
}

// Card Structure
type Card struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Status string `json:"status"`
	Port []Port `json:"Port"`
}


// DSLAM Structure
type DSLAM struct {
	Id string `json:"Id"`
	Name string `json:"Name"`
	Status string `json:"Status"`
	Address string `json:"Address"`
	User string `json:"User"`
	Password string `json:"Password"`
	Card []Card `json:"Card"`
}

// Data Structure
type Data struct {
	DSLAM []DSLAM `json:"DSLAM"`
}

// Session Response
type SessionResponse struct {
	ID string `json:"ID"`
	Status string `json:"Status"`
}

type CommandOut struct {
	CommandOut []string `json:"CommandOut"`
	ReturnCode string `json:"ReturnCode"`
}

func Run() {

	fmt.Println("WebSite Lunched")

	// When Client ask /
	http.HandleFunc("/", indexHandler)

}


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
	var file, _ = ioutil.ReadFile("web/data.json")
	var dataFile Data
	var _ = json.Unmarshal(file, &dataFile)

	// Return Files
	switch r.URL.Path {
		case "/":
			index := headerTemplate
			index += indexView
			index += footerTemplate
			response, _ := template.New("index").Parse(index)
			DSLAMList := ""
			for i := 0; i < len(dataFile.DSLAM); i++ {
				if dataFile.DSLAM[i].Status != "OK" {
					DSLAMList += "<button class=\"list-button pure-button pure-u-1\" disabled id=\""+ dataFile.DSLAM[i].Id +"\"><span class=\"fa-stack fa-lg\"><i class=\"fa fa-cube fa-stack-2x\"></i><i class=\"fa fa-close fa-stack-2x\"></i></span><span class=\"textButtonDSLAM\">"+ dataFile.DSLAM[i].Name +"</span></a>"
				} else {
					DSLAMList += "<button onclick=\"getDslam('"+dataFile.DSLAM[i].Id+"')\" class=\"list-button pure-button pure-u-1\"><span class=\"fa-stack fa-lg\"><i class=\"fa fa-cube fa-stack-2x\"></i></span><span class=\"textButtonDSLAM\">"+ dataFile.DSLAM[i].Name +"</span></a>"				}
			}
			CardList := "<button class=\"pure-button pure-u-1\"><span class=\"fa-stack fa-lg\"><i class=\"imageCard fa fa-hdd-o fa-stack-2x\"></i><i class=\"imageCard fa fa-exclamation fa-stack-2x\"></i></span><span class=\"textButtonCard\">NVLT-N</span></button>"
			PortList := "<button class=\"pure-button pure-u-1\"><span class=\"fa-stack fa-lg\"><i class=\"imagePort fa fa-caret-square-o-right\"></i></span><span class=\"textButtonPort\">Chambre 1</span></button>"
			response.Execute(w, map[string]string{"DSLAMList": DSLAMList, "CardList": CardList, "PortList": PortList})
		case "/DSLAM":
			switch r.Method {
				case "GET":
					option := headerTemplate
					option += optionView
					option += footerTemplate
					response, _ := template.New("option").Parse(option)
					re := regexp.MustCompile("[a-z0-9\\-]*$")
					DSLAMid := re.FindString(r.URL.RawQuery)
					dslamPos := getDslamPosById(dataFile, DSLAMid)
					DSLAMList := ""
					OptionList := ""
					for i := 0; i < len(dataFile.DSLAM); i++ {
						if dataFile.DSLAM[i].Status != "OK" {
							DSLAMList += "<a class=\"list-button pure-button pure-u-1\" href=\"/DSLAM?id="+ dataFile.DSLAM[i].Id +"\"><span class=\"fa-stack fa-lg\"><i class=\"fa fa-cube fa-stack-2x\"></i><i class=\"fa fa-close fa-stack-2x\"></i></span><span class=\"textButtonDSLAM\">"+ dataFile.DSLAM[i].Name +"</span></a>"
							} else {
								DSLAMList += "<a class=\"list-button pure-button pure-u-1\" href=\"/DSLAM?id="+ dataFile.DSLAM[i].Id +"\"><span class=\"fa-stack fa-lg\"><i class=\"fa fa-cube fa-stack-2x\"></i></span><span class=\"textButtonDSLAM\">"+ dataFile.DSLAM[i].Name +"</span></a>"
							}
					}
					OptionList = "<form style=\"margin-top:30px\" class=\"pure-form pure-form-aligned\" action=\"/DSLAM?\" method=\"POST\"><fieldset>"
					OptionList += "<div class=\"pure-control-group\"><label for=\"name\">Nom</label><input id=\"name\" name=\"name\" type=\"text\" value="+ dataFile.DSLAM[dslamPos].Name +"></div>"
					OptionList += "<div class=\"pure-control-group\"><label for=\"address\">Adresse</label><input id=\"address\" name=\"address\" type=\"text\" value="+ dataFile.DSLAM[dslamPos].Address +"></div>"
					OptionList += "<div class=\"pure-control-group\"><label for=\"user\">Utilisateur</label><input id=\"user\" name=\"user\" type=\"text\" value="+ dataFile.DSLAM[dslamPos].User +"></div>"
					OptionList += "<div class=\"pure-control-group\"><label for=\"password\">Mot de Passe</label><input id=\"password\" name=\"password\" type=\"password\" value="+ dataFile.DSLAM[dslamPos].Password +"></div>"
					OptionList += "<input type=\"hidden\" name=\"id\" id=\"id\" value="+ dataFile.DSLAM[dslamPos].Id +">"
					OptionList += "<button type=\"submit\" style=\"margin-left:180px\" class=\"pure-button-primary pure-button\">Envoyer</button>"
					OptionList += "</fieldset></form>"
					OptionList += "<button onclick=\"sendDelete()\" class=\"button-error pure-button\" style=\"margin-left:180px\">Supprimer</button>"
					response.Execute(w, map[string]string{"DSLAMList": DSLAMList, "Options": OptionList})
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
						newDSLAM.Id = r.Form.Get("id")
						dslamPos := getDslamPosById(dataFile, r.Form.Get("id"))
						var oldDSLAM DSLAM
						oldDSLAM = getDslamById(dataFile, newDSLAM.Id)
						newDSLAM.Id = r.Form.Get("id")
						newDSLAM.Card = oldDSLAM.Card
						//replace DSLAM
						dataFile.DSLAM[dslamPos] = newDSLAM
					} else {
						// Add DSLAM
						newDSLAM.Id = uuid.New()
						dataFile.DSLAM = append(dataFile.DSLAM, newDSLAM)
					}
					jsonIndent, _ := json.MarshalIndent(dataFile, "", "\t")
					ioutil.WriteFile("web/data.json", jsonIndent, 0777)
					go writeStatus(newDSLAM.Id, dataFile)
					// Create Response Body
					for i := 0; i < len(dataFile.DSLAM); i++ {
						if dataFile.DSLAM[i].Status != "OK" {
							DSLAMList += "<a class=\"list-button pure-button pure-u-1\" href=\"/DSLAM?id="+ dataFile.DSLAM[i].Id +"\"><span class=\"fa-stack fa-lg\"><i class=\"fa fa-cube fa-stack-2x\"></i><i class=\"fa fa-close fa-stack-2x\"></i></span><span class=\"textButtonDSLAM\">"+ dataFile.DSLAM[i].Name +"</span></a>"
						} else {
							DSLAMList += "<a class=\"list-button pure-button pure-u-1\" href=\"/DSLAM?id="+ dataFile.DSLAM[i].Id +"\"><span class=\"fa-stack fa-lg\"><i class=\"fa fa-cube fa-stack-2x\"></i></span><span class=\"textButtonDSLAM\">"+ dataFile.DSLAM[i].Name +"</span></a>"
						}
					}
					response.Execute(w, map[string]string{"DSLAMList": DSLAMList, "Options": ""})
				case "DELETE":
					// Extract ID
					re := regexp.MustCompile("[a-z0-9\\-]*$")
					// Run Regex on String
					DSLAMid := re.FindString(r.URL.RawQuery)
					// Get Pos
					dslamPos := getDslamPosById(dataFile, DSLAMid)
					// Remove DSLAM
					dataFile.DSLAM = append(dataFile.DSLAM[:dslamPos], dataFile.DSLAM[dslamPos+1:]...)
					// Indent JSON
					jsonIndent, _ := json.MarshalIndent(dataFile, "", "\t")
					// Write JSON in File
					ioutil.WriteFile("web/data.json", jsonIndent, 0777)
					del := ""
					response, _ := template.New("index").Parse(del)
					response.Execute(w, map[string]string{"Delete": "OK"})
			}
		case "/option":
			option := headerTemplate
			option += optionView
			option += footerTemplate
			response, _ := template.New("option").Parse(option)
			DSLAMList := ""
			for i := 0; i < len(dataFile.DSLAM); i++ {
				if dataFile.DSLAM[i].Status != "OK" {
					DSLAMList += "<a class=\"list-button pure-button pure-u-1\" href=\"/DSLAM?id="+ dataFile.DSLAM[i].Id +"\"><span class=\"fa-stack fa-lg\"><i class=\"fa fa-cube fa-stack-2x\"></i><i class=\"fa fa-close fa-stack-2x\"></i></span><span class=\"textButtonDSLAM\">"+ dataFile.DSLAM[i].Name +"</span></a>"
				} else {
					DSLAMList += "<a class=\"list-button pure-button pure-u-1\" href=\"/DSLAM?id="+ dataFile.DSLAM[i].Id +"\"><span class=\"fa-stack fa-lg\"><i class=\"fa fa-cube fa-stack-2x\"></i></span><span class=\"textButtonDSLAM\">"+ dataFile.DSLAM[i].Name +"</span></a>"
				}
			}
			Options := ""
			switch r.URL.RawQuery {
				case "add":
					Options = "<form style=\"margin-top:30px\" class=\"pure-form pure-form-aligned\" action=\"/DSLAM?\" method=\"POST\"><fieldset>"
					Options += "<div class=\"pure-control-group\"><label for=\"name\">Nom</label><input id=\"name\" name=\"name\" type=\"text\" placeholder=\"Name\"></div>"
					Options += "<div class=\"pure-control-group\"><label for=\"address\">Adresse</label><input id=\"address\" name=\"address\" type=\"text\" placeholder=\"Adresse\"></div>"
					Options += "<div class=\"pure-control-group\"><label for=\"user\">Utilisateur</label><input id=\"user\" name=\"user\" type=\"text\" placeholder=\"User\"></div>"
					Options += "<div class=\"pure-control-group\"><label for=\"password\">Mot de Passe</label><input id=\"password\" name=\"password\" type=\"password\" placeholder=\"Password\"></div>"
					Options += "<button type=\"submit\" style=\"margin-left:180px\" class=\"pure-button-primary pure-button\">Envoyer</button>"
					Options += "</fieldset></form>"
				default :
					Options = ""
			}
			response.Execute(w, map[string]string{"DSLAMList": DSLAMList, "Options": Options})
		case "/getDslam":
			re := regexp.MustCompile("[a-z0-9\\-]*$")
			DSLAMid := re.FindString(r.URL.RawQuery)
			dslam := getDslamById(dataFile, DSLAMid)
			out, _ := json.Marshal(dslam)
			sessionId := getSshSession(DSLAMid, dataFile)
			go searchCard(sessionId, dataFile)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("{\"dslam\": "+ string(out) +", \"sessionId\": \"" + sessionId + "\"}"))
		case "/session":
			re := regexp.MustCompile("[a-z0-9\\-]*$")
			DSLAMid := re.FindString(r.URL.RawQuery)
			sessionId := getSshSession(DSLAMid, dataFile)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("{\"sessionId\": \""+ sessionId +"\"}"))
		case "/command":
			reSession := regexp.MustCompile("sessionId=[a-z0-9\\-]*")
			session := reSession.FindString(r.URL.RawQuery)
			reSessionId := regexp.MustCompile("[a-z0-9\\-]*$")
			sessionId := reSessionId.FindString(session)
			reCommand := regexp.MustCompile("command=[a-zA-Z0-9\\/\\-%]*$")
			command := reCommand.FindString(r.URL.RawQuery)
			reCommandRaw := regexp.MustCompile("[a-zA-Z0-9\\/\\-%]*$")
			commandRaw := reCommandRaw.FindString(command)
			commandRawReplace := strings.Replace(commandRaw,"%20"," ",-1)
			commandOut := getCommandOut(sessionId, commandRawReplace)
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
		default:
	}

}


func getDslamById(data Data, id string) (DSLAM) {

	for i := 0; i < len(data.DSLAM); i++ {
		if data.DSLAM[i].Id == id {
			return data.DSLAM[i]
		}
	}
	var null DSLAM
	return null

}


func getDslamPosById(data Data, id string) (int) {

	for i := 0; i < len(data.DSLAM); i++ {
		if data.DSLAM[i].Id == id {
			return i
		}
	}
	return -1

}


func getSshSession(id string, dataFile Data) (string) {
	var dslam DSLAM
	dslam = getDslamById(dataFile, id)
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8080/API/session",  bytes.NewBufferString("{\"user\": \"" + dslam.User + "\", \"host\": \"" + dslam.Address + ":22\", \"password\": \"" + dslam.Password + "\"}"))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var response SessionResponse
	var _ = json.Unmarshal(body, &response)

	if response.Status != "OK" {
		return "SSH_KO"
	} else {
		return response.ID
	}
}


func getCommandOut(sessionId, command string) ([]string) {
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8080/API/command",  bytes.NewBufferString("{\"SessionID\": \"" + sessionId + "\", \"command\": \"" + command + "\"}"))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)


	var response CommandOut
	var _ = json.Unmarshal(body, &response)

	return response.CommandOut
}


func writeStatus(id string, dataFile Data) {
	var oldDSLAM DSLAM
	var newDSLAM DSLAM

	dslamPos := getDslamPosById(dataFile, id)
	oldDSLAM = getDslamById(dataFile, id)

	req, _ := http.NewRequest("POST", "http://127.0.0.1:8080/API/session",  bytes.NewBufferString("{\"user\": \"" + oldDSLAM.User + "\", \"host\": \"" + oldDSLAM.Address + ":22\", \"password\": \"" + oldDSLAM.Password + "\"}"))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var response SessionResponse
	var _ = json.Unmarshal(body, &response)

	dataFile.DSLAM = append(dataFile.DSLAM[:dslamPos], dataFile.DSLAM[dslamPos+1:]...)
	newDSLAM = oldDSLAM

	newDSLAM.Status = response.Status

	dataFile.DSLAM = append(dataFile.DSLAM, newDSLAM)
	jsonIndent, _ := json.MarshalIndent(dataFile, "", "\t")
	ioutil.WriteFile("web/data.json", jsonIndent, 0777)
}

func searchCard(sessionId string, dataFile Data) {
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8080/API/command",  bytes.NewBufferString("{\"SessionID\": \"" + sessionId + "\", \"command\": \"ls\"}"))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	//bodyS := string(body)

	//reXml := regexp.MustCompile("^[\\ \\<]*.*")
	//xml := reXml.FindAllString(bodyS, -1)


	var response CommandOut
	var _ = json.Unmarshal(body, &response)

}
