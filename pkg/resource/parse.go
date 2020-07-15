package resource

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Manifest struct {
	Filename  string // relative to kustomization.yaml
	Basedir   string
	Resources []*Resource
}

type Resource struct {
	rootType
	Node *yaml.Node // subtree of the resource
}

func (r *Resource) DesiredFilename() string {
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

func Parse(name string, basedir string) (*Manifest, error) {
	f, err := os.Open(filepath.Join(basedir, name))
	if err != nil {
		return nil, fmt.Errorf("could not open the file: %w", err)
	}
	defer f.Close()

	var rs []*Resource
	d := yaml.NewDecoder(f)
	for {
		var n yaml.Node
		if err := d.Decode(&n); err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("could not decode the file: %w", err)
		}
		r, err := parseNode(&n)
		if err != nil {
			return nil, fmt.Errorf("could not parse the resource YAML: %w", err)
		}
		rs = append(rs, r)
	}
	return &Manifest{
		Filename:  name,
		Basedir:   basedir,
		Resources: rs,
	}, nil
}

func parseNode(n *yaml.Node) (*Resource, error) {
	var v rootType
	if err := n.Decode(&v); err != nil {
		return nil, fmt.Errorf("could not decode the node: %w", err)
	}
	return &Resource{
		rootType: v,
		Node:     n,
	}, nil
}

type rootType struct {
	APIVersion string       `yaml:"apiVersion"`
	Kind       string       `yaml:"kind"`
	Metadata   metadataType `yaml:"metadata"`
}

type metadataType struct {
	Name string `yaml:"name"`
}
