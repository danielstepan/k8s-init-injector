package webhook

import (
	"encoding/json"

	apiv1 "k8s.io/api/core/v1"
)

type PatchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

func addInitContainer(pod *apiv1.Pod, container apiv1.Container) []PatchOperation {
	initContainers := pod.Spec.InitContainers

	initContainers = append(initContainers, container)

	patches := []PatchOperation{{
		Op:    "add",
		Path:  "/spec/initContainers",
		Value: initContainers,
	}}

	return patches
}

func addLabel(pod *apiv1.Pod) []PatchOperation {
	labels := pod.ObjectMeta.Labels
	labels["injected-init-container"] = "true"

	patches := []PatchOperation{{
		Op:    "add",
		Path:  "/metadata/labels",
		Value: labels,
	}}
	return patches
}

func CreatePodPatch(pod *apiv1.Pod, container apiv1.Container) ([]byte, error) {
	var patch []PatchOperation
	patch = append(patch, addLabel(pod)...)
	patch = append(patch, addInitContainer(pod, container)...)

	return json.Marshal(patch)
}
