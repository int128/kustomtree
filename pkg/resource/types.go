package resource

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Manifest struct {
	Path      string // relative to kustomization.yaml
	Basedir   string
	Resources []*Resource
}

type Resource struct {
	rootType
	Node *yaml.Node // subtree of the resource
}

func (r *Resource) DesiredPath() string {
	gvk := fmt.Sprintf("%s/%s", r.APIVersion, r.Kind)
	switch gvk {
	case "apps/v1/Deployment":
		return fmt.Sprintf("deployment/%s.yaml", r.Metadata.Name)
	case "argoproj.io/v1alpha1/Rollout":
		return fmt.Sprintf("rollout/%s.yaml", r.Metadata.Name)
	case "batch/v1beta1/CronJob":
		return fmt.Sprintf("cronjob/%s.yaml", r.Metadata.Name)
	case "v1/Service":
		return fmt.Sprintf("service/%s.yaml", r.Metadata.Name)
	case "autoscaling/v1/HorizontalPodAutoscaler":
		return fmt.Sprintf("hpa/%s.yaml", r.Metadata.Name)
	case "policy/v1beta1/PodDisruptionBudget":
		return fmt.Sprintf("pdb/%s.yaml", r.Metadata.Name)
	case "v1/ConfigMap":
		return fmt.Sprintf("configmap/%s.yaml", r.Metadata.Name)
	}
	return ""
}
