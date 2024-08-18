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

func CreatePodPatch(pod *apiv1.Pod) ([]byte, error) {
	labels := pod.ObjectMeta.Labels
	labels["example-webhook"] = "it-worked"

	patches := []PatchOperation{{
		Op:    "add",
		Path:  "/metadata/labels",
		Value: labels,
	}}

	return json.Marshal(patches)
}
