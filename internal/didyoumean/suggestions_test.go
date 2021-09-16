package didyoumean

import "testing"

var suggestions = []struct {
	in       []string
	given    string
	expected string
}{
	{[]string{"apply", "delete", "plan", "validate", "version"}, "verzion", "version"},
	{[]string{"apply", "delete", "plan", "validate", "version"}, "aply", "apply"},
	{[]string{"apply", "delete", "plan", "validate", "version"}, "pan", "plan"},
	{[]string{"apply", "delete", "plan", "validate", "version"}, "gibberish", ""},
}

func TestSuggestions(t *testing.T) {
	for i, tt := range suggestions {
		actual := NameSuggestion(tt.given, tt.in)
		if actual != tt.expected {
			t.Errorf("test %d: NameSuggestion(%s, %v): expected %s but got %s", i+1, tt.given, tt.in, tt.expected, actual)
		}
	}

}
