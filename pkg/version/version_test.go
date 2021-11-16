package version

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	for _, tt := range []struct {
		name              string
		cliVersion        string
		kubernetesVersion string
		commit            string
		expectedString    string
		expectedFull      string
	}{
		{
			name:              "outputs the correct string",
			cliVersion:        "v0.0.1",
			kubernetesVersion: "v1alpha1",
			commit:            "b6d81d10b34d75c85eea9fd3904298d768f91f4a",
			expectedString: fmt.Sprintf(
				"Combo version: %s\nGit commit: %s\nKubernetes version: %s",
				"v0.0.1",
				"b6d81d10b34d75c85eea9fd3904298d768f91f4a",
				"v1alpha1",
			),
			expectedFull: "v0.0.1-b6d81d10b34d75c85eea9fd3904298d768f91f4a-v1alpha1",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ComboVersion = tt.cliVersion
			KubernetesVersion = tt.kubernetesVersion
			GitCommit = tt.commit

			assert.Equal(t, tt.expectedString, String())
			assert.Equal(t, tt.expectedFull, Full())
		})
	}
}
