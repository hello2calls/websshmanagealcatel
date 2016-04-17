package web

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"text/template"

	"github.com/GeertJohan/go.rice"
	"github.com/gorilla/websocket"

	"bitbucket.org/nmontes/WebSSHManageAlcatel/web/controllers/dslam"
	"bitbucket.org/nmontes/WebSSHManageAlcatel/web/controllers/option"
	"bitbucket.org/nmontes/WebSSHManageAlcatel/web/controllers/siteapi"

	"bitbucket.org/nmontes/WebSSHManageAlcatel/web/pkg/file"
	"bitbucket.org/nmontes/WebSSHManageAlcatel/web/pkg/logger"
	S "bitbucket.org/nmontes/WebSSHManageAlcatel/web/pkg/structures"
)

// Run lunch Website Interface
func Run() {

	fmt.Println("WebSite Lunched")
	logger.Print("WebSite Lunched", nil)

	// When Client ask /
	http.HandleFunc("/", indexHandler)

}

// indexHandler define WebServer
func indexHandler(w http.ResponseWriter, r *http.Request) {

	// Used for WebSocket
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

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
			dslam.Get(w, r, dataFile)
		// POST new DSLAM when we are in option page
		case "POST":
			dslam.Post(w, r, dataFile)
		// DELETE DSLAM when we are in option page
		case "DELETE":
			dslam.Delete(w, r, dataFile)
		}

	// Return option page with DSLAM list
	case "/option":
		option.Get(w, r, dataFile)

	// Web API
	// GET all datas
	case "/SITEAPI/all":
		siteapi.All(w, r, dataFile)
	// Update data file on server (with real information -> SSH)
	case "/SITEAPI/update":
		siteapi.Update(w, r, dataFile)
	// Update services on DSLAM with post informations (interface name, internet, voip, iptv)
	case "/SITEAPI/services":
		siteapi.Services(w, r, dataFile)
	// Send commande to DSLAM
	case "/SITEAPI/command":
		siteapi.Command(w, r, dataFile)
	case "/SITEAPI/session":
		switch r.Method {
		// GET DSLAM Informations when we are in option page
		case "DELETE":
			siteapi.DeleteSession(w, r)
		}

	// WebSocket
	case "/ws":
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				return
			}
			fmt.Println(string(p))

			err = conn.WriteMessage(messageType, []byte("Send from webserver to JS"))
			if err != nil {
				fmt.Println("ERROR Write message from webserver to JS")
				return
			}
			fmt.Println("webserver : Message sent to JS")
		}
	case "/APISSHWS":
		wsURL := "ws://localhost:8080/API/ws"
		u, err := url.Parse(wsURL)
		if err != nil {
			fmt.Println(err)
		}
		rawConn, err := net.Dial("tcp", u.Host)
		if err != nil {
			fmt.Println(err)
		}
		wsHeaders := http.Header{
			"Origin": {wsURL},
		}
		wsConn, _, err := websocket.NewClient(rawConn, u, wsHeaders, 1024, 1024)
		if err != nil {
			fmt.Println("ERROR wsConn")
			fmt.Println(err)
		}
		err = wsConn.WriteMessage(websocket.TextMessage, []byte("Send from webserver to SSH API"))
		if err != nil {
			fmt.Println("ERROR Write message from webserver to SSH API ")
			return
		}
		fmt.Println("webserver : Message sent to SSH API")
		for {
			_, p, err := wsConn.ReadMessage()
			if err != nil {
				return
			}
			fmt.Println(string(p))
		}

	// Serve Static Files
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
	case "/logo.png":
		response, _ := imageBox.String("logo.png")
		w.Header().Set("Content-Type", "image/png")
		w.Write([]byte(response))
	default:
	}

}
