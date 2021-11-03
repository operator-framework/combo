package cmd

import "testing"

func TestRoot(t *testing.T) {
	t.Run("Execute", func(t *testing.T) {
		Execute()
	})
}
