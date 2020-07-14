package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
	"sigs.k8s.io/kustomize/api/types"
)

func Run(name string) error {
	if err := filepath.Walk(name, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == "kustomization.yaml" {
			return processKustomizationYAML(path)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("could not find YAMLs: %w", err)
	}
	return nil
}

func processKustomizationYAML(name string) error {
	log.Printf("reading %s", name)
	f, err := os.Open(name)
	if err != nil {
		return fmt.Errorf("could not open the file: %w", err)
	}
	defer f.Close()

	d := yaml.NewDecoder(f)
	var n types.Kustomization
	if err := d.Decode(&n); err != nil {
		return fmt.Errorf("could not decode: %w", err)
	}

	for _, resourceName := range n.Resources {
		if err := processResource(resourceName, filepath.Dir(name)); err != nil {
			return fmt.Errorf("could not process %s: %w", name, err)
		}
	}

	//e := yaml.NewEncoder(os.Stdout)
	//e.SetIndent(2)
	//if err := e.Encode(&n); err != nil {
	//	return fmt.Errorf("could not encode: %w", err)
	//}
	return nil
}
