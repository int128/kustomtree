package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
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
	log.Printf("processing %s", name)

	f, err := os.Open(name)
	if err != nil {
		return fmt.Errorf("could not open the file: %w", err)
	}
	defer f.Close()

	d := yaml.NewDecoder(f)
	var n struct {
		Resources []string `yaml:"resources"`
	}
	if err := d.Decode(&n); err != nil {
		return fmt.Errorf("could not decode: %w", err)
	}
	log.Printf("%+v", n)
	return nil
}
