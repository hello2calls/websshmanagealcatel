package dslam

import (
	"net/http"
	"regexp"

	"bitbucket.org/nmontes/WebSSHManageAlcatel/web/pkg/equipment"
	"bitbucket.org/nmontes/WebSSHManageAlcatel/web/pkg/file"
	S "bitbucket.org/nmontes/WebSSHManageAlcatel/web/pkg/structures"
	WD "bitbucket.org/nmontes/WebSSHManageAlcatel/web/pkg/writeData"

	"code.google.com/p/go-uuid/uuid"
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

// Get return DSLAM list page
func Get(w http.ResponseWriter, r *http.Request, dataFile S.Data) {
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
		OptionList += "<br><div class=\"pure-control-group\"><label>Internet</label></div>"
		OptionList += "<div class=\"pure-control-group\"><label for=\"internetVlan\">VLAN</label><input id=\"internetVlan\" name=\"internetVlan\" style=\"width:100px\" type=\"number\" max=\"4094\" placeholder=\"VLAN\" value=" + dataFile.DSLAM[dslamPos].Internet.Vlan + ">"
		OptionList += "<label style=\"width:auto; margin-left:20px\" for=\"internetVpi\">VPI</label><input id=\"internetVpi\" name=\"internetVpi\" style=\"width:100px\" type=\"number\" max=\"255\" placeholder=\"VPI\" value=" + dataFile.DSLAM[dslamPos].Internet.Vpi + ">"
		OptionList += "<label style=\"width:auto; margin-left:20px\" for=\"internetVci\">VCI</label><input id=\"internetVci\" name=\"internetVci\" style=\"width:100px\" type=\"number\" max=\"65535\" placeholder=\"VCI\" value=" + dataFile.DSLAM[dslamPos].Internet.Vci + "></div>"
		OptionList += "<br><div class=\"pure-control-group\"><label>Téléphonie</label></div>"
		OptionList += "<div class=\"pure-control-group\"><label for=\"voipVlan\">VLAN</label><input id=\"voipVlan\" name=\"voipVlan\" style=\"width:100px\" type=\"number\" max=\"4094\" placeholder=\"VLAN\" value=" + dataFile.DSLAM[dslamPos].Telephony.Vlan + ">"
		OptionList += "<label style=\"width:auto; margin-left:20px\" for=\"voipVpi\">VPI</label><input id=\"voipVpi\" name=\"voipVpi\" style=\"width:100px\" type=\"number\" max=\"255\" placeholder=\"VPI\" value=" + dataFile.DSLAM[dslamPos].Telephony.Vpi + ">"
		OptionList += "<label style=\"width:auto; margin-left:20px\" for=\"voipVci\">VCI</label><input id=\"voipVci\" name=\"voipVci\" style=\"width:100px\" type=\"number\" max=\"65535\" placeholder=\"VCI\" value=" + dataFile.DSLAM[dslamPos].Telephony.Vci + "></div>"
		OptionList += "<br><div class=\"pure-control-group\"><label>Vidéo</label></div>"
		OptionList += "<div class=\"pure-control-group\"><label for=\"videoVlan\">VLAN</label><input id=\"videoVlan\" name=\"videoVlan\" style=\"width:100px\" type=\"number\" max=\"4094\" placeholder=\"VLAN\" value=" + dataFile.DSLAM[dslamPos].Video.Vlan + ">"
		OptionList += "<label style=\"width:auto; margin-left:20px\" for=\"videoVpi\">VPI</label><input id=\"videoVpi\" name=\"videoVpi\" style=\"width:100px\" type=\"number\" max=\"255\" placeholder=\"VPI\" value=" + dataFile.DSLAM[dslamPos].Video.Vpi + ">"
		OptionList += "<label style=\"width:auto; margin-left:20px\" for=\"videoVci\">VCI</label><input id=\"videoVci\" name=\"videoVci\" style=\"width:100px\" type=\"number\" max=\"65535\" placeholder=\"VCI\" value=" + dataFile.DSLAM[dslamPos].Video.Vci + "></div>"
		OptionList += "<input type=\"hidden\" name=\"id\" id=\"id\" value=" + dataFile.DSLAM[dslamPos].ID + ">"
		OptionList += "<div class=\"pure-control-group pure-u-1-2\"><label></label><button type=\"submit\" class=\"pure-button-primary pure-button\">Envoyer</button></div>"
		OptionList += "<div class=\"pure-control-group pure-u-1-2\"><label></label><button onclick=\"sendDelete()\" class=\"button-error pure-button\">Supprimer</button></div>"
		OptionList += "</fieldset></form>"
	}
	response.Execute(w, map[string]string{"DSLAMList": DSLAMList, "Options": OptionList})
}

// Post update or create new DSLAM when we are in option page
func Post(w http.ResponseWriter, r *http.Request, dataFile S.Data) {
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
	newDSLAM.Internet.Vlan = r.Form.Get("internetVlan")
	newDSLAM.Internet.Vpi = r.Form.Get("internetVpi")
	newDSLAM.Internet.Vci = r.Form.Get("internetVci")
	newDSLAM.Telephony.Vlan = r.Form.Get("voipVlan")
	newDSLAM.Telephony.Vpi = r.Form.Get("voipVpi")
	newDSLAM.Telephony.Vci = r.Form.Get("voipVci")
	newDSLAM.Video.Vlan = r.Form.Get("videoVlan")
	newDSLAM.Video.Vpi = r.Form.Get("videoVpi")
	newDSLAM.Video.Vci = r.Form.Get("videoVci")
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
}

// Delete remove DSLAM when we are in option page
func Delete(w http.ResponseWriter, r *http.Request, dataFile S.Data) {
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
