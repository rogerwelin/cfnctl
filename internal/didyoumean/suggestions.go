package didyoumean

import "github.com/agext/levenshtein"

func NameSuggestion(given string, suggestions []string) string {
	for _, suggestion := range suggestions {
		dist := levenshtein.Distance(given, suggestion, nil)
		if dist < 3 {
			return suggestion
		}
	}
	return ""
}
