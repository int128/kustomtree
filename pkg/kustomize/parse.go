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
	Path                           string // full-path
	Kustomization                  *types.Kustomization
	ResourceManifests              []*resource.Manifest
	PatchesStrategicMergeManifests []*resource.Manifest
}

func (m *Manifest) DesiredResources() []string {
	var rs []string
	for _, resourceManifest := range m.ResourceManifests {
		for _, r := range resourceManifest.Resources {
			desiredFilename := r.DesiredPath()
			if desiredFilename == "" {
				continue
			}
			rs = append(rs, desiredFilename)
		}
	}
	return rs
}

func (m *Manifest) DesiredPatchesStrategicMerge() []types.PatchStrategicMerge {
	var rs []types.PatchStrategicMerge
	for _, resourceManifest := range m.PatchesStrategicMergeManifests {
		for _, r := range resourceManifest.Resources {
			desiredFilename := r.DesiredPath()
			if desiredFilename == "" {
				continue
			}
			rs = append(rs, types.PatchStrategicMerge(desiredFilename))
		}
	}
	return rs
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
	return parseNode(v, name)
}

func parseNode(v types.Kustomization, name string) (*Manifest, error) {
	basedir := filepath.Dir(name)
	var rs, ps []*resource.Manifest
	for _, resourcePath := range v.Resources {
		regular, err := isRegularFile(filepath.Join(basedir, resourcePath))
		if err != nil {
			return nil, fmt.Errorf("resource not found: %w", err)
		}
		if !regular {
			continue
		}
		m, err := resource.Parse(resourcePath, basedir)
		if err != nil {
			return nil, fmt.Errorf("could not load resource %s: %w", resourcePath, err)
		}
		rs = append(rs, m)
	}

	for _, patch := range v.PatchesStrategicMerge {
		resourcePath := string(patch)
		regular, err := isRegularFile(filepath.Join(basedir, resourcePath))
		if err != nil {
			return nil, fmt.Errorf("patchesStrategicMerge not found: %w", err)
		}
		if !regular {
			continue
		}
		m, err := resource.Parse(resourcePath, basedir)
		if err != nil {
			return nil, fmt.Errorf("could not load patchesStrategicMerge %s: %w", resourcePath, err)
		}
		ps = append(ps, m)
	}

	return &Manifest{
		Path:                           name,
		Kustomization:                  &v,
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
