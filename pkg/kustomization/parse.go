package kustomization

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
	"sigs.k8s.io/kustomize/api/types"

	"github.com/int128/kustomtree/pkg/resource"
)

func Parse(name string) (*Manifest, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("could not open the file: %w", err)
	}
	defer f.Close()

	d := yaml.NewDecoder(f)
	var v types.Kustomization
	if err := d.Decode(&v); err != nil {
		return nil, fmt.Errorf("could not decode: %w", err)
	}
	return parseNode(v, name)
}

func parseNode(v types.Kustomization, name string) (*Manifest, error) {
	basedir := filepath.Dir(name)
	var rs []ResourceRef
	for _, resourcePath := range v.Resources {
		ref, err := determineResourceReference(resourcePath, basedir)
		switch {
		case err != nil:
			return nil, fmt.Errorf("resource not found: %w", err)

		case ref == resourceReferenceFile:
			m, err := resource.Parse(resourcePath, basedir)
			if err != nil {
				return nil, fmt.Errorf("could not load resource %s: %w", resourcePath, err)
			}
			rs = append(rs, ResourceRef{Path: resourcePath, ResourceSet: m})

		default:
			rs = append(rs, ResourceRef{Path: resourcePath})
		}
	}

	var ps []PatchStrategicMergeRef
	//nolint:SA1019
	for _, patch := range v.PatchesStrategicMerge {
		resourcePath := string(patch)
		ref, err := determineResourceReference(resourcePath, basedir)
		switch {
		case err != nil:
			return nil, fmt.Errorf("patchesStrategicMerge not found: %w", err)

		case ref == resourceReferenceFile:
			m, err := resource.Parse(resourcePath, basedir)
			if err != nil {
				return nil, fmt.Errorf("could not load patchesStrategicMerge %s: %w", resourcePath, err)
			}
			ps = append(ps, PatchStrategicMergeRef{Path: resourcePath, ResourceSet: m})

		default:
			ps = append(ps, PatchStrategicMergeRef{Path: resourcePath})
		}
	}

	return &Manifest{
		Path:                  name,
		Kustomization:         &v,
		Resources:             rs,
		PatchesStrategicMerge: ps,
	}, nil
}

type resourceReference int

const (
	_ resourceReference = iota
	resourceReferenceFile
	resourceReferenceDir
	resourceReferenceURL
)

func determineResourceReference(path, baseDir string) (resourceReference, error) {
	// see https://github.com/kubernetes-sigs/kustomize/blob/master/examples/remoteBuild.md#url-format
	if strings.HasPrefix(path, "https://") {
		return resourceReferenceURL, nil
	}
	if strings.HasPrefix(path, "http://") {
		return resourceReferenceURL, nil
	}
	if strings.HasPrefix(path, "ssh://") {
		return resourceReferenceURL, nil
	}
	if strings.HasPrefix(path, "git:") {
		return resourceReferenceURL, nil
	}
	if strings.HasPrefix(path, "github.com/") {
		return resourceReferenceURL, nil
	}

	s, err := os.Stat(filepath.Join(baseDir, path))
	if err != nil {
		return 0, fmt.Errorf("stat error: %w", err)
	}
	if s.IsDir() {
		return resourceReferenceDir, nil
	}
	return resourceReferenceFile, nil
}
