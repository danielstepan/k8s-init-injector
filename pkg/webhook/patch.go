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

func addInitContainer(pod *apiv1.Pod) []PatchOperation {
	initContainers := pod.Spec.InitContainers

	ic := apiv1.Container{
		Name:  "init-container",
		Image: "busybox",
		Command: []string{
			"echo",
			"Hello, World!",
		},
	}

	initContainers = append(initContainers, ic)

	patches := []PatchOperation{{
		Op:    "add",
		Path:  "/spec/initContainers",
		Value: initContainers,
	}}

	return patches
}

func addLabel(pod *apiv1.Pod) []PatchOperation {
	labels := pod.ObjectMeta.Labels
	labels["example-webhook"] = "it-worked"

	patches := []PatchOperation{{
		Op:    "add",
		Path:  "/metadata/labels",
		Value: labels,
	}}
	return patches
}

func CreatePodPatch(pod *apiv1.Pod) ([]byte, error) {
	var patch []PatchOperation
	//udelat to configurovatgelny
	patch = append(patch, addLabel(pod)...)
	patch = append(patch, addInitContainer(pod)...)

	return json.Marshal(patch)
}
