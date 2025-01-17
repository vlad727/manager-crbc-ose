// Package getlen iterate over all cluster role and their items to get len for all rules
// example and the end of code
package getlen

import (
	"golang.org/x/net/context"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/strings/slices"
	"log"
	"time"
	"webapp/clientgo"
)

func GetLen(y []string) map[string]int {

	// execution time
	log.Println("Func GetLen started")
	start := time.Now()

	// temp var
	var lenForCrItems int

	// declare map which one will be store name for cluster role and len for their items
	// <cluster role name>: <len for all their items>
	mClusterRoleLen := make(map[string]int)

	// ----------------------------------------------------------------------------------------------------------------
	listCr, _ := clientgo.Сlientset.RbacV1().ClusterRoles().List(context.TODO(), v1.ListOptions{})
	//log.Println(listCr.Rules)
	for _, x := range listCr.Items {
		if slices.Contains(y, x.Name) {
			// iterate over cluster role rules
			for _, el := range x.Rules {
				// temporary slice for rules
				tempslice := [][]string{el.APIGroups, el.ResourceNames, el.Resources, el.Verbs, el.NonResourceURLs}
				// iterate over tempslice
				for _, y := range tempslice {
					// iterate over items for example APIGroups or ResourceNames then get sum of all len items
					for _, z := range y {
						lenForCrItems += len(z)
					}
				}

			}
			mClusterRoleLen[x.Name] = lenForCrItems
			// set len to nil, need for the next cluster role items
			lenForCrItems = 0
		}

	}
	log.Println(mClusterRoleLen)

	// Code to measure
	duration := time.Since(start)
	log.Printf("Time execution for Func GetLen  %s", duration)
	// return map to main
	return mClusterRoleLen
}

/*
If we have such cluster role and rules like below, so the len for all items will be 16
rules:
- apiGroups:
  - ""
  resources:
  - namespaces
  verbs:
  - delete

log.Printf("The len for string namespaces is %d", len("namespaces"))
log.Printf("The len for string delete is %d", len("delete"))
2024/10/22 13:26:13 The len for string namespaces is 10
2024/10/22 13:26:13 The len for string delete is 6
Note: "" won't count
So, after count all items we will know the size(len) for each cluster role and compare it with size(len) provided cluster role

panic: assignment to entry in nil map
You have to initialize the map using the make function (or a map literal) before you can add any elements:
m := make(map[string]float64)
m["pi"] = 3.1416
*/
