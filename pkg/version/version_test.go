package version

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	for _, tt := range []struct {
		name           string
		version        string
		commit         string
		expectedString string
		expectedFull   string
	}{
		{
			name:    "outputs the correct full string",
			version: "v0.0.1",
			commit:  "b6d81d10b34d75c85eea9fd3904298d768f91f4a",
			expectedString: fmt.Sprintf(
				"Combo version: %v\nGit commit: %v",
				"v0.0.1",
				"b6d81d10b34d75c85eea9fd3904298d768f91f4a",
			),
			expectedFull: "v0.0.1-b6d81d10b34d75c85eea9fd3904298d768f91f4a",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ComboVersion = tt.version
			GitCommit = tt.commit

			assert.Equal(t, tt.expectedString, String())
			assert.Equal(t, tt.expectedFull, Full())
		})
	}
}
