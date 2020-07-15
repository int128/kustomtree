package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/int128/kustomtree/pkg/kustomize"
)

type options struct {
	dryRun bool
}

func Run(osArgs []string) error {
	f := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	var o options
	f.BoolVar(&o.dryRun, "dry-run", false, "If set, do not write files actually")
	if err := f.Parse(osArgs[1:]); err != nil {
		return fmt.Errorf("wrong argument: %w", err)
	}
	if f.NArg() < 1 {
		return fmt.Errorf("you need to set at least 1 path")
	}
	for _, path := range f.Args() {
		if err := run(path, o); err != nil {
			return err
		}
	}
	return nil
}

func run(name string, o options) error {
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
			kustomize.Refactor(manifest, o.dryRun)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("could not find YAMLs: %w", err)
	}
	return nil
}
