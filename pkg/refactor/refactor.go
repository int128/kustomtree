package refactor

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"sigs.k8s.io/kustomize/api/types"

	"github.com/int128/kustomtree/pkg/kustomization"
	"github.com/int128/kustomtree/pkg/resource"
)

type Plan struct {
	KustomizationManifest *kustomization.Manifest

	Resources             []string
	PatchesStrategicMerge []types.PatchStrategicMerge

	// a map of filename and resources to create.
	Create map[string][]*resource.Resource

	// files to remove
	Remove []string
}

func (p *Plan) HasChange() bool {
	return len(p.Remove)+len(p.Create) > 0
}

func ComputePlan(m *kustomization.Manifest) Plan {
	plan := Plan{
		KustomizationManifest: m,
		Create:                make(map[string][]*resource.Resource),
	}

	removeSet := make(map[string]interface{})
	resourceSet := make(map[string]interface{})
	for _, ref := range m.Resources {
		if ref.ResourceSet == nil {
			resourceSet[ref.Path] = nil
			continue
		}
		for _, r := range ref.ResourceSet.Resources {
			desiredFilename := r.DesiredPath()
			if desiredFilename == "" || ref.Path == desiredFilename {
				resourceSet[ref.Path] = nil
				continue
			}
			resourceSet[desiredFilename] = nil
			plan.Create[desiredFilename] = append(plan.Create[desiredFilename], r)
			removeSet[ref.Path] = nil
		}
	}
	for k := range resourceSet {
		plan.Resources = append(plan.Resources, k)
	}

	//TODO: consider if resource and patch have same name
	patchSet := make(map[string]interface{})
	for _, ref := range m.PatchesStrategicMerge {
		if ref.ResourceSet == nil {
			patchSet[ref.Path] = nil
			continue
		}
		for _, r := range ref.ResourceSet.Resources {
			desiredFilename := r.DesiredPath()
			if desiredFilename == "" || ref.Path == desiredFilename {
				patchSet[ref.Path] = nil
				continue
			}
			patchSet[desiredFilename] = nil
			plan.Create[desiredFilename] = append(plan.Create[desiredFilename], r)
			removeSet[ref.Path] = nil
		}
	}
	for k := range patchSet {
		plan.PatchesStrategicMerge = append(plan.PatchesStrategicMerge, types.PatchStrategicMerge(k))
	}

	for k := range removeSet {
		plan.Remove = append(plan.Remove, k)
	}
	return plan
}

func Apply(plan Plan) error {
	for _, name := range plan.Remove {
		fullpath := filepath.Join(plan.KustomizationManifest.Basedir(), name)
		log.Printf("removing %s", fullpath)
		if err := os.Remove(fullpath); err != nil {
			return fmt.Errorf("could not remove %s: %w", fullpath, err)
		}
	}

	for name, resources := range plan.Create {
		fullpath := filepath.Join(plan.KustomizationManifest.Basedir(), name)
		log.Printf("creating %s", fullpath)
		if err := resource.Write(fullpath, resources); err != nil {
			return fmt.Errorf("could not create %s: %w", fullpath, err)
		}
	}

	if plan.HasChange() {
		log.Printf("writing to %s", plan.KustomizationManifest.Path)
		plan.KustomizationManifest.Kustomization.Resources = plan.Resources
		plan.KustomizationManifest.Kustomization.PatchesStrategicMerge = plan.PatchesStrategicMerge
		if err := kustomization.Write(plan.KustomizationManifest.Path, plan.KustomizationManifest.Kustomization); err != nil {
			return fmt.Errorf("could not update kustomization.yaml: %w", err)
		}
	}
	return nil
}
