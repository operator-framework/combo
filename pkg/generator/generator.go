package generator

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Generate accepts an args []map[string]string and generates a multidoc with
// each key/value pair specified within args. It then returns this multidoc in
// the []byte format.
func Generate(replacementCombos []map[string]string, file []byte) ([]byte, error) {
	// Error if attempting to read an empty file
	if len(file) == 0 {
		return nil, errors.New("cannot generate combinations for an empty file")
	}

	// Get each document
	documents := strings.Split(string(file), "---")

	// For each document specified, generate its combinations
	var generatedCombos [][]byte
	for _, document := range documents {
		if document == "" {
			continue
		}

		var added bool
		for _, replacementCombo := range replacementCombos {
			documentCombo := document
			for key, val := range replacementCombo {
				documentCombo = regexp.MustCompile(key+`\b`).ReplaceAllString(documentCombo, val)
			}

			if documentCombo != document || !added {
				generatedCombos = append(generatedCombos, []byte(documentCombo))
				added = true
			}
		}

	}

	// Build the multi-doc back up after processing
	return buildMultiDoc(generatedCombos), nil
}

// buildMultiDoc Takes a [][]byte and combines each []byte together to form
// a YAML multidoc with the needed seperator.
func buildMultiDoc(docs [][]byte) []byte {
	var multiDoc string
	for i := 0; i < len(docs); i++ {
		var docSeperator string
		if i != len(docs)-1 {
			docSeperator = "---\n"
		}

		multiDoc += fmt.Sprintf("%s\n%s", strings.TrimSpace(string(docs[i])), docSeperator)
	}
	return []byte(fmt.Sprintf("---\n%v", multiDoc))
}
