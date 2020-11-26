package resource

import (
	"fmt"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type Set struct {
	Resources []*Resource
}

// Resource represents resources defined in kustomization.yaml.
type Resource struct {
	APIVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`

	Node *yaml.Node `yaml:"-"` // subtree of the resource
}

// Metadata represents a metadata of resource.
type Metadata struct {
	Name string `yaml:"name"`
}

// DesiredPath returns a path for the resource.
// It should be form of {kind}/{metadata.name}.yaml.
// Any placeholders, i.e. ${}, are removed from the path.
func (r *Resource) DesiredPath() string {
	resourceFilename := sanitizeResourceFilename(r.Metadata.Name)
	if r.APIVersion == "autoscaling/v1" && r.Kind == "HorizontalPodAutoscaler" {
		return fmt.Sprintf("hpa/%s.yaml", resourceFilename)
	}
	if r.APIVersion == "autoscaling.k8s.io/v1" && r.Kind == "VerticalPodAutoscaler" {
		return fmt.Sprintf("vpa/%s.yaml", resourceFilename)
	}
	if r.APIVersion == "policy/v1beta1" && r.Kind == "PodDisruptionBudget" {
		return fmt.Sprintf("pdb/%s.yaml", resourceFilename)
	}
	lowerKind := strings.ToLower(r.Kind)
	return fmt.Sprintf("%s/%s.yaml", lowerKind, resourceFilename)
}

var rePlaceholderInYAML = regexp.MustCompile(`-?\${.+?}`)

func sanitizeResourceFilename(name string) string {
	return rePlaceholderInYAML.ReplaceAllString(name, "")
}
