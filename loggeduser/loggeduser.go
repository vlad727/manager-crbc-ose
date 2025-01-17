package loggeduser

import (
	"encoding/json"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"slices"
	"webapp/clientgo"
)

type GroupStruct struct {
	Items []struct {
		Metadata struct {
			Name string `json:"name"`
		} `json:"metadata"`
		Users []string `json:"users"`
	} `json:"items"`
}

// LoggedUserRun parse request and return map with user and them groups
func LoggedUserRun(r *http.Request) map[string][]string {

	log.Println("Func LoggedUserRun started... ")
	var loggedUser string // temporary var for user name
	r.ParseForm()         // Анализирует переданные параметры url, затем анализирует пакет ответа для тела POST (тела запроса)
	// внимание: без вызова метода ParseForm последующие данные не будут получены
	log.Println(r.Header)
	log.Println(r)
	// Loop over header names
	for name, values := range r.Header {
		//log.Println(name, values)
		if name == "X-Forwarded-User" {
			log.Println(values)
			log.Printf("Got username %s", values)
			for _, y := range values {
				loggedUser = y

			}
		}
	}
	slGroups := GroupCollect(loggedUser)
	mUserGroups := make(map[string][]string)
	mUserGroups[loggedUser] = slGroups

	return mUserGroups

}

func GroupCollect(LoggedUser string) []string {

	// get groups with RestClient
	listgroups, err := clientgo.Сlientset.AppsV1().RESTClient().Get().AbsPath("/apis/user.openshift.io/v1/groups").DoRaw(context.TODO())
	if err != nil {
		log.Printf("Failed %s", listgroups)
		log.Println(err)
	}

	// init struct
	dataObjet := GroupStruct{}

	jsonErr := json.Unmarshal(listgroups, &dataObjet)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	log.Println(dataObjet.Items)

	var sliceGroupForUser []string

	for _, x := range dataObjet.Items {
		log.Println(x.Metadata.Name) // list group name
		log.Println(x.Users)         // list of slice users
		if slices.Contains(x.Users, LoggedUser) {
			sliceGroupForUser = append(sliceGroupForUser, x.Metadata.Name)
		}
	}

	// logged user and them groups

	log.Printf("Collectted data for user: %s", sliceGroupForUser)
	return sliceGroupForUser
}
