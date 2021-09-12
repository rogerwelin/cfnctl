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

	ctl := &client.Cfnctl{
		Svc:           svc,
		TemplateBody:  string(templateBody),
		TemplatePath:  templateFile,
		StackName:     stackName,
		ChangesetName: stackName,
		VarsFile:      varsFile,
		Output:        os.Stdout,
		AutoApprove:   autoApprove,
	}

	return ctl, nil
}
