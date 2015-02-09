package web

import (
	"fmt"
	"net/http"
	"text/template"
	"github.com/GeertJohan/go.rice"
	"code.google.com/p/go-uuid/uuid"
	"regexp"
	"encoding/json"
	//"os"
	"io/ioutil"
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
	Card []Card `json:"Card"`
}

// Data Structure
type Data struct {
	DSLAM []DSLAM `json:"DSLAM"`
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
			DSLAMList := "<button class=\"pure-button pure-u-1\"><i class=\"imageDSLAM fa fa-cube fa-3x\"></i><span class=\"textButtonDSLAM\">DSLAM 2</span></button>"
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
						DSLAMList += "<a class=\"list-button pure-button pure-u-1\" href=\"/DSLAM?id="+ dataFile.DSLAM[i].Id +"\"><i class=\"imageDSLAM fa fa-cube fa-3x\"></i><span class=\"textButtonDSLAM\">"+ dataFile.DSLAM[i].Name +"</span></a>"
					}
					OptionList = "<form style=\"margin-top:30px\" class=\"pure-form pure-form-aligned\" action=\"/DSLAM?\" method=\"POST\"><fieldset>"
					OptionList += "<div class=\"pure-control-group\"><label for=\"name\">Nom</label><input id=\"name\" name=\"name\" type=\"text\" value="+ dataFile.DSLAM[dslamPos].Name +"></div><div class=\"pure-control-group\"><label for=\"password\">Adresse</label><input id=\"address\" name=\"address\" type=\"text\" value="+ dataFile.DSLAM[dslamPos].Address +"></div><input type=\"hidden\" name=\"id\" id=\"id\" value="+ dataFile.DSLAM[dslamPos].Id +"><button type=\"submit\" style=\"margin-left:180px\" class=\"pure-button-primary pure-button\">Envoyer</button></fieldset></form>"
					OptionList += "<button onclick=\"sendDelete()\" class=\"button-error pure-button\" style=\"margin-left:180px\">Supprimer</button>"
					response.Execute(w, map[string]string{"DSLAMList": DSLAMList, "Options": OptionList})
				case "POST":
					option := headerTemplate
					option += optionView
					option += footerTemplate
					response, _ := template.New("option").Parse(option)
					DSLAMList := ""
					r.ParseForm()
					if r.Form.Get("id") != "" {
						dslamPos := getDslamPosById(dataFile, r.Form.Get("id"))
						//Save Old DSLAM
						var oldDSLAM DSLAM
						oldDSLAM = getDslamById(dataFile, r.Form.Get("id"))
						//Delete DSLAM
						dataFile.DSLAM = append(dataFile.DSLAM[:dslamPos], dataFile.DSLAM[dslamPos+1:]...)
						// Add DSLAM
						var newDSLAM DSLAM
						newDSLAM.Id = r.Form.Get("id")
						newDSLAM.Name = r.Form.Get("name")
						newDSLAM.Address = r.Form.Get("address")
						newDSLAM.Card = oldDSLAM.Card
						dataFile.DSLAM = append(dataFile.DSLAM, newDSLAM)
					} else {
						// Add DSLAM
						var newDSLAM DSLAM
						newDSLAM.Id = uuid.New()
						newDSLAM.Name = r.Form.Get("name")
						newDSLAM.Address = r.Form.Get("address")
						dataFile.DSLAM = append(dataFile.DSLAM, newDSLAM)
					}
					jsonIndent, _ := json.MarshalIndent(dataFile, "", "\t")
					ioutil.WriteFile("web/data.json", jsonIndent, 0777)
					//Connexion au DSLAM pour la liste des Cartes et des Ports.
					//rows, _ := database.Query("SELECT * FROM DSLAM;")
					for i := 0; i < len(dataFile.DSLAM); i++ {
						DSLAMList += "<a class=\"list-button pure-button pure-u-1\" href=\"/DSLAM?id="+ dataFile.DSLAM[i].Id +"\"><i class=\"imageDSLAM fa fa-cube fa-3x\"></i><span class=\"textButtonDSLAM\">"+ dataFile.DSLAM[i].Name +"</span></a>"
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
			}
		case "/option":
			option := headerTemplate
			option += optionView
			option += footerTemplate
			response, _ := template.New("option").Parse(option)
			DSLAMList := ""
			for i := 0; i < len(dataFile.DSLAM); i++ {
				DSLAMList += "<a class=\"list-button pure-button pure-u-1\" href=\"/DSLAM?id="+ dataFile.DSLAM[i].Id +"\"><i class=\"imageDSLAM fa fa-cube fa-3x\"></i><span class=\"textButtonDSLAM\">"+ dataFile.DSLAM[i].Name +"</span></a>"
			}
			Options := ""
			switch r.URL.RawQuery {
				case "add":
					Options = "<form style=\"margin-top:30px\" class=\"pure-form pure-form-aligned\" action=\"/DSLAM?\" method=\"POST\"><fieldset><div class=\"pure-control-group\"><label for=\"name\">Nom</label><input id=\"name\" name=\"name\" type=\"text\" placeholder=\"Name\"></div><div class=\"pure-control-group\"><label for=\"password\">Adresse</label><input id=\"address\" name=\"address\" type=\"text\" placeholder=\"Adresse\"></div><button type=\"submit\" style=\"margin-left:180px\" class=\"pure-button-primary pure-button\">Envoyer</button></fieldset></form>"
				default :
					Options = ""
			}
			response.Execute(w, map[string]string{"DSLAMList": DSLAMList, "Options": Options})

		//Static Files
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
			w.Header().Set("Content-Type", "text/javascrAddresst")
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
