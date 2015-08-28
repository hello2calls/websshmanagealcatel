package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

type TestURL struct {
	url    string
	method string
}

var testsDslamID = []string{
	"4883886d-e8e4-4eb8-8df7-c147e0504571",
	"0",
	"0azerty",
	"azerty",
	"shf67!('é&&dfbè-è!é$",
}

var testsURL = []TestURL{
	{"http://127.0.0.1:8080/", "GET"},
	{"http://127.0.0.1:8080/DSLAM", "GET"},
	{"http://127.0.0.1:8080/DSLAM?id=4883886d-e8e4-4eb8-8df7-c147e0504571", "GET"},
	{"http://127.0.0.1:8080/DSLAM", "POST"},
	{"http://127.0.0.1:8080/DSLAM", "DELETE"},
	{"http://127.0.0.1:8080/option", "GET"},
	{"http://127.0.0.1:8080/SITEAPI/all", "GET"},
	{"http://127.0.0.1:8080/SITEAPI/update", "GET"},
	{"http://127.0.0.1:8080/SITEAPI/services", "GET"},
	{"http://127.0.0.1:8080/SITEAPI/command", "GET"},
	{"http://127.0.0.1:8080/css/pure-min.css", "GET"},
	{"http://127.0.0.1:8080/css/grids-responsive-min.css", "GET"},
	{"http://127.0.0.1:8080/css/grids-responsive-old-ie-min.css", "GET"},
	{"http://127.0.0.1:8080/css/WebManageAlcatel.css", "GET"},
	{"http://127.0.0.1:8080/css/WebManageAlcatel-old-ie.css", "GET"},
	{"http://127.0.0.1:8080/css/font-awesome.min.css", "GET"},
	{"http://127.0.0.1:8080/fonts/fontawesome-webfont.woff", "GET"},
	{"http://127.0.0.1:8080/fonts/fontawesome-webfont.ttf", "GET"},
	{"http://127.0.0.1:8080/fonts/fontawesome-webfont.eot", "GET"},
	{"http://127.0.0.1:8080/js/WebManageAlcatel.js", "GET"},
	{"http://127.0.0.1:8080/favicon.png", "GET"},
	{"http://127.0.0.1:8080/background.jpg", "GET"},
	{"http://127.0.0.1:8080/mitel.png", "GET"},
}

var dataFile Data

// Read File
func readFile() {
	var file, _ = ioutil.ReadFile("data.json")
	var err = json.Unmarshal(file, &dataFile)
	if err != nil {
		fmt.Println("JSON Unmarshal file in indexHandler error:", err)
	}
}

// Test function run
func TestRun(t *testing.T) {
	Run()
	go http.ListenAndServe("127.0.0.1:8080", nil)
	res, _ := http.Get("http://127.0.0.1:8080/")
	switch res.Status {
	default:
		t.Errorf("#ERROR : %s", res.Status)
	case "200 OK":
	}
}

// Test function indexHandler
func TestIndexHandler(t *testing.T) {
	for i, test := range testsURL {
		req, _ := http.NewRequest(test.method, test.url, bytes.NewBufferString(""))
		client := &http.Client{}
		res, _ := client.Do(req)
		switch res.Status {
		default:
			t.Errorf("#ERROR : %d : for url : %s : %s", i, test.url, res.Status)
		case "200 OK":
		}
	}
}

// Test function getDslamById
func TestGetDslamByID(t *testing.T) {
	readFile()
	for i, test := range testsDslamID {
		var dslam interface{}
		dslam = getDslamByID(dataFile, test)
		switch dslam := dslam.(type) {
		default:
			t.Errorf("#ERROR : %d : DSLAM(%s)=%s", i, test, dslam)
		case DSLAM:
		}
	}
}

// Test function getDslamPosById
func TestGetDslamPosByID(t *testing.T) {
	readFile()
	for i, test := range testsDslamID {
		var dslam interface{}
		dslam = getDslamPosByID(dataFile, test)
		switch dslam := dslam.(type) {
		default:
			t.Errorf("#ERROR : %d : DSLAM(%s)=%s", i, test, dslam)
		case int:
		}
	}
}

// Test function getDslamPosById
func TestGetSSHSession(t *testing.T) {
	readFile()
	for i, test := range testsDslamID {
		var dslam interface{}
		dslam = getSSHSession(dataFile, test)
		switch dslam := dslam.(type) {
		default:
			t.Errorf("#ERROR : %d : DSLAM(%s)=%s", i, test, dslam)
		case string:
		}
	}
}
