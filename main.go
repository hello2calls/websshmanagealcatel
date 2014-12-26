package main

import (
	"fmt"
	"net/http"
	"bitbucket.org/nmontes/WebSSHManageAlcatel/api"
	"bitbucket.org/nmontes/WebSSHManageAlcatel/web"
	"github.com/skratchdot/open-golang/open"
	)

	func main() {

		fmt.Println("WebServer Listen on all interfaces, port 8080")

		api.Run()
		web.Run()

		open.Start("http://127.0.0.1:8080")

		http.ListenAndServe(":8080", nil)

	}
