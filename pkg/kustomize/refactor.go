package kustomize

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/int128/kustomtree/pkg/resource"
	"gopkg.in/yaml.v3"
	"sigs.k8s.io/kustomize/api/types"
)

func Refactor(manifest *Manifest, dryRun bool) error {
	for _, resourceManifest := range manifest.ResourceManifests {
		for _, r := range resourceManifest.Resources {
			desiredPath := r.DesiredPath()
			if desiredPath == "" {
				log.Printf("SKIP: resource %s: unknown kind %s", resourceManifest.Path, r.Kind)
				continue
			}
			if resourceManifest.Path != desiredPath {
				if dryRun {
					log.Printf("DRYRUN: move resource %s -> %s", resourceManifest.Path, desiredPath)
					continue
				}
				log.Printf("WRITE: resource %s -> %s", resourceManifest.Path, desiredPath)
				fullpath := filepath.Join(resourceManifest.Basedir, desiredPath)
				if err := resource.Write(fullpath, []*resource.Resource{r}); err != nil {
					return fmt.Errorf("could not write to %s: %w", desiredPath, err)
				}
			}
		}
	}

	for _, patchManifest := range manifest.PatchesStrategicMergeManifests {
		for _, r := range patchManifest.Resources {
			desiredPath := r.DesiredPath()
			if desiredPath == "" {
				log.Printf("SKIP: patchesStrategicMerge %s: unknown kind %s", patchManifest.Path, r.Kind)
				continue
			}
			if patchManifest.Path != desiredPath {
				if dryRun {
					log.Printf("DRYRUN: move patchesStrategicMerge %s -> %s", patchManifest.Path, desiredPath)
					continue
				}
				// TODO: aggregate nodes to same file
				log.Printf("WRITE: patchesStrategicMerge %s -> %s", patchManifest.Path, desiredPath)
				fullpath := filepath.Join(patchManifest.Basedir, desiredPath)
				if err := resource.Write(fullpath, []*resource.Resource{r}); err != nil {
					return fmt.Errorf("could not write to %s: %w", desiredPath, err)
				}
			}
		}
	}

	manifest.Kustomization.Resources = manifest.DesiredResources()
	manifest.Kustomization.PatchesStrategicMerge = manifest.DesiredPatchesStrategicMerge()
	if dryRun {
		log.Printf("DRYRUN: update %s", manifest.Path)
		return nil
	}
	log.Printf("WRITE: update %s", manifest.Path)
	if err := write(manifest.Path, manifest.Kustomization); err != nil {
		return fmt.Errorf("could not write to %s: %w", manifest.Path, err)
	}
	return nil
}

func write(name string, k *types.Kustomization) error {
	f, err := os.Create(name)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}
	defer f.Close()
	e := yaml.NewEncoder(f)
	e.SetIndent(2)
	if err := e.Encode(k); err != nil {
		return fmt.Errorf("could not encode: %w", err)
	}
	if err := e.Close(); err != nil {
		return fmt.Errorf("close error: %w", err)
	}
	return nil
}
