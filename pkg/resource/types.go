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

type Resource struct {
	rootType
	Node *yaml.Node // subtree of the resource
}

// DesiredPath returns a path for the resource.
// It should be form of {kind}/{metadata.name}.yaml.
// Any placeholders, i.e. ${}, are removed from the path.
func (r *Resource) DesiredPath() string {
	resourceFilename := sanitizeResourceFilename(r.Metadata.Name)
	if r.APIVersion == "autoscaling/v1" && r.Kind == "HorizontalPodAutoscaler" {
		return fmt.Sprintf("hpa/%s.yaml", resourceFilename)
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
