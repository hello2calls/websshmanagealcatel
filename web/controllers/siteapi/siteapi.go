package siteapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"bitbucket.org/nmontes/WebSSHManageAlcatel/web/pkg/command"
	"bitbucket.org/nmontes/WebSSHManageAlcatel/web/pkg/equipment"
	"bitbucket.org/nmontes/WebSSHManageAlcatel/web/pkg/sshsession"
	S "bitbucket.org/nmontes/WebSSHManageAlcatel/web/pkg/structures"
	WD "bitbucket.org/nmontes/WebSSHManageAlcatel/web/pkg/writeData"

	"github.com/GeertJohan/go.rice"
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

// All return all DSLAM data
func All(w http.ResponseWriter, r *http.Request, dataFile S.Data) {
	w.Header().Set("Content-Type", "application/json")
	data, _ := json.MarshalIndent(dataFile, "", "\t")
	w.Write(data)
}

// Update data file on server (with real information -> SSH)
func Update(w http.ResponseWriter, r *http.Request, dataFile S.Data) {
	w.Header().Set("Content-Type", "application/json")
	i := 0
	for i = 0; i < len(dataFile.DSLAM); i++ {
		var dslamID = dataFile.DSLAM[i].ID
		fmt.Println("Update DSLAM " + dslamID)
		WD.WriteStatus(dataFile, dslamID)
		if dataFile.DSLAM[i].Status == "OK" {
			sessionID := sshsession.Get(dataFile, dslamID)
			if sessionID != "SSH_KO" {
				WD.WriteCard(dataFile, sessionID, dslamID)
			}
		}
	}
	fmt.Println("Data File Updated")
	w.Write([]byte("OK"))
}

// Services update services on DSLAM with post informations (interface name, internet, voip, iptv)
func Services(w http.ResponseWriter, r *http.Request, dataFile S.Data) {
	if r.URL.RawQuery != "" {
		params, _ := url.ParseQuery(r.URL.RawQuery)
		name := params.Get("portName")
		internet := params.Get("internetSwitch")
		voip := params.Get("voipSwitch")
		iptv := params.Get("iptvSwitch")
		dslamID := params.Get("dslamID")
		slot := params.Get("slot")
		portIndex := params.Get("portIndex")
		sessionID := sshsession.Get(dataFile, dslamID)
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
			command.Set(sessionID, "configure bridge port "+portIndex+":8:35")
			command.Set(sessionID, "configure bridge port "+portIndex+":8:35 vlan-id 10")
			command.Set(sessionID, "configure bridge port "+portIndex+":8:35 pvid 10")
			WD.WriteServiceOnePort(dataFile, sessionID, dslamID, portIndex)
		} else if oldInternet == false {
		} else if internet == "true" && oldInternet == true {
		} else if oldInternet == true {
			command.Set(sessionID, "configure bridge no port "+portIndex+":8:35")
			WD.WriteServiceOnePort(dataFile, sessionID, dslamID, portIndex)
		}
		// Update voip service
		if voip == "true" && oldVoip == false {
			command.Set(sessionID, "configure bridge port "+portIndex+":8:36")
			command.Set(sessionID, "configure bridge port "+portIndex+":8:36 vlan-id 20")
			command.Set(sessionID, "configure bridge port "+portIndex+":8:36 pvid 20")
			WD.WriteServiceOnePort(dataFile, sessionID, dslamID, portIndex)
		} else if oldVoip == false {
		} else if voip == "true" && oldVoip == true {
		} else if oldVoip == true {
			command.Set(sessionID, "configure bridge no port "+portIndex+":8:36")
			WD.WriteServiceOnePort(dataFile, sessionID, dslamID, portIndex)
		}
		// Update iptv service
		if iptv == "true" && oldIptv == false {
			command.Set(sessionID, "configure bridge port "+portIndex+":8:37")
			command.Set(sessionID, "configure bridge port "+portIndex+":8:37 vlan-id 30")
			command.Set(sessionID, "configure bridge port "+portIndex+":8:37 pvid 30")
			WD.WriteServiceOnePort(dataFile, sessionID, dslamID, portIndex)
		} else if oldIptv == false {
		} else if iptv == "true" && oldIptv == true {
		} else if oldIptv == true {
			command.Set(sessionID, "configure bridge no port "+portIndex+":8:37")
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
}

// Command send commande to DSLAM
func Command(w http.ResponseWriter, r *http.Request, dataFile S.Data) {
	if r.URL.RawQuery != "" {
		reSession := regexp.MustCompile("sessionID=[a-z0-9\\-]*")
		session := reSession.FindString(r.URL.RawQuery)
		resessionID := regexp.MustCompile("[a-z0-9\\-]*$")
		sessionID := resessionID.FindString(session)
		reCommand := regexp.MustCompile("command=[a-zA-Z0-9\\/\\-%]*$")
		commandURL := reCommand.FindString(r.URL.RawQuery)
		reCommandRaw := regexp.MustCompile("[a-zA-Z0-9\\/\\-%]*$")
		commandRaw := reCommandRaw.FindString(commandURL)
		commandRawReplace := strings.Replace(commandRaw, "%20", " ", -1)
		commandOut := command.GetOut(sessionID, commandRawReplace)
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
}
