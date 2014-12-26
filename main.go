package main

import (
	"fmt"
	"net/http"
	"bitbucket.org/nmontes/SSHManage/ssh"
)

	func main() {

		fmt.Println("WebServer Listen on all interfaces, port 8080")

		http.HandleFunc("/", sshConnect.IndexHandler)
		http.HandleFunc("/API/session", sshConnect.SessionHandler)
		http.HandleFunc("/API/command", sshConnect.CommandHandler)
		http.ListenAndServe(":8080", nil)

	}
