package kustomize

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func Refactor(manifest *Manifest, dryRun bool) error {
	for _, resourceManifest := range manifest.ResourceManifests {
		for _, resource := range resourceManifest.Resources {
			desiredFilename := resource.DesiredFilename()
			if desiredFilename == "" {
				log.Printf("SKIP: resource %s: unknown kind %s", resourceManifest.Filename, resource.Kind)
				continue
			}
			if resourceManifest.Filename != desiredFilename {
				if dryRun {
					log.Printf("DRYRUN: move resource %s -> %s", resourceManifest.Filename, desiredFilename)
					continue
				}
				log.Printf("WRITE: resource %s -> %s", resourceManifest.Filename, desiredFilename)
				fullpath := filepath.Join(resourceManifest.Basedir, desiredFilename)
				if err := writeNode(fullpath, resource.Node); err != nil {
					return fmt.Errorf("could not write to %s: %w", desiredFilename, err)
				}
			}
		}
	}
	for _, patchManifest := range manifest.PatchesStrategicMergeManifests {
		for _, resource := range patchManifest.Resources {
			desiredFilename := resource.DesiredFilename()
			if desiredFilename == "" {
				log.Printf("SKIP: patchesStrategicMerge %s: unknown kind %s", patchManifest.Filename, resource.Kind)
				continue
			}
			if patchManifest.Filename != desiredFilename {
				if dryRun {
					log.Printf("DRYRUN: move patchesStrategicMerge %s -> %s", patchManifest.Filename, desiredFilename)
					continue
				}
				// TODO: aggregate nodes to same file
				log.Printf("WRITE: patchesStrategicMerge %s -> %s", patchManifest.Filename, desiredFilename)
				fullpath := filepath.Join(patchManifest.Basedir, desiredFilename)
				if err := writeNode(fullpath, resource.Node); err != nil {
					return fmt.Errorf("could not write to %s: %w", desiredFilename, err)
				}
			}
		}
	}
	// TODO: update kustomize.yaml
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
