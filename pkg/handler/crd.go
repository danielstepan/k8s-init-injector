package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	config "github.com/danielstepan/k8s-init-injector/pkg/config"
)

type InitContainer struct {
	Name    string   `json:"name"`
	Image   string   `json:"image"`
	Command []string `json:"command,omitempty"`
	Args    []string `json:"args,omitempty"`
}

// Metadata struct to hold the name field
type Metadata struct {
	Name string `json:"name"`
}

// InjectableInitContainer struct now includes Metadata
type InjectableInitContainer struct {
	Metadata Metadata      `json:"metadata"`
	Spec     InitContainer `json:"spec"`
}

// InjectableInitContainerList remains the same
type InjectableInitContainerList struct {
	Items []InjectableInitContainer `json:"items"`
}

func HandleCRD(w http.ResponseWriter, r *http.Request) {

	clientSet := config.ClientSet

	pods, err := clientSet.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	d, err := clientSet.RESTClient().Get().AbsPath("/apis/danielstepan.cz/v1/injectableinitcontainers").DoRaw(context.TODO())
	if err != nil {
		panic(err)
	}
	fmt.Printf("CRD: %s\n", string(d))

	var initContainerList InjectableInitContainerList
	if err := json.Unmarshal(d, &initContainerList); err != nil {
		panic(err)
	}

	// Extract the Metadata.Name and Spec fields
	var initContainers []InitContainer
	for _, item := range initContainerList.Items {
		fmt.Printf("Name: %s\n", item.Metadata.Name)
		initContainers = append(initContainers, item.Spec)
	}

	fmt.Println("InitContainers: ", initContainers)

}
