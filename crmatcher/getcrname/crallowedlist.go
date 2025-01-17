// Package getcrname get list all cluster role names
// compare it with allowed list from configmap
package getcrname

import (
	"golang.org/x/net/context"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/strings/slices"
	"log"
	"os"
	"strings"
	"time"
	"webapp/clientgo"
	"webapp/crmatcher/getlen"
)

// func CrAllowedList run from func CrMatcherResult package handlers file handlerresult.go
// should return map like below
// map[app-shard-resources-role:253 application-configurator-admin-role:395 application-configurator-editor-role:334
// application-configurator-leader-election-role:57
func CrAllowedList() map[string]int {
	// execution time
	start := time.Now()

	// logging
	log.Println("Func CrAllowedList started")

	// get all cluster roles from kubernetes API
	listAllCrs := GetCrNameList()

	// read file with not allowed cluster roles
	listForbiddenCrs := ReadFileCrNames()

	slAllowed := []string{}

	// iterate over slice with all cluster role names
	for _, x := range listAllCrs {
		if !slices.Contains(listForbiddenCrs, x) { // if x not in slice forbiddenCr add it to allowed slice
			slAllowed = append(slAllowed, x) // we will get all allowed cluster roles
		}

	}

	// map contain cluster role name and len for it
	mapCRLen := getlen.GetLen(slAllowed) // <<<<<<<<<<<  To change!!!

	//log.Println(mapCRLen) // the end for crallowedlist !!!

	// Code to measure
	duration := time.Since(start)
	log.Printf("Time execution for func CrAllowedList is %s", duration)
	return mapCRLen
}

// GetCrNameList collect cluster role names and return it to crcheck.
func GetCrNameList() []string {
	start := time.Now()
	log.Println("Func GetCrNameList started ")
	listCr, err := clientgo.Сlientset.RbacV1().ClusterRoles().List(context.TODO(), v1.ListOptions{})
	if err != nil {
		log.Printf("Failed %s", listCr)
		log.Println(err)
	}

	slClusterRoles := []string{}

	// iterate over cluster roles and append it to slice
	for _, cr := range listCr.Items {
		//log.Println(cr.Name)
		slClusterRoles = append(slClusterRoles, cr.Name)
	}
	// Code to measure
	duration := time.Since(start)
	log.Printf("Time execution for Func GetCrNameList  %s", duration)
	return slClusterRoles
}

// Func ReadFileCrNames read file clusterroles which one not allowed to show
func ReadFileCrNames() []string {

	start := time.Now()
	// logging readFile
	log.Println("Func ReadFileCrNames started")

	// read file with user admin
	crNames, err := os.ReadFile("/files/clusterroles")
	if err != nil {
		log.Printf("Error message: %s", err)
		log.Println("Can't read file ")

	}
	strCrNames := string(crNames)

	slClusterRolesNotAllowed := []string{}

	s := strings.Split(strCrNames, "\n")
	for _, x := range s {
		slClusterRolesNotAllowed = append(slClusterRolesNotAllowed, x)
	}
	// Code to measure
	duration := time.Since(start)
	log.Printf("Time execution for Func ReadFileCrNames  %s", duration)

	return slClusterRolesNotAllowed
}
