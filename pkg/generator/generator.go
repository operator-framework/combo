package generator

import (
	"fmt"
	"strings"
)

// Generate accepts an args []map[string]string and generates a multidoc with
// each key/value pair specified within args. It then returns this multidoc in
// the []byte format.
func Generate(args []map[string]string, file []byte) ([]byte, error) {
	stringData := string(file)

	var generatedCombos [][]byte
	for _, combo := range args {
		currentFileCombo := stringData
		for key, val := range combo {
			currentFileCombo = strings.ReplaceAll(currentFileCombo, key, val)
		}
		generatedCombos = append(generatedCombos, []byte(currentFileCombo))
	}

	return buildMultiDoc(generatedCombos), nil
}

// buildMultiDoc Takes a [][]byte and combines each []byte together to form
// a YAML multidoc with the needed seperator.
func buildMultiDoc(docs [][]byte) []byte {
	var multiDoc string
	for i := 0; i < len(docs); i++ {
		var docSeperator string
		if i != len(docs)-1 {
			docSeperator = "\n---"
		}

		multiDoc += fmt.Sprintf("%s%s\n", string(docs[i]), docSeperator)
	}
	return []byte(multiDoc)
}
