package main

import (
	"log"

	"github.com/operator-framework/combo/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
