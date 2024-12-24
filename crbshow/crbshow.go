package crbshow

import (
	"net/http"
	"text/template"
	"webapp/home/loggeduser"
	"webapp/parsepost"
)

func CrbShow(w http.ResponseWriter, r *http.Request) {
	//parse html
	t, _ := template.ParseFiles("tmpl/createcrbshowcrb.html")
	// send request to parse and get logged user string
	LoggedUser := loggeduser.LoggedUserRun(r)

	// init struct
	Msg := struct {
		Message           string
		MessageLoggedUser string
	}{
		Message:           parsepost.Crbname, //show created cluster role binding
		MessageLoggedUser: LoggedUser,
	}
	// send string to web page
	err := t.Execute(w, Msg)
	if err != nil {
		return
	}
}
