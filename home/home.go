package home

import (
	"log"
	"net/http"
	"text/template"
	"webapp/loggeduser"
)

// HomeFunc the main page, get request, send to parse, get logged user, parse html
func HomeFunc(w http.ResponseWriter, r *http.Request) {

	// send request to parse, trim and decode jwt, get map with user and groups
	UserAndGroups := loggeduser.LoggedUserRun(r)

	// empty var for name of logged user
	username := ""
	// get logged user name from map and skip groups
	for k := range UserAndGroups {
		username = k
		break // we do need to proceed only get username
	}

	t, err := template.ParseFiles("tmpl/getresp.html")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// create and init struct
	userStruct := struct {
		Message string
	}{
		Message: username, // set logged user login and put to html
	}

	err = t.Execute(w, userStruct)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// send
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}
