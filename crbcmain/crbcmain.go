// Package crbcmain the main page for tab "Create Cluster Role Binding"
// read file with allowed label, collect allowed cluster roles, provides select element with service account and namespace
// also provide select element with cluster role for description cluster role and counter (how much crb has been created)
package crbcmain

import (
	"fmt"
	"golang.org/x/net/context"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"net/http"
	"os"
	"text/template"
	"webapp/clientgo"
	"webapp/counter"
	"webapp/getsacollect"
	"webapp/loggeduser"
)

type DataStruct struct {
	CrbSlice          []string
	SaMap             []string
	MessageLoggedUser string
	NumberOfEntities  int
}

func CrbcMain(w http.ResponseWriter, r *http.Request) {

	// logging
	log.Println("Func CrbcMain started")

	// send request to parse, trim and decode jwt, get map with user and groups
	UserAndGroups := loggeduser.LoggedUserRun(r)

	var username string
	// get logged user name from map
	for k := range UserAndGroups {
		username = k
		break
	}
	// func Counter collect and count all cluster role bindings which one has substring "crbc"
	numberOfEntities := counter.Counter()

	// empty slice for service account and name namespace
	var sliceSaName []string

	// get slice with map and change it
	_, slSaAndNs := getsacollect.GetSaCollect(UserAndGroups)

	// change slice add colon and space to slice
	for _, x := range slSaAndNs {
		for k, v := range x {
			result := fmt.Sprint(k + ":" + " " + v)
			sliceSaName = append(sliceSaName, result)
		}
	}

	// read file with label
	labelCrBytes, err := os.ReadFile("/files/allowedlabel")
	if err != nil {
		log.Fatalf("Can't read file: %v", err)
	}

	// convert bytes to string
	labelString := string(labelCrBytes)

	// list cluster role binding only with allowed label
	listClusterRoles, err := clientgo.Сlientset.RbacV1().ClusterRoles().List(context.Background(), v1.ListOptions{LabelSelector: labelString})
	if err != nil {
		log.Println(err)
	}

	// slice for allowed Cluster Roles
	var sliceCrAllowed []string

	// iterate over items to get name cluster role will use for creating cluster role binding
	for _, el := range listClusterRoles.Items {
		sliceCrAllowed = append(sliceCrAllowed, el.Name)
	}
	// logging cluster roles
	log.Println("Slice cluster roles requested and collected")

	// struct for output
	DataProvider := DataStruct{
		CrbSlice:          sliceCrAllowed,   // output slice with cluster roles with allowed label
		SaMap:             sliceSaName,      // output map example: my-sa: my-namespace
		MessageLoggedUser: username,         // logged user string
		NumberOfEntities:  numberOfEntities, // number of cluster role bindings created via manager-crbc
	}

	// parse template
	t, err := template.ParseFiles("tmpl/crbcmain.html")
	if err != nil {
		log.Printf("Can't parse file: %v", err)
	}

	err = t.Execute(w, DataProvider)
	if err != nil {
		log.Printf("Error execute %v", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	log.Println("End of crbcmain func...")

}
