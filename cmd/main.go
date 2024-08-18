package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/danielstepan/k8s-init-injector/pkg/config"
	"github.com/danielstepan/k8s-init-injector/pkg/handler"
)

func main() {
	parameters := config.NewServerParameters()

	err := config.LoadKubeConfig()

	if err != nil {
		panic(err.Error())
	}

	http.HandleFunc("/", handler.HandleRoot)
	http.HandleFunc("/mutate", handler.HandleMutate)
	log.Fatal(http.ListenAndServeTLS(":"+strconv.Itoa(parameters.Port), parameters.CertFile, parameters.KeyFile, nil))
}
