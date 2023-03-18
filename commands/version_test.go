package commands

import (
	"bytes"
	"testing"
)

func TestVersion(t *testing.T) {

	expectedStr := "Cfnctl version v0.1.1\n"
	buf := &bytes.Buffer{}
	err := OutputVersion("0.1.1", buf)
	if err != nil {
		t.Errorf("Expected ok but got: %v", err)
	}

	if buf.String() != expectedStr {
		t.Errorf("Expected str:\n %s but got:\n %s", expectedStr, buf.String())
	}

	expectedStr2 := "Cfnctl version v1.0.0\n"
	buf.Reset()

	err = OutputVersion("1.0.0", buf)
	if err != nil {
		t.Errorf("Expected ok but got: %v", err)
	}

	if buf.String() != expectedStr2 {
		t.Errorf("Expected str:\n%s but got:\n %s", expectedStr2, buf.String())
	}

}
