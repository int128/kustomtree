package kustomize

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
	"sigs.k8s.io/kustomize/api/types"
)

func Refactor(manifest *Manifest, dryRun bool) error {
	for _, resourceManifest := range manifest.ResourceManifests {
		for _, resource := range resourceManifest.Resources {
			desiredPath := resource.DesiredPath()
			if desiredPath == "" {
				log.Printf("SKIP: resource %s: unknown kind %s", resourceManifest.Path, resource.Kind)
				continue
			}
			if resourceManifest.Path != desiredPath {
				if dryRun {
					log.Printf("DRYRUN: move resource %s -> %s", resourceManifest.Path, desiredPath)
					continue
				}
				log.Printf("WRITE: resource %s -> %s", resourceManifest.Path, desiredPath)
				fullpath := filepath.Join(resourceManifest.Basedir, desiredPath)
				if err := writeNode(fullpath, resource.Node); err != nil {
					return fmt.Errorf("could not write to %s: %w", desiredPath, err)
				}
			}
		}
	}

	for _, patchManifest := range manifest.PatchesStrategicMergeManifests {
		for _, resource := range patchManifest.Resources {
			desiredPath := resource.DesiredPath()
			if desiredPath == "" {
				log.Printf("SKIP: patchesStrategicMerge %s: unknown kind %s", patchManifest.Path, resource.Kind)
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
				if err := writeNode(fullpath, resource.Node); err != nil {
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
	if err := writeKustomizationManifest(manifest.Path, manifest.Kustomization); err != nil {
		return fmt.Errorf("could not write to %s: %w", manifest.Path, err)
	}
	return nil
}

func writeKustomizationManifest(name string, k *types.Kustomization) error {
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

func writeNode(name string, n *yaml.Node) error {
	basedir := filepath.Dir(name)
	if err := os.MkdirAll(basedir, 0700); err != nil {
		return fmt.Errorf("could not mkdir: %w", err)
	}
	f, err := os.Create(name)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}
	defer f.Close()
	e := yaml.NewEncoder(f)
	e.SetIndent(2)
	if err := e.Encode(n); err != nil {
		return fmt.Errorf("could not encode: %w", err)
	}
	if err := e.Close(); err != nil {
		return fmt.Errorf("close error: %w", err)
	}
	return nil
}
