package kustomize

import "log"

func Refactor(manifest *Manifest, dryRun bool) {
	for _, resourceManifest := range manifest.ResourceManifests {
		for _, resource := range resourceManifest.Resources {
			desiredFilename := resource.DesiredFilename()
			if desiredFilename == "" {
				log.Printf("SKIP: resource %s: unknown kind %s", resourceManifest.Filename, resource.Kind)
				continue
			}
			if resourceManifest.Filename != desiredFilename {
				if dryRun {
					log.Printf("DRYRUN: move resource %s -> %s", resourceManifest.Filename, desiredFilename)
					continue
				}
				//TODO: write files
			}
		}
	}

	for _, patchManifest := range manifest.PatchesStrategicMergeManifests {
		for _, resource := range patchManifest.Resources {
			desiredFilename := resource.DesiredFilename()
			if desiredFilename == "" {
				log.Printf("SKIP: patchesStrategicMerge %s: unknown kind %s", patchManifest.Filename, resource.Kind)
				continue
			}
			if patchManifest.Filename != desiredFilename {
				if dryRun {
					log.Printf("DRYRUN: move patchesStrategicMerge %s -> %s", patchManifest.Filename, desiredFilename)
					continue
				}
				//TODO: write files
			}
		}
	}
}
