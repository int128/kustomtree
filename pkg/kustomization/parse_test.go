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
	want := &Manifest{
		Path: "testdata/kustomization.yaml",
		Resources: []ResourceRef{
			{Path: "github.com/octocat/Spoon-Knife/index.html"},
			{Path: "base"},
			{Path: "https://raw.githubusercontent.com/octocat/Spoon-Knife/master/README.md"},
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
							Node: &yaml.Node{
								Kind:   yaml.DocumentNode,
								Line:   1,
								Column: 1,
							},
						},
					},
				},
			},
		},
		PatchesStrategicMerge: []PatchStrategicMergeRef{
			{Path: "github.com/octocat/Spoon-Knife/README.md"},
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
							Node: &yaml.Node{
								Kind:   yaml.DocumentNode,
								Line:   1,
								Column: 1,
							},
						},
					},
				},
			},
			{Path: "https://raw.githubusercontent.com/octocat/Spoon-Knife/master/index.html"},
		},
		Kustomization: &types.Kustomization{
			NamePrefix: "cluster-a-",
			PatchesStrategicMerge: []types.PatchStrategicMerge{
				"github.com/octocat/Spoon-Knife/README.md",
				"deployment.yaml",
				"https://raw.githubusercontent.com/octocat/Spoon-Knife/master/index.html",
			},
			Resources: []string{
				"github.com/octocat/Spoon-Knife/index.html",
				"base",
				"https://raw.githubusercontent.com/octocat/Spoon-Knife/master/README.md",
				"service.yaml",
			},
		},
	}
	o := []cmp.Option{
		cmpopts.IgnoreFields(yaml.Node{}, "Content"),
	}
	if diff := cmp.Diff(want, got, o...); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}
