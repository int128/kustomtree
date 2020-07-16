package kustomize

import (
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

func (m *Manifest) DesiredResources() []string {
	var a []string
	for _, ref := range m.Resources {
		if ref.ResourceSet == nil {
			a = append(a, ref.Path)
			continue
		}
		for _, r := range ref.ResourceSet.Resources {
			desiredFilename := r.DesiredPath()
			if desiredFilename == "" {
				a = append(a, ref.Path)
				continue
			}
			a = append(a, desiredFilename)
		}
	}
	return a
}

func (m *Manifest) DesiredPatchesStrategicMerge() []types.PatchStrategicMerge {
	var a []types.PatchStrategicMerge
	for _, ref := range m.Resources {
		if ref.ResourceSet == nil {
			a = append(a, types.PatchStrategicMerge(ref.Path))
			continue
		}
		for _, r := range ref.ResourceSet.Resources {
			desiredFilename := r.DesiredPath()
			if desiredFilename == "" {
				a = append(a, types.PatchStrategicMerge(ref.Path))
				continue
			}
			a = append(a, types.PatchStrategicMerge(desiredFilename))
		}
	}
	return a
}
