package generate

import (
	"errors"
	"io"
	"io/ioutil"
	"regexp"
	"strings"
)

var (
	ErrCouldNotReadFile = errors.New("could not read file")
)

// Template contains an array of documents that can be
// interacted with with its various functions.
type template struct {
	documents          []string
	processedDocuments []string
}

func buildTemplate(file io.Reader) (template, error) {
	// Separate the documents by the yaml separator and build a template with them
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return template{}, ErrCouldNotReadFile
	}

	docs := strings.Split(string(fileBytes), "---")

	var constructedTemplate template
	for _, doc := range docs {
		trimmedDoc := strings.TrimSpace(doc)
		if trimmedDoc != "" {
			constructedTemplate.documents = append(constructedTemplate.documents, trimmedDoc)
		}
	}

	return constructedTemplate, nil
}

// has determines if any of the documents for the template
// contains the specified string.
func (t *template) has(searchDocument string) bool {
	for _, document := range t.processedDocuments {
		if document == searchDocument {
			return true
		}
	}
	return false
}

// with builds the template documents with the combination set specified
func (t *template) with(combo map[string]string) {
	// For each document in the template evaluate the current combination set
	for _, doc := range t.documents {
		for key, val := range combo {
			doc = regexp.MustCompile(key+`\b`).ReplaceAllString(doc, val)
		}

		// Add the document if it isn't empty and doesn't already exist in the template
		shouldAdd := doc != "" && !t.has(doc)
		if shouldAdd {
			t.processedDocuments = append(t.processedDocuments, doc)
		}
	}
}
