package loggeduser

import (
	"log"
	"net/http"
)

func LoggedUserRun(r *http.Request) string {
	log.Println("Func LoggedUserRun started... ")
	var LoggedUser string // temporary var for user name
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
				LoggedUser = y

			}
		}
	}
	return LoggedUser
}
