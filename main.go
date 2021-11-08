package main

import (
	"fmt"
	"os"

	"github.com/operator-framework/combo/cmd"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func main() {
	if err := cmd.Execute(zap.New()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
