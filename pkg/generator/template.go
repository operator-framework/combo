package generate

import (
	"regexp"
	"strings"
)

// Template contains an array of documents that can be
// interacted with with its various functions.
type template struct {
	documents []string
}

func newTemplate(file string) template {
	// Separate the documents by the yaml separator and build a template with them
	docs := strings.Split(file, "---")

	constructedTemplate := template{}
	constructedTemplate.documents = append(constructedTemplate.documents, docs...)

	return constructedTemplate
}

// has determines if any of the documents for the template
// contains the specified string.
func (t *template) has(find string) bool {
	for _, document := range t.documents {
		if document == find {
			return true
		}
	}
	return false
}

// with builds the template documents with the combination set specified
func (t *template) with(combo map[string]string, to template) template {
	// For each document in the template evaluate the current combination set
	for _, doc := range t.documents {
		incDoc := strings.TrimSpace(doc)
		for key, val := range combo {
			incDoc = regexp.MustCompile(key+`\b`).ReplaceAllString(incDoc, val)
		}

		// Add the document if it isn't empty and doesn't already exist in the template
		shouldAdd := incDoc != "" && !to.has(incDoc)
		if shouldAdd {
			to.documents = append(to.documents, incDoc)
		}
	}
	return to
}
