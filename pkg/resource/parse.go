package resource

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

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
		Path:      name,
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
