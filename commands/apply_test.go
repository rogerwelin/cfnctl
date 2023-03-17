package commands

import (
	"bytes"
	"testing"

	"github.com/rogerwelin/cfnctl/internal/mock"
	"github.com/rogerwelin/cfnctl/pkg/client"
)

func TestApply(t *testing.T) {

	var expectedStr = `
Cfnctl will perform the following actions:


Plan: 0 to add, 0 to change, 0 to destroy


No changes. Your infrastructure matches the configuration

Cfnctl has compared your real infrastructure against your configuration and found no differences, so no changes are needed.

Apply complete! Resources: 0 added, 0 changed, 0 destroyed
`
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

	err := Apply(ctl)
	if err != nil {
		t.Errorf("Expected str:\n%s but got:\n %s", expectedStr, buf.String())
	}

	if buf.String() != expectedStr {
		t.Fatal("apply output was not expected string")
	}
}
