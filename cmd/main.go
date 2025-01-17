// Package main register handlers and run listener
package main

import (
	"log"
	"net/http"
	"webapp/crbcmain"
	"webapp/crmatcher/handlers"
	errormsg "webapp/error"
	"webapp/getcrb"
	"webapp/getcrdesc"
	"webapp/getsa"
	"webapp/health"
	"webapp/home"
	"webapp/parsepost"
)

const (
	porthttp = ":8080"
)

func main() {

	// logging
	log.Printf("Starting server on port %s", porthttp)

	// handlers
	registerHandlers()

	// listen http
	http.ListenAndServe(porthttp, nil)
}

func registerHandlers() {
	http.HandleFunc("/", home.HomeFunc)
	http.HandleFunc("/getcrb", getcrb.GetCrb)
	http.HandleFunc("/getsa", getsa.GetSa)
	http.HandleFunc("/crbcmain", crbcmain.CrbcMain)
	http.HandleFunc("/createcrbmanager", parsepost.ParsePostRequest)
	http.HandleFunc("/error", errormsg.ErrorOut)
	http.HandleFunc("/getcrdesc", getcrdesc.GetCrDesc)
	http.HandleFunc("/uploadfile", handlers.UploadFile)
	http.HandleFunc("/uploadedfile", handlers.HandlePost)
	http.HandleFunc("/crmatcherresult", handlers.CrMatcherResult)
	http.HandleFunc("/health", health.Health)
}
