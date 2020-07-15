package kustomize

import "log"

func Refactor(manifest *Manifest) {
	for _, resourceManifest := range manifest.ResourceManifests {
		for _, resource := range resourceManifest.Resources {
			desiredFilename := resource.DesiredFilename()
			if desiredFilename == "" {
				log.Printf("SKIP: %s", resourceManifest.Filename)
				continue
			}
			if resourceManifest.Filename != desiredFilename {
				log.Printf("TODO: move %s -> %s", resourceManifest.Filename, desiredFilename)
			}
		}
	}
}
