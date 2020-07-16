package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/int128/kustomtree/pkg/kustomize"
	"github.com/int128/kustomtree/pkg/refactor"
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

func run(dirname string, o options) error {
	kustomizationYAMLs, err := find(dirname, "kustomization.yaml")
	if err != nil {
		return fmt.Errorf("could not find: %w", err)
	}
	for _, kustomizationYAML := range kustomizationYAMLs {
		manifest, err := kustomize.Parse(kustomizationYAML)
		if err != nil {
			return fmt.Errorf("could not parse YAML: %w", err)
		}
		log.Printf("== PLAN")
		plan := refactor.ComputePlan(manifest)
		log.Print(plan)
		if !o.dryRun {
			log.Printf("== APPLY")
			if err := refactor.Apply(plan); err != nil {
				return fmt.Errorf("could not apply: %w", err)
			}
		}
	}
	return nil
}

func find(dirname, filename string) ([]string, error) {
	var a []string
	if err := filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == filename {
			a = append(a, path)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return a, nil
}
