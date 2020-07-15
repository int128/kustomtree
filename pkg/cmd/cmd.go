package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/int128/kustomtree/pkg/kustomize"
)

func Run(name string) error {
	if err := filepath.Walk(name, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == "kustomization.yaml" {
			log.Printf("reading %s", path)
			manifest, err := kustomize.Parse(path)
			if err != nil {
				return fmt.Errorf("could not parse YAML: %w", err)
			}
			kustomize.Refactor(manifest)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("could not find YAMLs: %w", err)
	}
	return nil
}
