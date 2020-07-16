package kustomize

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
	"sigs.k8s.io/kustomize/api/types"
)

func Refactor(m *Manifest, dryRun bool) error {
	for _, ref := range m.Resources {
		if ref.ResourceSet == nil {
			continue
		}
		for _, mr := range ref.ResourceSet.Resources {
			desiredPath := mr.DesiredPath()
			if desiredPath == "" {
				log.Printf("SKIP: resource %s: unknown kind %s", m.Path, mr.Kind)
				continue
			}
			if ref.Path != desiredPath {
				log.Printf("DRYRUN: move resource %s -> %s", m.Path, desiredPath)
				//TODO: write
			}
		}
	}

	for _, ref := range m.PatchesStrategicMerge {
		if ref.ResourceSet == nil {
			continue
		}
		for _, mr := range ref.ResourceSet.Resources {
			desiredPath := mr.DesiredPath()
			if desiredPath == "" {
				log.Printf("SKIP: patchesStrategicMerge %s: unknown kind %s", m.Path, mr.Kind)
				continue
			}
			if ref.Path != desiredPath {
				log.Printf("DRYRUN: move patchesStrategicMerge %s -> %s", m.Path, desiredPath)
				//TODO: write
			}
		}
	}

	m.Kustomization.Resources = m.DesiredResources()
	m.Kustomization.PatchesStrategicMerge = m.DesiredPatchesStrategicMerge()
	if dryRun {
		log.Printf("DRYRUN: update %s", m.Path)
		return nil
	}
	log.Printf("WRITE: update %s", m.Path)
	if err := write(m.Path, m.Kustomization); err != nil {
		return fmt.Errorf("could not write to %s: %w", m.Path, err)
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
