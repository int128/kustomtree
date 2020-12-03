package refactor

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"sigs.k8s.io/kustomize/api/types"

	"github.com/int128/kustomtree/pkg/kustomization"
	"github.com/int128/kustomtree/pkg/resource"
)

func TestComputePlan(t *testing.T) {
	t.Run("HappyPath", func(t *testing.T) {
		m := kustomization.Manifest{
			Path: "testdata/kustomization.yaml",
			Resources: []kustomization.ResourceRef{
				{
					Path: "service.yaml",
					ResourceSet: &resource.Set{
						Resources: []*resource.Resource{
							{
								APIVersion: "v1",
								Kind:       "Service",
								Metadata: resource.Metadata{
									Name: "echoserver",
								},
							},
						},
					},
				},
				{
					Path: "deployment.yaml",
					ResourceSet: &resource.Set{
						Resources: []*resource.Resource{
							{
								APIVersion: "apps/v1",
								Kind:       "Deployment",
								Metadata: resource.Metadata{
									Name: "helloworld",
								},
							},
						},
					},
				},
			},
			PatchesStrategicMerge: []kustomization.PatchStrategicMergeRef{
				{
					Path: "patches.yaml",
					ResourceSet: &resource.Set{
						Resources: []*resource.Resource{
							{
								APIVersion: "autoscaling/v1",
								Kind:       "HorizontalPodAutoscaler",
								Metadata: resource.Metadata{
									Name: "example",
								},
							},
							{
								APIVersion: "v1",
								Kind:       "Service",
								Metadata: resource.Metadata{
									Name: "example",
								},
							},
						},
					},
				},
			},
			Kustomization: &types.Kustomization{
				PatchesStrategicMerge: []types.PatchStrategicMerge{},
				Resources: []string{
					"service.yaml",
					"deployment.yaml",
				},
			},
		}
		got := ComputePlan(&m)
		want := Plan{
			KustomizationManifest: &m,
			Resources: []string{
				"service/echoserver.yaml",
				"deployment/helloworld.yaml",
			},
			PatchesStrategicMerge: []types.PatchStrategicMerge{
				"hpa/example.yaml",
				"service/example.yaml",
			},
			Create: map[string][]*resource.Resource{
				"service/echoserver.yaml": {
					{
						APIVersion: "v1",
						Kind:       "Service",
						Metadata: resource.Metadata{
							Name: "echoserver",
						},
					},
				},
				"deployment/helloworld.yaml": {
					{
						APIVersion: "apps/v1",
						Kind:       "Deployment",
						Metadata: resource.Metadata{
							Name: "helloworld",
						},
					},
				},
				"hpa/example.yaml": {
					{
						APIVersion: "autoscaling/v1",
						Kind:       "HorizontalPodAutoscaler",
						Metadata: resource.Metadata{
							Name: "example",
						},
					},
				},
				"service/example.yaml": {
					{
						APIVersion: "v1",
						Kind:       "Service",
						Metadata: resource.Metadata{
							Name: "example",
						},
					},
				},
			},
			Remove: []string{
				"service.yaml",
				"deployment.yaml",
				"patches.yaml",
			},
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	})
}
