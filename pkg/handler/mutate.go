package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"k8s.io/api/admission/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"

	"github.com/danielstepan/k8s-init-injector/pkg/webhook"
)

var (
	universalDeserializer = serializer.NewCodecFactory(runtime.NewScheme()).UniversalDeserializer()
)

func HandleMutate(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "could not read request body", http.StatusInternalServerError)
		return
	}

	var admissionReviewReq v1beta1.AdmissionReview

	if _, _, err := universalDeserializer.Decode(body, nil, &admissionReviewReq); err != nil {
		http.Error(w, "could not deserialize request", http.StatusBadRequest)
		return
	}

	var pod apiv1.Pod
	if err := json.Unmarshal(admissionReviewReq.Request.Object.Raw, &pod); err != nil {
		http.Error(w, "could not unmarshal pod on admission request", http.StatusInternalServerError)
		return
	}

	patchBytes, err := webhook.CreatePodPatch(&pod)
	if err != nil {
		http.Error(w, "could not create patch", http.StatusInternalServerError)
		return
	}

	admissionReviewResponse := v1beta1.AdmissionReview{
		Response: &v1beta1.AdmissionResponse{
			UID:     admissionReviewReq.Request.UID,
			Allowed: true,
			Patch:   patchBytes,
		},
	}

	respBytes, err := json.Marshal(&admissionReviewResponse)
	if err != nil {
		http.Error(w, "could not marshal response", http.StatusInternalServerError)
		return
	}

	w.Write(respBytes)
}
