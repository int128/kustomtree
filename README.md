# kustomtree

This is a tool for refactoring [Kustomize](https://github.com/kubernetes-sigs/kustomize) manifests.

It finds `kustomization.yaml` and dependent manifests.
For now `resources` and `patchesStrategicMerge` are supported.
It sorts the manifest into directories by kind.

For example,

```
.
├── deployment
│   └── helloworld.yaml
├── service
│   └── helloworld.yaml
└── kustomization.yaml
```
