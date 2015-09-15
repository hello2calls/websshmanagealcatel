package main

import (
	"fmt"
	"net/http"

	"bitbucket.org/nmontes/WebSSHManageAlcatel/web"
	"bitbucket.org/nmontes/WebSSHManageAlcatel/web/pkg/logger"
	"bitbucket.org/nmontes/sshapi/api"
	"github.com/skratchdot/open-golang/open"
)

func main() {

	fmt.Println("WebServer Listen on all interfaces, port 8080")
	logger.Print("WebServer Listen on all interfaces, port 8080", nil)

	// RUN API
	SSHAPI.Run()
	// RUN WebSite Server
	web.Run()

	// OPEN Browser
	open.Start("http://127.0.0.1:8080")

	// RUN WebServer
	http.ListenAndServe("127.0.0.1:8080", nil)

}
