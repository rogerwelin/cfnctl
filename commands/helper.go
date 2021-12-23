package commands

import (
	"os"

	"github.com/rogerwelin/cfnctl/aws"
	"github.com/rogerwelin/cfnctl/pkg/client"
	"github.com/rogerwelin/cfnctl/utils"
)

// CommandBuilder returns a new Cfnctl svc
func CommandBuilder(templateFile, varsFile string, autoApprove bool) (*client.Cfnctl, error) {

	svc, err := aws.NewAWS()
	if err != nil {
		return nil, err
	}

	templateBody, err := utils.ReadFile(templateFile)
	if err != nil {
		return nil, err
	}

	stackName := utils.TrimFileSuffix(templateFile)

	ctl := client.New(
		client.WithSvc(svc),
		client.WithTemplateBody(string(templateBody)),
		client.WithTemplatePath(templateFile),
		client.WithStackName(stackName),
		client.WithChangesetName(stackName),
		client.WithVarsFile(varsFile),
		client.WithAutoApprove(autoApprove),
		client.WithOutput(os.Stdout),
	)

	return ctl, nil
}
