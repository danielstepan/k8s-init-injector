package webhook

import (
	"strings"

	v1beta1 "k8s.io/api/admission/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	annotationInjectKey         = "k8s-init-injector/inject"
	annotationInitContainerName = "k8s-init-injector/container"
)

func getPodAnnotations(metadata *metav1.ObjectMeta) map[string]string {
	annotations := metadata.GetAnnotations()
	if annotations == nil {
		return map[string]string{}
	}
	return annotations
}

func isMutationNeeded(metadata *metav1.ObjectMeta) bool {
	annotations := getPodAnnotations(metadata)

	shouldInject := strings.ToLower(annotations[annotationInjectKey])
	return shouldInject == "true"
}

func NewAdmissionResponse(pod apiv1.Pod, admissionReviewReq v1beta1.AdmissionReview) (*v1beta1.AdmissionResponse, error) {
	patchBytes, err := CreatePodPatch(&pod)
	if err != nil {
		return nil, err
	}

	if !isMutationNeeded(&pod.ObjectMeta) {
		return &v1beta1.AdmissionResponse{
			UID:     admissionReviewReq.Request.UID,
			Allowed: true,
		}, nil
	}
	return &v1beta1.AdmissionResponse{
		UID:     admissionReviewReq.Request.UID,
		Allowed: true,
		Patch:   patchBytes,
	}, nil

}
