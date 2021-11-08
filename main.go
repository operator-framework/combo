package main

import (
	"fmt"
	"os"

	"github.com/operator-framework/combo/cmd"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func main() {
	log := zap.New()
	if err := cmd.Execute(log); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
