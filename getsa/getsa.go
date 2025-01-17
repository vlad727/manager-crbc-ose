package getsa

import (
	"log"
	"net/http"
	"os"
	"sigs.k8s.io/yaml"
	"text/template"
	"webapp/getsacollect"
	"webapp/loggeduser"
)

var (
	logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
)

func GetSa(w http.ResponseWriter, r *http.Request) {
	defer logger.Println("INFO: Func GetSa finished")
	logger.Println("INFO: Func GetSa started")
	// send request to parse, trim and decode jwt, get map with user and groups
	UserAndGroups := loggeduser.LoggedUserRun(r) // get logged user and groups, ex:  map[ose.test.user:[ipausers tuz-endless]]

	logger.Println("INFO: Got message form LoggedUserRun")
	log.Println(UserAndGroups)

	// name of logged user
	var username string

	// get logged user name from map
	for k, _ := range UserAndGroups {
		username = k
		break
	}

	// Создаем новую карту для каждого вызова функции таким образом мы очистим данные из mNsAndSa
	mNsAndSa := make(map[string][]string)

	// send map to func GetSaCollect and return M3 map and Sl1 slice
	mNsAndSa, _ = getsacollect.GetSaCollect(UserAndGroups)
	// Sl1 will be skipped  because we don't need here, Sl1 it's slice with namespace name and service account name example below:
	// my-test-ns: my-test-sa
	// M3 it's map with namespace name and service account name like below:
	// ose-test-ns:
	// - default
	// - ose-sa

	// Marshal to yaml for out to web page
	yamlFile, err := yaml.Marshal(mNsAndSa)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// convert to string for struct if you do not convert it will be in bytes
	str := string(yamlFile)
	// parse html
	t, err := template.ParseFiles("tmpl/getsa.html")
	if err != nil {
		log.Printf("Error parsing template: %v\n", err)
		return
	}
	// init struct and var
	Msg := struct {
		Message           string `yaml:"message"`
		MessageLoggedUser string
	}{
		Message:           str,
		MessageLoggedUser: username,
	}
	// execute
	err = t.Execute(w, Msg)
	if err != nil {
		return
	}
	mNsAndSa = nil
}
