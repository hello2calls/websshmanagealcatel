package api

import (
	"fmt"
	"net/http"
	"bitbucket.org/nmontes/WebSSHManageAlcatel/ssh"
	)

	func Run() {

		fmt.Println("API Lunched")

		// RUN SessionHandler when client ask session
		http.HandleFunc("/API/session", sshConnect.SessionHandler)
		// RUN CommandHandler when client send command
		http.HandleFunc("/API/command", sshConnect.CommandHandler)

	}
