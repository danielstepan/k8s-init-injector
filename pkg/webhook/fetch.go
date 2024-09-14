package webhook

import (
	"context"
	"encoding/json"

	config "github.com/danielstepan/k8s-init-injector/pkg/config"
	apiv1 "k8s.io/api/core/v1"
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

	d, err := clientSet.RESTClient().Get().AbsPath("/apis/danielstepan.cz/v1/injectableinitcontainers").DoRaw(context.TODO())
	if err != nil {
		panic(err)
	}

	var initContainerList InjectableInitContainerList
	if err := json.Unmarshal(d, &initContainerList); err != nil {
		panic(err)
	}

	return initContainerList
}
