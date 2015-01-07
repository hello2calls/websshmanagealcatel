package web

import (
	"fmt"
	"net/http"
	"text/template"
	"github.com/GeertJohan/go.rice"
	"database/sql"
	"github.com/mattn/go-sqlite3"
)

type DSLAM struct {
	id		sql.NullInt64
	name	sql.NullString
	IP		sql.NullString
}

type Card struct {
	id				sql.NullInt64
	name			sql.NullString
	DSLAM_id	sql.NullInt64
}

type Port struct {
	id			sql.NullInt64
	name		sql.NullString
	Card_id	sql.NullInt64
}

func Run() {

	var DB_DRIVER string
	sql.Register(DB_DRIVER, &sqlite3.SQLiteDriver{})
	database, err := sql.Open(DB_DRIVER, "mysqlite_3")
	if err != nil {
		fmt.Println("--Failed to create the handle")
	}
	if err2 := database.Ping(); err2 != nil {
		fmt.Println("Failed to keep connection alive")
	}

	fmt.Println("WebSite Lunched")

	http.HandleFunc("/", indexHandler)

}


func indexHandler(w http.ResponseWriter, r *http.Request) {

	templateBox, _ := rice.FindBox("templates")
	viewBox, _ := rice.FindBox("views")
	cssBox, _ := rice.FindBox("static-files/css")
	fontBox, _ := rice.FindBox("static-files/fonts")
	imageBox, _ := rice.FindBox("static-files/images")

	switch r.URL.Path {

		case "/":
			headerTemplate, _ := templateBox.String("header.tmpl")
			footerTemplate, _ := templateBox.String("footer.tmpl")
			indexView, _ := viewBox.String("index.tmpl")
			index := headerTemplate
			index += indexView
			index += footerTemplate
			response, _ := template.New("index").Parse(index)
			DSLAMList := "<button class=\"pure-button pure-u-1\"><i class=\"imageDSLAM fa fa-cube fa-3x\"></i><span class=\"textButtonDSLAM\">DSLAM 2</span></button>"
			CardList := "<button class=\"pure-button pure-u-1\"><span class=\"fa-stack fa-lg\"><i class=\"imageCard fa fa-hdd-o fa-stack-2x\"></i><i class=\"imageCard fa fa-exclamation fa-stack-2x\"></i></span><span class=\"textButtonCard\">NVLT-N</span></button>"
			PortList := "<button class=\"pure-button pure-u-1\"><span class=\"fa-stack fa-lg\"><i class=\"imagePort fa fa-caret-square-o-right\"></i></span><span class=\"textButtonPort\">Chambre 1</span></button>"
			response.Execute(w, map[string]string{"DSLAMList": DSLAMList, "CardList": CardList, "PortList": PortList})
		case "/option":
			headerTemplate, _ := templateBox.String("header.tmpl")
			footerTemplate, _ := templateBox.String("footer.tmpl")
			optionView, _ := viewBox.String("options.tmpl")
			option := headerTemplate
			option += optionView
			option += footerTemplate
			response, _ := template.New("option").Parse(option)
			DSLAMList := "<button class=\"pure-button pure-u-1\"><i class=\"imageDSLAM fa fa-cube fa-3x\"></i><span class=\"textButtonDSLAM\">DSLAM 2</span></button>"
			Options := ""
			if r.URL.RawQuery == "add" {
				Options = "ADD"
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
		case "/favicon.png":
			response, _ := imageBox.String("favicon.png")
			w.Header().Set("Content-Type", "image/png")
			w.Write([]byte(response))
		default:
	}

}
