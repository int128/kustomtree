package resource

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type Set struct {
	Resources []*Resource
}

type Resource struct {
	rootType
	Node *yaml.Node // subtree of the resource
}

func (r *Resource) DesiredPath() string {
	gvk := fmt.Sprintf("%s/%s", r.APIVersion, r.Kind)
	switch gvk {
	case "autoscaling/v1/HorizontalPodAutoscaler":
		return fmt.Sprintf("hpa/%s.yaml", r.Metadata.Name)
	case "policy/v1beta1/PodDisruptionBudget":
		return fmt.Sprintf("pdb/%s.yaml", r.Metadata.Name)
	}
	lowerKind := strings.ToLower(r.Kind)
	return fmt.Sprintf("%s/%s.yaml", lowerKind, r.Metadata.Name)
}
