package main

import (
	"github.com/operator-framework/combo/internal/cli/root"
)

func main() {
	root.NewCmd().Execute()
}
