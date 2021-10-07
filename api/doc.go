//go:generate -command controller-gen go run sigs.k8s.io/controller-tools/cmd/controller-gen
//go:generate controller-gen object:headerFile=../hack/boilerplate.go.txt paths=./...

// Package api contains type definitions for all external versions of combo's APIs.
package api
