# kustomtree ![build](https://github.com/int128/kustomtree/workflows/build/badge.svg)

This is a command line tool for refactoring of [Kustomize](https://github.com/kubernetes-sigs/kustomize) manifests.

It sorts manifests into kind based directories.
For example,

```
.
├── deployment
│   └── hello-world.yaml
├── service
│   └── hello-world.yaml
├── ingress
│   └── hello-world.yaml
└── kustomization.yaml
```

It finds `kustomization.yaml` and performs the following steps for each `kustomization.yaml`:

1. Run `kustomize build`.
1. Find Kubernetes resources in files specified in `resources` and `patchesStrategicMerge`.
1. Move a Kubernetes resource into a file of `KIND/NAME.yaml`.
   - Placeholder (like `-${FOO}`) is removed.
   - Special character (`:`) is replaced with `-`.
1. Run `kustomize build`.
1. Verify that the rendered manifests of 2 and 4 are same.
   This ensures no breaking change in refactoring.


## Getting Started

```sh
go get github.com/int128/kustomtree
```

### Exclude path

You can exclude manifest(s) from refactoring by `-exclude-path-regexp` flag.

For example, pass `-exclude-path-regexp=^vendor/` to exclude files in `vendor` directory, like:

```
.
├── deployment
│   └── hello-world.yaml
├── vendor
│   └── generated.yaml
└── kustomization.yaml
```


## Contributions

This is an open source software.
Feel free to open issues and pull requests.
