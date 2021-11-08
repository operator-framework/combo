package cmd

import "testing"

func TestRoot(t *testing.T) {
	t.Run("Execute", func(t *testing.T) {
		if err := Execute(nil); err != nil {
			t.Fatalf("error occurred while executing command: %v", err)
		}
	})
}
