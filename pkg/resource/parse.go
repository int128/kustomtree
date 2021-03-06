package resource

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func Parse(name string, basedir string) (*Set, error) {
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
	return &Set{
		Resources: rs,
	}, nil
}

func parseNode(n *yaml.Node) (*Resource, error) {
	r := Resource{Node: n}
	if err := n.Decode(&r); err != nil {
		return nil, fmt.Errorf("could not decode the node: %w", err)
	}
	return &r, nil
}
