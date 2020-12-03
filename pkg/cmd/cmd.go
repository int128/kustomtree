package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/go-cmp/cmp"

	"github.com/int128/kustomtree/pkg/kustomization"
	"github.com/int128/kustomtree/pkg/kustomize"
	"github.com/int128/kustomtree/pkg/refactor"
)

type options struct {
	dryRun            bool
	excludePathRegexp string
}

func Run(osArgs []string) error {
	f := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	var o options
	f.BoolVar(&o.dryRun, "dry-run", false, "If set, do not write files actually")
	f.StringVar(&o.excludePathRegexp, "exclude-path-regexp", "", "If set, exclude the path from refactoring (e.g. ^vendor/)")
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
	var refactorOption refactor.Option
	if o.excludePathRegexp != "" {
		var err error
		refactorOption.ExcludePathRegexp, err = regexp.Compile(o.excludePathRegexp)
		if err != nil {
			return fmt.Errorf("invalid exclude-path-regexp: %w", err)
		}
	}

	kustomizationYAMLs, err := find(dirname, "kustomization.yaml")
	if err != nil {
		return fmt.Errorf("could not find: %w", err)
	}
	for _, kustomizationYAML := range kustomizationYAMLs {
		manifest, err := kustomization.Parse(kustomizationYAML)
		if err != nil {
			return fmt.Errorf("could not parse YAML: %w", err)
		}
		log.Printf("== PLAN: %s", manifest.Path)
		plan := refactor.ComputePlan(manifest, refactorOption)
		if !plan.HasChange() {
			continue
		}
		log.Println(refactor.Format(plan))
		if o.dryRun {
			continue
		}
		log.Printf("== APPLY: %s", manifest.Path)
		if err := applyAndVerify(plan); err != nil {
			return err
		}
	}
	return nil
}

func applyAndVerify(plan refactor.Plan) error {
	originalBuild, err := kustomize.Build(plan.KustomizationManifest.Basedir())
	if err != nil {
		return fmt.Errorf("kustomize build error: %w", err)
	}
	if err := refactor.Apply(plan); err != nil {
		return fmt.Errorf("could not apply: %w", err)
	}
	refactoredBuild, err := kustomize.Build(plan.KustomizationManifest.Basedir())
	if err != nil {
		return fmt.Errorf("kustomize build error: %w", err)
	}
	if originalBuild != refactoredBuild {
		originalBuildLines := strings.Split(originalBuild, "\r\n")
		refactoredBuildLines := strings.Split(refactoredBuild, "\r\n")
		diff := cmp.Diff(originalBuildLines, refactoredBuildLines)
		return fmt.Errorf("refactoring caused breaking change(s):\n%s", diff)
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
