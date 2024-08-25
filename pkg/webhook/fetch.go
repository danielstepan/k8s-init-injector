package webhook

import (
	"context"
	"encoding/json"
	"fmt"

	config "github.com/danielstepan/k8s-init-injector/pkg/config"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type InjectableInitContainerList struct {
	Items []struct {
		Metadata struct {
			Name string `json:"name"`
		} `json:"metadata"`
		Spec apiv1.Container
	} `json:"items"`
}

func FetchInjectableInitContainers() InjectableInitContainerList {

	clientSet := config.ClientSet

	pods, err := clientSet.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("There are %d pods in the cluster\n", len(pods.Items))

	d, err := clientSet.RESTClient().Get().AbsPath("/apis/danielstepan.cz/v1/injectableinitcontainers").DoRaw(context.TODO())
	if err != nil {
		panic(err)
	}
	fmt.Printf("CRD: %s\n", string(d))

	var initContainerList InjectableInitContainerList
	if err := json.Unmarshal(d, &initContainerList); err != nil {
		panic(err)
	}

	fmt.Println("InitContainer List: ", initContainerList)
	return initContainerList
}
