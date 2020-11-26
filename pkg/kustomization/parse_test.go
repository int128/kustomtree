package kustomization

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"gopkg.in/yaml.v3"
	"sigs.k8s.io/kustomize/api/types"

	"github.com/int128/kustomtree/pkg/resource"
)

func TestParse(t *testing.T) {
	got, err := Parse("testdata/kustomization.yaml")
	if err != nil {
		t.Fatalf("Parse error: %+v", err)
	}
	var wantResource resource.Resource
	wantResource.APIVersion = "apps/v1"
	wantResource.Kind = "Deployment"
	wantResource.Metadata.Name = "helloworld"
	want := &Manifest{
		Path: "testdata/kustomization.yaml",
		Resources: []ResourceRef{
			{Path: "base"},
		},
		PatchesStrategicMerge: []PatchStrategicMergeRef{
			{
				Path: "deployment.yaml",
				ResourceSet: &resource.Set{
					Resources: []*resource.Resource{
						&wantResource,
					},
				},
			},
		},
	}
	o := []cmp.Option{
		cmp.AllowUnexported(
			resource.Resource{},
		),
		cmpopts.IgnoreTypes(
			new(types.Kustomization),
			new(yaml.Node),
		),
	}
	if diff := cmp.Diff(want, got, o...); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}
