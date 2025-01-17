// Package getcrdesc get request with cluster role name then parse it
// iterate over it, and provide yaml to web page
package getcrdesc

import (
	"golang.org/x/net/context"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"net/http"
	"sigs.k8s.io/yaml"
	"text/template"
	"webapp/clientgo"
	"webapp/loggeduser"
)

func GetCrDesc(w http.ResponseWriter, r *http.Request) {

	// logging
	log.Println("Func GetCrDes started...")

	// send request to parse, trim and decode jwt, get map with user and groups
	UserAndGroups := loggeduser.LoggedUserRun(r)

	var username string            // name of logged user
	for k := range UserAndGroups { // get logged user name from map
		username = k
		break
	}

	// parse post request
	err := r.ParseForm()
	if err != nil {
		log.Println("Can't parse request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// logging
	log.Println(r.Form)
	log.Println("Path: ", r.URL.Path)

	// cluster role name
	var clusterRoleName string

	// iterate over request with for to get requested cluster role name
	// example post request: map[choice1:[1nd-line-support]]
	for k, v := range r.Form {
		log.Printf("key: %s value: %s", k, v)
		for _, el := range v {
			clusterRoleName = el
		}

	}

	// get cluster role from k8s
	getCr, err := clientgo.Сlientset.RbacV1().ClusterRoles().Get(context.Background(), clusterRoleName, v1.GetOptions{})
	if err != nil {
		log.Println(err)
		log.Printf("Cluster is unavailable %s", err)
	}

	// main map
	var outSlice []map[string][]string

	// temp slices and temp maps
	var sl1, sl2, sl3, sl4 []string

	m0 := map[string][]string{}
	m1 := map[string][]string{}
	m2 := map[string][]string{}
	m3 := map[string][]string{}
	m4 := map[string][]string{}
	
	// iterate over cluster role
	for _, el := range getCr.Rules {

		tempslice := [][]string{el.APIGroups, el.ResourceNames, el.Resources, el.Verbs, el.NonResourceURLs}
		for x, item := range tempslice {

			switch x {
			case 0:
				sl0 := []string{}
				for _, a := range item {
					if len(a) == 0 {
						a += "\"\""
					}
					sl0 = append(sl0, a)

				}
				// add data to map
				m0["apiGroups"] = sl0
				// add data to main slice
				outSlice = append(outSlice, m0)

				// clear map
				m0 = make(map[string][]string)
				// clear slice
				sl0 = nil

			case 1:
				if len(item) == 0 {
					log.Println("No resourceNames")
				} else {
					for _, a := range item {

						sl1 = append(sl1, a)

					}
					m1["resourceNames"] = sl1
					outSlice = append(outSlice, m1)

					// clear map
					m1 = make(map[string][]string)
					// clear slice
					sl1 = nil
				}
			case 2:
				for _, a := range item {

					sl2 = append(sl2, a)

				}
				m2["resources"] = sl2
				outSlice = append(outSlice, m2)

				// clear map
				m2 = make(map[string][]string)
				// clear slice
				sl2 = nil
			case 3:
				for _, a := range item {

					sl3 = append(sl3, a)

				}
				m3["verbs"] = sl3
				outSlice = append(outSlice, m3)
				// clear map
				m3 = make(map[string][]string)
				// clear slice
				sl3 = nil
			case 4:

				if len(item) == 0 {
					log.Println("No nonResourceURLs")
				} else {
					for _, a := range item {

						sl4 = append(sl4, a)

					}
					m4["nonResourceURLs"] = sl4
					outSlice = append(outSlice, m4)
					// clear map
					m4 = make(map[string][]string)
					// clear slice
					sl4 = nil
				}

			}

		}

	}
	// logging slice
	log.Println(outSlice)

	// Marshal to yaml for out to web page
	yamlFile, err := yaml.Marshal(outSlice)
	if err != nil {
		panic(err)
	}

	// convert to string
	itemsClusterRole := string(yamlFile)

	// init struct with var
	Msg := struct {
		ClusterRoleName   string
		Items             string `yaml:"out"`
		MessageLoggedUser string
	}{
		ClusterRoleName:   clusterRoleName,
		Items:             itemsClusterRole,
		MessageLoggedUser: username,
	}

	// parse html
	t, _ := template.ParseFiles("tmpl/descshow.html")

	// execute
	err = t.Execute(w, Msg)
	if err != nil {
		return
	}

	// redirect to page with description
	http.Redirect(w, r, "/descshow", http.StatusSeeOther)
}
