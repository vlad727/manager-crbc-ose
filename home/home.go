package home

import (
	"log"
	"net/http"
	"text/template"
	"webapp/home/loggeduser"
)

// HomeFunc the main page
func HomeFunc(w http.ResponseWriter, r *http.Request) {

	log.Println("Func HomeFunc started")

	// send request to parse and get logged username as string
	LoggedUser := loggeduser.LoggedUserRun(r)

	// parse html
	t, _ := template.ParseFiles("tmpl/getresp.html")

	// create and init struct
	Marketing := struct {
		Message string
	}{
		Message: LoggedUser, // get logged user name
	}

	err := t.Execute(w, Marketing)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")

	// set string to nil
	LoggedUser = ""

}
