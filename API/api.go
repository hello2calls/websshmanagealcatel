package api

import (
	"fmt"
	"net/http"
	"bitbucket.org/nmontes/WebSSHManageAlcatel/ssh"
	)

	func Run() {

		fmt.Println("API Lunched")

		http.HandleFunc("/API/session", sshConnect.SessionHandler)
		http.HandleFunc("/API/command", sshConnect.CommandHandler)

	}
