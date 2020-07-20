# kustomtree ![build](https://github.com/int128/kustomtree/workflows/build/badge.svg)

This is a command line tool for refactoring of [Kustomize](https://github.com/kubernetes-sigs/kustomize) manifests.

It sorts manifests into kind based directories.
For example,

```
.
├── deployment
│   └── helloworld.yaml
├── service
│   └── helloworld.yaml
└── kustomization.yaml
```

It will perform the following steps.

1. Find `kustomization.yaml`.
1. Find its dependencies such as `resources` and `patchesStrategicMerge`.
1. Rename the dependencies to `KIND/NAME.yaml`.


## Getting Started

```sh
go get github.com/int128/kustomtree
```


## Contributions

This is an open source software.
Feel free to open issues and pull requests.
