package commands

import (
	"github.com/rogerwelin/cfnctl/internal/aws"
	"github.com/rogerwelin/cfnctl/internal/utils"
	"github.com/rogerwelin/cfnctl/pkg/client"
)

func Validate(templatePath string) error {
	svc, err := aws.NewAWS()
	if err != nil {
		return err
	}
	dat, err := utils.ReadFile(templatePath)
	if err != nil {
		return err
	}

	ctl := &client.Cfnctl{
		Svc:          svc,
		TemplateBody: string(dat),
	}

	err = ctl.ValidateCFTemplate()
	if err != nil {
		return err
	}

	return nil
}
