package commands

import (
	"bytes"
	"testing"

	"github.com/rogerwelin/cfnctl/internal/mock"
	"github.com/rogerwelin/cfnctl/pkg/client"
)

func TestDestroy(t *testing.T) {

	expectedStr := "\nNo changes. No objects need to be destroyed\n\nEither you have not created any objects yet, there is no Stack named stack or the existing objects were already deleted outside of Cfnctl\n\nDestroy complete! Resources: 0 destroyed\n"

	svc := mock.NewMockAPI()
	buf := &bytes.Buffer{}

	ctl := client.New(
		client.WithSvc(svc),
		client.WithStackName("stack"),
		client.WithChangesetName("change-stack"),
		client.WithTemplatePath("testdata/template.yaml"),
		client.WithAutoApprove(true),
		client.WithOutput(buf),
	)

	err := Destroy(ctl)
	if err != nil {
		t.Errorf("Expected err to be nil but got: %v", err)
	}

	if buf.String() != expectedStr {
		t.Errorf("Expected str:\n%s but got:\n %s", expectedStr, buf.String())
	}
}
