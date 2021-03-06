//go:generate -command controller-gen go run sigs.k8s.io/controller-tools/cmd/controller-gen
//go:generate controller-gen crd:crdVersions=v1 output:crd:dir=../manifests/crds/ paths=./...
//go:generate controller-gen schemapatch:manifests=../manifests/crds/ output:dir=../manifests/crds/ paths=./...
//go:generate controller-gen object:headerFile=../hack/boilerplate.go.txt paths=./...

// Package api contains type definitions for all external versions of combo's APIs.
package api
