package kustomize

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
	"sigs.k8s.io/kustomize/api/types"
)

func Write(name string, k *types.Kustomization) error {
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
