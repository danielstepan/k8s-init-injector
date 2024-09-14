package handler

import (
	"encoding/json"
	"io"
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
	body, err := io.ReadAll(r.Body)
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

	admissionReviewResponse, err := webhook.NewAdmissionResponse(pod, admissionReviewReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	admissionReview := v1beta1.AdmissionReview{
		Response: admissionReviewResponse,
	}

	respBytes, err := json.Marshal(&admissionReview)
	if err != nil {
		http.Error(w, "could not marshal response", http.StatusInternalServerError)
		return
	}

	w.Write(respBytes)
}
