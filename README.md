# K8s Operators

## Installation

### Setup local k8s cluster
Ensure you have a k8s cluster installed locally.
I use `kind` for my local cluster needs.

To create the cluster, run
```bash
kind create cluster --name operators
```

### Install kubebuilder
```bash
curl -L -o kubebuilder https://go.kubebuilder.io/dl/latest/$(go env GOOS)/$(go env GOARCH) && chmod +x kubebuilder && mv kubebuilder /usr/local/bin/
kubebuilder version

# This should give an output like the following
# Version: main.version{KubeBuilderVersion:"3.10.0", KubernetesVendor:"1.26.1", GitCommit:"0fa57405d4a892efceec3c5a902f634277e30732", BuildDate:"2023-04-15T08:10:35Z", GoOs:"darwin", GoArch:"amd64"}
```

### Setup the domain, API and the groups
```bash
kubebuilder init --domain sarmag.co --repo sarmag.co/todo
```

