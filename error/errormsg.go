package errormsg

import (
	"net/http"
	"text/template"
	"webapp/home/loggeduser"
	"webapp/parsepost"
)

func ErrorOut(w http.ResponseWriter, r *http.Request) {

	// send request to parse and get logged user string
	LoggedUser := loggeduser.LoggedUserRun(r)
	t, _ := template.ParseFiles("tmpl/error.html")
	// init struct
	Msg := struct {
		Message           string
		MessageLoggedUser string
	}{
		Message:           parsepost.ErrorMsg,
		MessageLoggedUser: LoggedUser,
	}
	// send string to web page
	err := t.Execute(w, Msg)
	if err != nil {
		return
	}

}
