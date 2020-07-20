package resource

import "testing"

func TestResource_DesiredPath(t *testing.T) {
	for name, c := range map[string]struct {
		resource    Resource
		desiredName string
	}{
		"generic": {
			Resource{
				rootType: rootType{
					APIVersion: "apps/v1",
					Kind:       "Deployment",
					Metadata:   metadataType{Name: "hello-world"},
				},
			},
			"deployment/hello-world.yaml",
		},
		"hpa": {
			Resource{
				rootType: rootType{
					APIVersion: "autoscaling/v1",
					Kind:       "HorizontalPodAutoscaler",
					Metadata:   metadataType{Name: "hello-world"},
				},
			},
			"hpa/hello-world.yaml",
		},
		"pdb": {
			Resource{
				rootType: rootType{
					APIVersion: "policy/v1beta1",
					Kind:       "PodDisruptionBudget",
					Metadata:   metadataType{Name: "hello-world"},
				},
			},
			"pdb/hello-world.yaml",
		},
		"placeholder": {
			Resource{
				rootType: rootType{
					APIVersion: "v1",
					Kind:       "Service",
					Metadata:   metadataType{Name: "hello-world-${ENV_KEY}"},
				},
			},
			"service/hello-world.yaml",
		},
	} {
		t.Run(name, func(t *testing.T) {
			got := c.resource.DesiredPath()
			if got != c.desiredName {
				t.Errorf("DesiredPath wants %s but was %s", c.desiredName, got)
			}
		})
	}
}
