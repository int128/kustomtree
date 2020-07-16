package kustomize

import (
	"path/filepath"

	"github.com/int128/kustomtree/pkg/resource"
	"sigs.k8s.io/kustomize/api/types"
)

type ResourceRef struct {
	Path        string        // relative to kustomization.yaml
	ResourceSet *resource.Set // nil if path is directory
}

type PatchStrategicMergeRef struct {
	Path        string        // relative to kustomization.yaml
	ResourceSet *resource.Set // nil if path is directory
}

type Manifest struct {
	Path                  string // full-path
	Kustomization         *types.Kustomization
	Resources             []ResourceRef
	PatchesStrategicMerge []PatchStrategicMergeRef
}

func (m *Manifest) Basedir() string {
	return filepath.Dir(m.Path)
}
