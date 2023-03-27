package utils

import (
	"encoding/json"
	"fmt"
	"io"

	diff "github.com/yudai/gojsondiff"
	"github.com/yudai/gojsondiff/formatter"
)

func JsonDiff(expected, actual string, writer io.Writer) error {
	// convert to byte slice
	ex := []byte(expected)
	ac := []byte(actual)

	differ := diff.New()
	d, err := differ.Compare(ex, ac)
	if err != nil {
		return err
	}

	// Output the result
	var diffString string
	var aJson map[string]interface{}
	json.Unmarshal(ex, &aJson)

	config := formatter.AsciiFormatterConfig{
		ShowArrayIndex: true,
		Coloring:       true,
	}

	formatter := formatter.NewAsciiFormatter(aJson, config)
	diffString, err = formatter.Format(d)
	fmt.Fprintln(writer, diffString)

	return nil
}
