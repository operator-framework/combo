package template

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	comboErrors "github.com/operator-framework/combo/pkg/errors"
	"gopkg.in/yaml.v2"
)

var (
	ErrEmptyFile = errors.New("empty file")
)

// template contains an array of manifests that can be
// interacted with with its various functions.
type template struct {
	manifests          []string
	processedManifests []string
}

// validateFile is a simple wrapper to ensure the file we're using exists, is readable,
// and is valid YAML
func (t *template) validate() error {
	if len(t.manifests) == 0 {
		return ErrEmptyFile
	}

	for _, manifest := range t.manifests {
		var holder interface{}
		if err := yaml.Unmarshal([]byte(manifest), &holder); err != nil {
			return err
		}
	}
	return nil
}

func newTemplate(file io.Reader) (template, error) {
	// Separate the manifests by the yaml separator and build a template with them
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return template{}, fmt.Errorf("%w: %s", comboErrors.ErrCouldNotReadFile, err.Error())
	}
	manifests := strings.Split(string(fileBytes), "---")

	var constructedTemplate template
	for _, manifest := range manifests {
		trimmedManifest := strings.TrimSpace(manifest)
		if trimmedManifest != "" {
			constructedTemplate.manifests = append(constructedTemplate.manifests, trimmedManifest)
		}
	}

	// Validate template that is constructed to ensure proper yaml file
	if err := constructedTemplate.validate(); err != nil && !errors.Is(err, ErrEmptyFile) {
		return constructedTemplate, fmt.Errorf("failed to validate file specified: %w", err)
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
