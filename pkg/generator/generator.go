package generator

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/operator-framework/combo/pkg/combination"
)

// document stores a value string to represent the file
// and a seen boolean to be used in the document generation
// algorithm
type document struct {
	value string
	seen  bool
}

type documents []*document

// Template contains an array of documents that can be
// interacted with with its various functions.
type Template struct {
	documents documents
}

// has determines if any of the documents for the template
// contains the specified string.
func (t *Template) has(find string) bool {
	for _, document := range t.documents {
		if document.value == find {
			return true
		}
	}
	return false
}

// build combines all of the template's documents to be
// a valid yaml multidocument.
func (t *Template) build() string {
	var result string
	for _, document := range t.documents {
		result += fmt.Sprintf("---\n%v\n", strings.TrimSpace(document.value))
	}

	return result
}

// with builds the template documents with the combination set specified
func (t *Template) with(combo combination.Set, to Template) Template {
	// For each document in the template evaluate the current combination set
	for _, doc := range t.documents {
		incDoc := doc.value
		for key, val := range combo {
			incDoc = regexp.MustCompile(key+`\b`).ReplaceAllString(incDoc, val)
		}

		// Add the document if it had replacements or wasn't seen, wasn't empty and the to document
		// didn't already have it.
		shouldAdd := (incDoc != doc.value || !doc.seen) && strings.TrimSpace(incDoc) != "" && !to.has(incDoc)
		if shouldAdd {
			doc.seen = true
			to.documents = append(to.documents, &document{value: incDoc})
		}
	}
	return to
}

// Evaluate uses specified template and combination stream to build/return the combinations of
// documents built together
func Evaluate(ctx context.Context, stringTemplate string, combinations combination.Stream) (string, error) {
	// Separate the documents by the yaml seperator and build a template with them
	docs := strings.Split(stringTemplate, "---")
	var splitTemplate Template
	for _, doc := range docs {
		splitTemplate.documents = append(splitTemplate.documents, &document{value: doc})
	}

	var result Template

	// Wait for the context to end or the combinations to be done
	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
			combination, err := combinations.Next(ctx)
			if err != nil {
				return "", err
			}

			if combination == nil {
				return result.build(), nil
			}

			result = splitTemplate.with(combination, result)
		}
	}
}
