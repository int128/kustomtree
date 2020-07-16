package resource

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func Write(name string, rs []*Resource) error {
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
	for _, r := range rs {
		if err := e.Encode(r.Node); err != nil {
			return fmt.Errorf("could not encode: %w", err)
		}
	}
	if err := e.Close(); err != nil {
		return fmt.Errorf("close error: %w", err)
	}
	return nil
}
