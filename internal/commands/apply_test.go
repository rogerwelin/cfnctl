package commands

import (
	"bytes"
	"testing"

	"github.com/rogerwelin/cfnctl/pkg/client"
	"github.com/rogerwelin/cfnctl/pkg/mock"
)

func TestApply(t *testing.T) {

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
		t.Errorf("Expected err to be nil but got: %v", err)
	}

}
