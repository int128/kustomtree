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
	ResourceManifests              []*resource.Manifest
	PatchesStrategicMergeManifests []*resource.Manifest
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
	var rs, ps []*resource.Manifest
	for _, resourceFilename := range v.Resources {
		regular, err := isRegularFile(filepath.Join(basedir, resourceFilename))
		if err != nil {
			return nil, fmt.Errorf("resource not found: %w", err)
		}
		if !regular {
			continue
		}
		m, err := resource.Parse(resourceFilename, basedir)
		if err != nil {
			return nil, fmt.Errorf("could not load resource %s: %w", resourceFilename, err)
		}
		rs = append(rs, m)
	}

	for _, patch := range v.PatchesStrategicMerge {
		resourceFilename := string(patch)
		regular, err := isRegularFile(filepath.Join(basedir, resourceFilename))
		if err != nil {
			return nil, fmt.Errorf("patchesStrategicMerge not found: %w", err)
		}
		if !regular {
			continue
		}
		m, err := resource.Parse(resourceFilename, basedir)
		if err != nil {
			return nil, fmt.Errorf("could not load patchesStrategicMerge %s: %w", resourceFilename, err)
		}
		ps = append(ps, m)
	}

	return &Manifest{
		ResourceManifests:              rs,
		PatchesStrategicMergeManifests: ps,
	}, nil
}

func isRegularFile(name string) (bool, error) {
	s, err := os.Stat(name)
	if err != nil {
		return false, fmt.Errorf("could not stat: %w", err)
	}
	return s.Mode().IsRegular(), nil
}
