// Package counter count number of cluster role binding with end "crbc" get int and return it to package crbcmain
package counter

import (
	"golang.org/x/net/context"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"strings"
	"webapp/clientgo"
)

func Counter() int {
	// logging
	log.Println("Func Counter started")

	// get all cluster role bindings
	listCRB, err := clientgo.Сlientset.RbacV1().ClusterRoleBindings().List(context.Background(), v1.ListOptions{})
	if err != nil {
		log.Println(err)
		log.Println("Can't list cluster role bindings")
	}

	// count all cluster role bindings
	var numberOfEntities int
	for _, x := range listCRB.Items {
		if strings.Contains(x.Name, "crbc") && x.Name != "manager-crbc-clusterrolebinding-admin" {
			numberOfEntities++

		}
	}
	return numberOfEntities
}
