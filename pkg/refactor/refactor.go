package refactor

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"sigs.k8s.io/kustomize/api/types"

	"github.com/int128/kustomtree/pkg/kustomization"
	"github.com/int128/kustomtree/pkg/refactor/orderedset"
	"github.com/int128/kustomtree/pkg/resource"
)

type Option struct {
	ExcludePathRegexp *regexp.Regexp
}

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

func ComputePlan(m *kustomization.Manifest, o Option) Plan {
	if o.ExcludePathRegexp == nil {
		o.ExcludePathRegexp = regexp.MustCompile(`^$`)
	}

	plan := Plan{
		KustomizationManifest: m,
		Create:                make(map[string][]*resource.Resource),
	}
	var removeSet orderedset.Strings

	var resourceSet orderedset.Strings
	for _, ref := range m.Resources {
		if o.ExcludePathRegexp.MatchString(ref.Path) {
			resourceSet.Append(ref.Path)
			continue
		}
		if ref.ResourceSet == nil {
			resourceSet.Append(ref.Path)
			continue
		}
		for _, r := range ref.ResourceSet.Resources {
			desiredFilename := r.DesiredPath()
			if desiredFilename == "" || ref.Path == desiredFilename {
				resourceSet.Append(ref.Path)
				continue
			}
			resourceSet.Append(desiredFilename)
			plan.Create[desiredFilename] = append(plan.Create[desiredFilename], r)
			removeSet.Append(ref.Path)
		}
	}
	plan.Resources = append(plan.Resources, resourceSet.Get()...)

	//TODO: consider if resource and patch have same name
	var patchSet orderedset.Strings
	for _, ref := range m.PatchesStrategicMerge {
		if o.ExcludePathRegexp.MatchString(ref.Path) {
			patchSet.Append(ref.Path)
			continue
		}
		if ref.ResourceSet == nil {
			patchSet.Append(ref.Path)
			continue
		}
		for _, r := range ref.ResourceSet.Resources {
			desiredFilename := r.DesiredPath()
			if desiredFilename == "" || ref.Path == desiredFilename {
				patchSet.Append(ref.Path)
				continue
			}
			patchSet.Append(desiredFilename)
			plan.Create[desiredFilename] = append(plan.Create[desiredFilename], r)
			removeSet.Append(ref.Path)
		}
	}
	for _, k := range patchSet.Get() {
		plan.PatchesStrategicMerge = append(plan.PatchesStrategicMerge, types.PatchStrategicMerge(k))
	}

	plan.Remove = append(plan.Remove, removeSet.Get()...)
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
