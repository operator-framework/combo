package template

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"
)

var (
	ErrCouldNotReadFile = errors.New("could not read file")
)

// template contains an array of manifests that can be
// interacted with with its various functions.
type template struct {
	manifests          []string
	processedManifests []string
}

func newTemplate(file io.Reader) (template, error) {
	// Separate the manifests by the yaml separator and build a template with them
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return template{}, fmt.Errorf("%w: %s", ErrCouldNotReadFile, err.Error())
	}

	manifests := strings.Split(string(fileBytes), "---")

	var constructedTemplate template
	for _, manifest := range manifests {
		trimmedManifest := strings.TrimSpace(manifest)
		if trimmedManifest != "" {
			constructedTemplate.manifests = append(constructedTemplate.manifests, trimmedManifest)
		}
	}

	return constructedTemplate, nil
}

// has determines if any of the manifests for the template
// contains the specified string.
func (t *template) has(searchManifest string) bool {
	for _, manifest := range t.processedManifests {
		if manifest == searchManifest {
			return true
		}
	}
	return false
}

// with builds the template manifests with the combination set specified
func (t *template) with(combo map[string]string) {
	// For each manifest in the template evaluate the current combination set
	for _, manifest := range t.manifests {
		for key, val := range combo {
			manifest = regexp.MustCompile(key+`\b`).ReplaceAllString(manifest, val)
		}

		// Add the manifest if it isn't empty and doesn't already exist in the template
		shouldAdd := manifest != "" && !t.has(manifest)
		if shouldAdd {
			t.processedManifests = append(t.processedManifests, manifest)
		}
	}
}
