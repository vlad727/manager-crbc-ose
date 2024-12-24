// Package handlers func UplodaFile show page with button upload and allow you get file from your local machine
package handlers

import (
	"log"
	"net/http"
	"text/template"
	"time"
	"webapp/home/loggeduser"
)

func UploadFile(w http.ResponseWriter, r *http.Request) {

	// send request to parse and get logged user string
	LoggedUser := loggeduser.LoggedUserRun(r)
	// execution time
	start := time.Now()
	//logging
	log.Println("Func UploadFile started ")
	//parse html
	t, _ := template.ParseFiles("tmpl/crmatcher.html")

	// init struct
	Msg := struct {
		MessageLoggedUser string
	}{
		MessageLoggedUser: LoggedUser, //home.LoggedUser,
	}
	// send string to web page execute
	err := t.Execute(w, Msg)
	if err != nil {
		return
	}

	// Code to measure
	duration := time.Since(start)
	log.Printf("Time execution for func UploadFile is %s", duration)
}
