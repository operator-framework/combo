package cmd

import (
	"testing"

	logr "github.com/go-logr/logr/testing"
)

func TestRoot(t *testing.T) {
	t.Run("Execute", func(t *testing.T) {
		if err := Execute(logr.TestLogger{T: t}); err != nil {
			t.Fatalf("error occurred while executing command: %v", err)
		}
	})
}
