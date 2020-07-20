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

It will perform the following steps.

1. Find `kustomization.yaml`.
1. Find manifest files in `resources` and `patchesStrategicMerge`.
1. Rename the manifest files to `KIND/NAME.yaml`.
   If filename contains a placeholder (e.g. `-${FOO}`), it will be removed.


## Getting Started

```sh
go get github.com/int128/kustomtree
```


## Contributions

This is an open source software.
Feel free to open issues and pull requests.
