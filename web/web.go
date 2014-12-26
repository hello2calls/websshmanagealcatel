package web

import (
	"fmt"
	"net/http"
	)

func Run() {

	fmt.Println("WebSite Lunched")

	http.HandleFunc("/", indexHandler)

}


func indexHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case "GET":
		fmt.Println("GET /")
		w.Write([]byte("GET /"))
		case "POST":
		w.Write([]byte("Not Implemented"))
		case "PUT":
		w.Write([]byte("Not Implemented"))
		case "DELETE":
		w.Write([]byte("Not Implemented"))
	}
}
