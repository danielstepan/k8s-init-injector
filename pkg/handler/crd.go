package handler

import (
	"net/http"

	"github.com/danielstepan/k8s-init-injector/pkg/webhook"
)

func HandleCRD(w http.ResponseWriter, r *http.Request) {
	webhook.FetchInjectableInitContainers()
	w.Write([]byte("HandleCRD!"))
}
