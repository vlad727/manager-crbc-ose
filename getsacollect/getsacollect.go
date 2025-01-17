package getsacollect

import (
	"golang.org/x/net/context"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/strings/slices"
	"log"
	"webapp/clientgo"
	"webapp/readfiles"
)

type StructGetSa struct {
	Collection map[string]string
}

type StructForRbSubject struct {
	Subject []string
}

// func compareSlice search items from slice in another slice
func compareSlice(groups, roleBindingSubjects []string) bool {

	for _, group := range groups {
		if slices.Contains(roleBindingSubjects, group) {
			return true
		}
	}
	return false
}

// func return 2 values for getsa func and crbcmain
func GetSaCollect(LoggedUser map[string][]string) (map[string][]string, []map[string]string) {
	// logging
	log.Println("Func GetSaCollect run ...")

	// Run func for read file with user admin
	UserAdmin, err := readfiles.ReadFile()
	if err != nil {
		log.Printf("Cant get data from about user admin; %s", err)
	}

	// logging
	log.Printf("Got it from func read file %s", UserAdmin)

	var UserName string
	var Groups []string

	// iterate over map to assign data to new client-go
	for k, v := range LoggedUser {
		UserName = k
		Groups = v
	}

	// Get list RoleBindings in namespaces
	listRB, err := clientgo.Сlientset.RbacV1().RoleBindings("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		log.Printf("Failed to get rolebindings. %s", err)
		log.Println(err)
	}
	// Temp string for name from RoleBinding Subject
	var RoleBindingSubjects []string

	// AllowedNsSlice slice for allowed namespaces
	var AllowedNsSlice []string

	// Iterate over role-bindings
	for _, el := range listRB.Items { // get rolebinding from list -> then iterate over this rb -> get el.Subject
		// Iterate over Subjects to get name (also it contains: apiGroup, kind, namespace )
		for _, x := range el.Subjects {
			log.Println(x.Name)
			RoleBindingSubjects = append(RoleBindingSubjects, x.Name) //May be group or username from ldap example or may array of users

		}
		// check condition: if clusterRole == admin and linked with user or group add namespace to allowed list
		// readfiles.UserAdmin <-- get var from configmap func ReadFile
		if el.RoleRef.Name == UserAdmin && slices.Contains(RoleBindingSubjects, UserName) || compareSlice(Groups, RoleBindingSubjects) {
			AllowedNsSlice = append(AllowedNsSlice, el.Namespace)
			RoleBindingSubjects = nil // clear slice
		}

	}

	// logging to know which one namespace we got
	log.Printf("Allowed namespaces: %s", AllowedNsSlice)

	// slice with allowed namespaces example
	//ns := []string{"ose-test-namespace-1", "ose-test-namespace-11", "ose-test-ns", "ose-groups"}

	// get all service accounts
	ListSa, _ := clientgo.Сlientset.CoreV1().ServiceAccounts("").List(context.Background(), v1.ListOptions{})

	// get all service accounts and their namespaces

	var Sl1 []map[string]string //temporary slice for namespaces

	for _, x := range ListSa.Items {
		if slices.Contains(AllowedNsSlice, x.Namespace) { // check that x.Namespace allowed
			M2 := map[string]string{ // add namespace name and service account name to map
				x.Namespace: x.Name}
			Sl1 = append(Sl1, M2)        // add map to slice
			M2 = make(map[string]string) // clear map M1
		}
	}
	// logging
	log.Printf("Slice Sl1 %s", Sl1) // <<< slice with map ns:sa
	// здесь у нас slice в нем мапа и каждому sa свой sa, что нам не подходит для "отдачи" на страницу
	// [map[ose-groups:default] map[ose-test-namespace-1:default] map[ose-test-namespace-1:test-sa] map[ose-test-namespace-11:default] map[ose-test-ns:default] map[ose-test-ns:ose-sa]]
	M3 := make(map[string][]string) // init empty map
	for _, x := range Sl1 {         //iterate over slice which one contain maps
		for k, v := range x {
			if _, ok := M3[k]; !ok { // при первой итерации у нас нет ключа и это false
				M3[k] = make([]string, 0) // <<< значит мы создаем срез " make([]string, 0)" в нашей мапе M3 и в key добавляем наше имя namespace "ose-groups"
			}
			M3[k] = append(M3[k], v) // здесь мы в нашей мапе в срез добавляем наш service account "default" а ключ у нас имя namespace
			/*
				if ok {
					M3[k] = append(M3[k], v)
				} else {
				M3[k] = make([]string, 0)
				}
			*/
		}
	}
	/*
			При второй итерации берем уже след namespace name и если имя ns у нас небыло, в нашем случае это будет ose-test-namespace-1,
			видим, что его так же нет в нашей мап, повторяем процедуру, след ns так ose-test-namespace-1 в данном случае у нас уже есть ключ со значемнием
			ose-test-namespace-1 и значит мы не создаем новый срез а сразу делаем append с этим ключем ose-test-namespace-1 в slice и получаем
		    ose-test-namespace-1:[default test-sa] <<<  namespace и его 2 service account и т.д. в итоге получаем то, что нам было необходимо

	*/
	// logging
	log.Printf("Map M3 %s", M3)
	/*
		What I want to see in M3:
			ose-test-ns:
			- deployer
			- default
			- builder
	*/
	return M3, Sl1
	//  map[ose-groups:[default] ose-test-namespace-1:[default test-sa] ose-test-namespace-11:[default] ose-test-ns:[default ose-sa]]
}

/*      How subject looks like

roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: admin
subjects:
- apiGroup: rbac.authorization.k8s.io
  kind: User
  name: ose.test.user      << this is x.Name var or, it may like below array of users
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: admin
subjects:
- apiGroup: rbac.authorization.k8s.io
  kind: User
  name: ose.test.user      << this is x.Name var
- apiGroup: rbac.authorization.k8s.io
  kind: User
  name: alice     << this is x.Name var or like below with group
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: admin
subjects:
- apiGroup: rbac.authorization.k8s.io
  kind: group
  name: admins      << this is x.Name var
*/
