package web

import (
	"fmt"
	"net/http"
	"text/template"
	"github.com/GeertJohan/go.rice"
	"database/sql"
	"github.com/mattn/go-sqlite3"
	"code.google.com/p/go-uuid/uuid"
	"regexp"
)

var database *sql.DB
var tx *sql.Tx

func Run() {

	// Register the driver
	var DB_DRIVER string
	sql.Register(DB_DRIVER, &sqlite3.SQLiteDriver{})
	//Open Database
	var err error
	database, err = sql.Open(DB_DRIVER, "AlcatelDSLAM.database")
	if err != nil {
		fmt.Println("Failed to create the handle")
	}
	//Ping Database
	if err2 := database.Ping(); err2 != nil {
		fmt.Println("Failed to keep connection alive")
	}
	//Start Database
	tx, err = database.Begin()
	if err != nil {
		fmt.Println("Error to create Database")
		fmt.Println(err)
	}
	//Create DSLAM Table
	_, err = database.Exec("CREATE TABLE IF NOT EXISTS DSLAM (id VARCHAR(100) PRIMARY KEY,name VARCHAR(100) NOT NULL,address VARCHAR(100) NOT NULL)",)
	if err != nil {
		fmt.Println("Error to create DSLAM Table")
		fmt.Println(err)
	}
	//Create Card Table
	_, err = database.Exec("CREATE TABLE IF NOT EXISTS Card (id VARCHAR(100) PRIMARY KEY,name VARCHAR(100) NOT NULL, active BOOLEAN NOT NULL,dslam_id INT NOT NULL,FOREIGN KEY (dslam_id) REFERENCES DSLAM(id))",)
	if err != nil {
		fmt.Println("Error to create Card Table")
		fmt.Println(err)
	}
	//Create Port Table
	_, err = database.Exec("CREATE TABLE IF NOT EXISTS DSLAM (id VARCHAR(100) PRIMARY KEY,name VARCHAR(100) NOT NULL, active BOOLEAN NOT NULL,card_id INT NOT NULL,FOREIGN KEY (card_id) REFERENCES Card(id))",)
	if err != nil {
		fmt.Println("Error to create Port Table")
		fmt.Println(err)
	}
	// Commit Changes in Database
	tx.Commit()

	fmt.Println("WebSite Lunched")

	// When Client ask /
	http.HandleFunc("/", indexHandler)

}


func indexHandler(w http.ResponseWriter, r *http.Request) {

	// Define Folders
	templateBox, _ := rice.FindBox("templates")
	viewBox, _ := rice.FindBox("views")
	cssBox, _ := rice.FindBox("static-files/css")
	fontBox, _ := rice.FindBox("static-files/fonts")
	imageBox, _ := rice.FindBox("static-files/images")
	jsBox, _ := rice.FindBox("static-files/js")

	headerTemplate, _ := templateBox.String("header.tmpl")
	footerTemplate, _ := templateBox.String("footer.tmpl")

	indexView, _ := viewBox.String("index.tmpl")
	optionView, _ := viewBox.String("options.tmpl")

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
					DSLAMList := ""
					OptionList := ""
					rows, _ := database.Query("SELECT * FROM DSLAM")
					for rows.Next() {
						var id sql.NullString
						var name sql.NullString
						var address sql.NullString
						rows.Scan(&id, &name, &address,)
						DSLAMList += "<a class=\"list-button pure-button pure-u-1\" href=\"/DSLAM?id="+ id.String +"\"><i class=\"imageDSLAM fa fa-cube fa-3x\"></i><span class=\"textButtonDSLAM\">"+ name.String +"</span></a>"
					}
					rows, _ = database.Query("SELECT * FROM DSLAM WHERE id=?", DSLAMid)
					var id sql.NullString
					var name sql.NullString
					var address sql.NullString
					rows.Next()
					rows.Scan(&id, &name, &address,)
					OptionList = "<form style=\"margin-top:30px\" class=\"pure-form pure-form-aligned\" action=\"/DSLAM?\" method=\"POST\"><fieldset><div class=\"pure-control-group\"><label for=\"name\">Nom</label><input id=\"name\" name=\"name\" type=\"text\" value="+ name.String +"></div><div class=\"pure-control-group\"><label for=\"password\">Adresse</label><input id=\"address\" name=\"address\" type=\"text\" value="+ address.String +"></div><input type=\"hidden\" id=\"id\" value="+ id.String +"><button type=\"submit\" style=\"margin-left:180px\" class=\"pure-button-primary pure-button\">Envoyer</button></fieldset></form>"
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
						database.Exec("UPDATE DSLAM SET name=?, address=? WHERE id=?", r.Form.Get("name"), r.Form.Get("address"), r.Form.Get("id"))
					} else {
						database.Exec("INSERT INTO DSLAM (id, name, address) VALUES (?,?,?)", uuid.New(), r.Form.Get("name"), r.Form.Get("address"))
					}
					tx.Commit()
					//Connexion au DSLAM pour la liste des Cartes et des Ports.
					rows, _ := database.Query("SELECT * FROM DSLAM")
					for rows.Next() {
						var id sql.NullString
						var name sql.NullString
						var address sql.NullString
						rows.Scan(&id, &name, &address,)
						DSLAMList += "<a class=\"list-button pure-button pure-u-1\" href=\"/DSLAM?id="+ id.String +"\"><i class=\"imageDSLAM fa fa-cube fa-3x\"></i><span class=\"textButtonDSLAM\">"+ name.String +"</span></a>"
					}
					response.Execute(w, map[string]string{"DSLAMList": DSLAMList, "Options": ""})
				case "DELETE":
					re := regexp.MustCompile("[a-z0-9\\-]*$")
					DSLAMid := re.FindString(r.URL.RawQuery)
					fmt.Println(DSLAMid)
					test, err := database.Exec("DELETE FROM DSLAM WHERE id=?", DSLAMid)
					txerr := tx.Commit()
					fmt.Println(test)
					fmt.Println(err)
					fmt.Println(txerr)
			}
		case "/option":
			option := headerTemplate
			option += optionView
			option += footerTemplate
			response, _ := template.New("option").Parse(option)
			DSLAMList := ""
			rows, _ := database.Query("SELECT * FROM DSLAM")
			for rows.Next() {
				var id sql.NullString
				var name sql.NullString
				var address sql.NullString
				rows.Scan(&id, &name, &address,)
				DSLAMList += "<a class=\"list-button pure-button pure-u-1\" href=\"/DSLAM?id="+ id.String +"\"><i class=\"imageDSLAM fa fa-cube fa-3x\"></i><span class=\"textButtonDSLAM\">"+ name.String +"</span></a>"
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
			w.Header().Set("Content-Type", "text/javascript")
			w.Write([]byte(response))
		case "/favicon.png":
			response, _ := imageBox.String("favicon.png")
			w.Header().Set("Content-Type", "image/png")
			w.Write([]byte(response))
		default:
	}

}
