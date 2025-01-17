// Package readfiles reads the contents of a file and returns it as a string without newline characters
// func ReadFile run from getsacollect package
package readfiles

import (
	"log"
	"os"
	"strings"
)

func ReadFile() (string, error) {
	// read file with user admin
	fileContent, err := os.ReadFile("/files/user-admin")
	if err != nil {
		log.Printf("Error message: %s", err)
		log.Println("Can't read file ")
	}
	// convert bytes to string
	userAdmin := string(fileContent)
	// logging
	log.Printf("Got username %s", userAdmin)
	userAdminString := strings.ReplaceAll(userAdmin, "\n", "")
	// return string
	return userAdminString, nil

}
