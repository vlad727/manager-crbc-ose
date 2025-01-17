package getcrb

import (
	"fmt"
	"golang.org/x/net/context"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/strings/slices"
	"log"
	"net/http"
	"os"
	"sigs.k8s.io/yaml"
	"strings"
	"text/template"
	"webapp/clientgo"
	"webapp/loggeduser"
)

// GetCrb execute after press button "Get Cluster Role Binding"
// read file with allowed cluster role bindings
// provide list of cluster role bindings which one allowed to show
func GetCrb(w http.ResponseWriter, r *http.Request) {
	log.Println("Func GetCrb started")

	// send request to parse, trim and decode jwt, get map with user and groups
	UserAndGroups := loggeduser.LoggedUserRun(r)

	username := ""                 // name of logged user
	for k := range UserAndGroups { // get logged user name from map
		username = k
		break
	}

	// Reads file with cluster role bindings which should hide
	clusterRoleFromFile, err := os.ReadFile("/files/clusterroles")
	if err != nil {
		log.Printf("Error message: %v", err)
	}

	// convert bytes to string
	clusterRoleFromFileString := string(clusterRoleFromFile)

	// split string and put it to slice
	listCrbNotAllowed := strings.Split(clusterRoleFromFileString, "\n")

	// list cluster role binding
	listCRB, err := clientgo.Сlientset.RbacV1().ClusterRoleBindings().List(context.Background(), v1.ListOptions{})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Cluster is unavailable: %s", err)
		return
	}
	// empty slice for append cluster role bindings
	var listOfCRB []string

	// map for appending slice of strings with names
	mapListCrb := map[string][]string{}

	// iterate over items to get name for cluster role binding and linked cluster role
	for _, el := range listCRB.Items {
		if slices.Contains(listCrbNotAllowed, el.RoleRef.Name) {
			//log.Println("Not allowed to show ")
		} else {
			listOfCRB = append(listOfCRB, "<b>"+el.Name+"</b>"+" "+el.RoleRef.Name)
			mapListCrb["List"] = listOfCRB
		}
	}
	// logging
	log.Println("Iteration over cluster role bindings finished")

	// Marshal to yaml for out to web page
	yamlFile, err := yaml.Marshal(mapListCrb)
	if err != nil {
		log.Printf("Failed to marshal YAML: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to process the request: %v", err)
		return
	}

	// convert to string for struct if you do not convert it will be in bytes
	str := string(yamlFile)

	//parse html
	t, err := template.ParseFiles("tmpl/getcrb.html")
	if err != nil {
		log.Fatalf("Failed to parse template: %v", err)
	}
	// init struct
	Msg := struct {
		Message           string `yaml:"message"`
		MessageLoggedUser string
	}{
		Message:           str,
		MessageLoggedUser: username,
	}
	// send string to web page execute
	err = t.Execute(w, Msg)
	if err != nil {
		return
	}
	// set slice to nil to prevent add new items after page refresh
	listOfCRB = nil
}
