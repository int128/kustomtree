package kustomize

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
	"sigs.k8s.io/kustomize/api/types"

	"github.com/int128/kustomtree/pkg/resource"
)

type Manifest struct {
	ResourceManifests []*resource.Manifest
}

func Parse(name string) (*Manifest, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("could not open the file: %w", err)
	}
	defer f.Close()

	d := yaml.NewDecoder(f)
	var v types.Kustomization
	if err := d.Decode(&v); err != nil {
		return nil, fmt.Errorf("could not decode: %w", err)
	}
	return parseNode(v, filepath.Dir(name))
}

func parseNode(v types.Kustomization, basedir string) (*Manifest, error) {
	var rs []*resource.Manifest
	for _, resourceName := range v.Resources {
		fullpath := filepath.Join(basedir, resourceName)
		s, err := os.Stat(fullpath)
		if err != nil {
			return nil, fmt.Errorf("resource not found: %w", err)
		}
		if s.IsDir() {
			continue
		}
		m, err := resource.Parse(resourceName, basedir)
		if err != nil {
			return nil, fmt.Errorf("could not load %s: %w", resourceName, err)
		}
		rs = append(rs, m)
	}
	return &Manifest{
		ResourceManifests: rs,
	}, nil
}
