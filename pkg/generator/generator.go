package generator

import (
	"fmt"
	"strings"
)

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

func buildMultiDoc(docs [][]byte) []byte {
	var multiDoc string
	for i := 0; i < len(docs); i++ {
		// fmt.Println(string(docs[i]))
		var docSeperator string
		if i != len(docs)-1 {
			docSeperator = "\n---"
		}

		multiDoc += fmt.Sprintf("%s%s\n", string(docs[i]), docSeperator)
	}
	return []byte(multiDoc)
}
