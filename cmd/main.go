package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/danielstepan/k8s-init-injector/pkg/config"
	"github.com/danielstepan/k8s-init-injector/pkg/handler"
)

func main() {
	config.InitializeFlags()

	err := config.LoadKubeConfig()
	if err != nil {
		panic(err.Error())
	}

	parameters := config.NewServerParameters()

	http.HandleFunc("/", handler.HandleRoot)
	http.HandleFunc("/mutate", handler.HandleMutate)
	http.HandleFunc("/crd", handler.HandleCRD)
	log.Fatal(http.ListenAndServeTLS(":"+strconv.Itoa(parameters.Port), parameters.CertFile, parameters.KeyFile, nil))
}
