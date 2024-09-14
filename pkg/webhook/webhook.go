package webhook

import (
	"fmt"
	"log/slog"
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

func getRequiredInitContainer(pod *apiv1.Pod) (apiv1.Container, error) {
	annotations := getPodAnnotations(&pod.ObjectMeta)

	initContainerName, ok := annotations[annotationInitContainerName]
	if ok == false || initContainerName == "" {
		return apiv1.Container{}, fmt.Errorf("No init container name provided")

	}
	slog.Info("Looking for init container", slog.String("initContainerName", initContainerName))
	presentContainers := FetchInjectableInitContainers().Items
	for _, container := range presentContainers {
		if container.Metadata.Name == initContainerName {
			slog.Info("Found init container", slog.String("initContainerName", initContainerName))
			return container.Spec, nil
		}
	}
	return apiv1.Container{}, fmt.Errorf("Init container %s not found", initContainerName)
}

func NewAdmissionResponse(pod apiv1.Pod, admissionReviewReq v1beta1.AdmissionReview) (*v1beta1.AdmissionResponse, error) {

	if !isMutationNeeded(&pod.ObjectMeta) {
		return &v1beta1.AdmissionResponse{
			UID:     admissionReviewReq.Request.UID,
			Allowed: true,
		}, nil
	}

	containerToInject, err := getRequiredInitContainer(&pod)
	if err != nil {
		return &v1beta1.AdmissionResponse{
			UID:     admissionReviewReq.Request.UID,
			Allowed: false,
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}, nil
	}
	slog.Info("About to create patch for pod", slog.String("pod", pod.Name))
	patchBytes, err := CreatePodPatch(&pod, containerToInject)
	if err != nil {
		return nil, err
	}

	return &v1beta1.AdmissionResponse{
		UID:     admissionReviewReq.Request.UID,
		Allowed: true,
		Patch:   patchBytes,
	}, nil
}
