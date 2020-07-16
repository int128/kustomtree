package kustomize

import (
	"github.com/int128/kustomtree/pkg/resource"
	"sigs.k8s.io/kustomize/api/types"
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
