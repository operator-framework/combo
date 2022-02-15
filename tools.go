//go:build tools
// +build tools

package tools

import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint" // Better linting
	_ "github.com/onsi/ginkgo/ginkgo"
)
