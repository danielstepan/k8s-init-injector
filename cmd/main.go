package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/danielstepan/k8s-init-injector/pkg/config"
	"github.com/danielstepan/k8s-init-injector/pkg/handler"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	slog.Info("Starting the server...")
	config.InitializeFlags()

	slog.Info("Loading the kubeconfig...")
	err := config.LoadKubeConfig()
	if err != nil {
		panic(err.Error())
	}

	parameters := config.NewServerParameters()

	http.HandleFunc("/", handler.HandleRoot)
	http.HandleFunc("/mutate", handler.HandleMutate)
	log.Fatal(http.ListenAndServeTLS(":"+strconv.Itoa(parameters.Port), parameters.CertFile, parameters.KeyFile, nil))
}
