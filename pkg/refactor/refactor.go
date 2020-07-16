package refactor

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"sigs.k8s.io/kustomize/api/types"

	"github.com/int128/kustomtree/pkg/kustomize"
	"github.com/int128/kustomtree/pkg/resource"
)

type Plan struct {
	KustomizeManifest *kustomize.Manifest

	Resources             []string
	PatchesStrategicMerge []types.PatchStrategicMerge

	// a map of filename and resources to create.
	Create map[string][]*resource.Resource

	// files to remove
	Remove []string
}

func (p Plan) String() string {
	var s strings.Builder
	_, _ = fmt.Fprintf(&s, "manifest: %s\n", p.KustomizeManifest.Path)
	_, _ = fmt.Fprintf(&s, "  resources: %s\n", p.Resources)
	_, _ = fmt.Fprintf(&s, "  patchesStrategicMerge: %s\n", p.PatchesStrategicMerge)
	_, _ = fmt.Fprintf(&s, "files:\n")
	for name, resources := range p.Create {
		fullpath := filepath.Join(p.KustomizeManifest.Basedir(), name)
		_, _ = fmt.Fprintf(&s, "  + %s %+v\n", fullpath, resources)
	}
	for _, name := range p.Remove {
		fullpath := filepath.Join(p.KustomizeManifest.Basedir(), name)
		_, _ = fmt.Fprintf(&s, "  - %s\n", fullpath)
	}
	return s.String()
}

func ComputePlan(m *kustomize.Manifest) Plan {
	plan := Plan{
		KustomizeManifest: m,
		Create:            make(map[string][]*resource.Resource),
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
	log.Printf("writing to %s", plan.KustomizeManifest.Path)
	plan.KustomizeManifest.Kustomization.Resources = plan.Resources
	plan.KustomizeManifest.Kustomization.PatchesStrategicMerge = plan.PatchesStrategicMerge
	if err := kustomize.Write(plan.KustomizeManifest.Path, plan.KustomizeManifest.Kustomization); err != nil {
		return fmt.Errorf("could not update kustomization.yaml: %w", err)
	}

	for _, name := range plan.Remove {
		fullpath := filepath.Join(plan.KustomizeManifest.Basedir(), name)
		log.Printf("removing %s", fullpath)
		if err := os.Remove(fullpath); err != nil {
			return fmt.Errorf("could not remove %s: %w", fullpath, err)
		}
	}

	for name, resources := range plan.Create {
		fullpath := filepath.Join(plan.KustomizeManifest.Basedir(), name)
		log.Printf("creating %s", fullpath)
		if err := resource.Write(fullpath, resources); err != nil {
			return fmt.Errorf("could not create %s: %w", fullpath, err)
		}
	}
	return nil
}
