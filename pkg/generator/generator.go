package generator

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/operator-framework/combo/pkg/combination"
)

type document struct {
	value string
	seen  bool
}

type documents []*document

type Template struct {
	documents documents
}

// with builds the template documents with the combination set specified
func (t *Template) with(combo combination.Set) string {
	var result string

	// For each document in the template evaluate the current combination set
	for _, document := range t.documents {
		incDocument := document.value
		for key, val := range combo {
			incDocument = regexp.MustCompile(key+`\b`).ReplaceAllString(incDocument, val)
		}
		// Add the document if it had replacements or hadn't been seen and is a valid string
		shouldAdd := (incDocument != document.value || !document.seen) && strings.TrimSpace(incDocument) != ""
		if shouldAdd {
			document.seen = true
			result += fmt.Sprintf("---%v", incDocument)
		}
	}
	return result
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

	var result string

	// Wait for the context to end or the combinations to be done
	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
			combination, err := combinations.Next(ctx)
			if err != nil || combination == nil {
				return result, err
			}

			result += splitTemplate.with(combination)
		}
	}
}
