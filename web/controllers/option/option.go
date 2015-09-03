package option

import (
	"net/http"

	S "bitbucket.org/nmontes/WebSSHManageAlcatel/web/pkg/structures"

	"github.com/GeertJohan/go.rice"
	"github.com/golang/go/src/pkg/text/template"
)

var templateBox, _ = rice.FindBox("../../templates")
var viewBox, _ = rice.FindBox("../../views")
var cssBox, _ = rice.FindBox("../../static-files/css")
var fontBox, _ = rice.FindBox("../../static-files/fonts")
var imageBox, _ = rice.FindBox("../../static-files/images")
var jsBox, _ = rice.FindBox("../../static-files/js")
var headerTemplate, _ = templateBox.String("header.tmpl")
var footerTemplate, _ = templateBox.String("footer.tmpl")
var indexView, _ = viewBox.String("index.tmpl")
var optionView, _ = viewBox.String("options.tmpl")

// Get return option page
func Get(w http.ResponseWriter, r *http.Request, dataFile S.Data) {
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
}
