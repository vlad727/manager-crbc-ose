package groups

import (
	"golang.org/x/net/context"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/utils/strings/slices"
	"log"
	"webapp/globalvar"
)

type GroupStruct struct {
	Items []struct {
		Metadata struct {
			Name string `json:"name"`
		} `json:"metadata"`
		Users []string `json:"users"`
	} `json:"items"`
}

func GroupCollect(LoggedUser string) map[string][]string {

	// get groups with RestClient
	listgroups, err := globalvar.Clientset.AppsV1().RESTClient().Get().AbsPath("/apis/user.openshift.io/v1/groups").DoRaw(context.TODO())
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

	var SliceGroupForUser []string

	for _, x := range dataObjet.Items {
		log.Println(x.Metadata.Name) // list group name
		log.Println(x.Users)         // list of slice users
		if slices.Contains(x.Users, LoggedUser) {
			SliceGroupForUser = append(SliceGroupForUser, x.Metadata.Name)
		}
	}

	// logged user and them groups
	M1 := map[string][]string{
		LoggedUser: SliceGroupForUser,
	}
	log.Printf("Collectted data for user: %s", M1)
	return M1
}
