package clientgo

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	/*
		// outside cluster client
		config, _    = clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
		Clientset, _ = kubernetes.NewForConfig(config)
	*/

	// inside cluster client, creates the in-cluster config
	Config, _ = rest.InClusterConfig()
	// creates the clientset
	Сlientset, _ = kubernetes.NewForConfig(Config)
)
