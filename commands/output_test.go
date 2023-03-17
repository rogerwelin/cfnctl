package commands

import (
	"bytes"
	"strings"
	"testing"

	"github.com/rogerwelin/cfnctl/internal/mock"
	"github.com/rogerwelin/cfnctl/pkg/client"
)

func TestOutput(t *testing.T) {

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

	err := Output(ctl)
	if err != nil {
		t.Errorf("Expected err to be nil but got: %v\n", err)
	}

	if !strings.Contains(buf.String(), "Bucket") {
		t.Error("Output did not contain expected name: Bucket")
	}

	if !strings.Contains(buf.String(), "TestBucket") {
		t.Error("Output did not contain expected value: TestBucket")
	}

}
