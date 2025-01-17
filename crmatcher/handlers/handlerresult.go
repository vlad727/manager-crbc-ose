// Package handlers func CrMatcherResult compare len for all cluster roles and uploaded cluster role
package handlers

import (
	"log"
	"net/http"
	"os"
	"text/template"
	"webapp/crmatcher/getcrname"
	"webapp/crmatcher/readfile"
	"webapp/loggeduser"
)

func CrMatcherResult(w http.ResponseWriter, r *http.Request) {

	// logging
	log.Println("Fund CrMatcherResult started")

	// send request to parse, trim and decode jwt, get map with user and groups
	UserAndGroups := loggeduser.LoggedUserRun(r)

	var LoggedUser string // username for logged user
	for k, v := range UserAndGroups {
		k = LoggedUser
		log.Println(k, v)
	}

	// send dst dir to read file
	readfile.ReadFileYaml(DstDirName) // Note this variable is visible for all func because it's global var and another code in the same package, declared in handlerfile.go

	// compare cluster roles with provided cluster role from yaml
	var ResultForCheck string

	for k, v := range getcrname.CrAllowedList() {
		if v == readfile.LenForCr { // если находим одинаковую длину то у нас match таким образом мы находим совпадение или не совпадение
			log.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
			log.Printf("Looks like %s cluster role is the same as %s ", k, readfile.Cr.Metadata.Name)
			log.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
			ResultForCheck = "Looks like your cluster role " + "<b>" + k + "</b>" + " the same as " + "<b>" + readfile.Cr.Metadata.Name + "</b>" + ". You don't need add this cluster role"
			break
		} else {
			//log.Printf("Branch else get values and keys, value: %s key: %d", k, v)
			ResultForCheck = "Such cluster role does not exist. You may ask administrator to add your cluster role for this cluster"

		}
	}

	//parse html
	t, _ := template.ParseFiles("tmpl/crmatcherresult.html")

	// init struct
	Msg := struct {
		MessageLoggedUser string
		MessageResult     string
	}{
		MessageLoggedUser: LoggedUser, //home.LoggedUser,
		MessageResult:     ResultForCheck,
	}
	// send string to web page execute
	err := t.Execute(w, Msg)
	if err != nil {
		return
	}

	// set string to nill
	ResultForCheck = ""
	// set len for cluster role to 0
	readfile.LenForCr = 0

	// clear dir with uploaded files
	pathString := "/app/uploads"
	err = os.RemoveAll(pathString)
	if err != nil {
		log.Printf("Can't remove files from dir %s", pathString)
		log.Println(err)
	} else {
		log.Println("Dir uploads has been removed")
	}
}
