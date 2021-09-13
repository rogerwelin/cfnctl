package commands

import (
	"os"

	"github.com/rogerwelin/cfnctl/internal/aws"
	"github.com/rogerwelin/cfnctl/internal/utils"
	"github.com/rogerwelin/cfnctl/pkg/client"
)

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
