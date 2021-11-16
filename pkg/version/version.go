package version

import "fmt"

// ComboVersion indicates what version of Combo the binary belongs to
var ComboVersion string

// GitCommit indicates which git commit the binary was built from
var GitCommit string

var KubernetesVersion string

// String returns a pretty string concatenation of ComboVersion and GitCommit
func String() string {
	return fmt.Sprintf(
		"Combo version: %s\nGit commit: %s\nKubernetes version: %s",
		ComboVersion, GitCommit, KubernetesVersion)
}
