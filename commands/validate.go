package commands

import (
	"github.com/rogerwelin/cfnctl/aws"
	"github.com/rogerwelin/cfnctl/pkg/client"
	"github.com/rogerwelin/cfnctl/utils"
)

/*
	Mighty Thor grips the snake
	Firmly by its tongue
	Lifts his hammer high to strike
	Soon his work is done
	Vingtor sends the giant snake
	Bleeding to the depth
	Twilight of the thundergod
	Ragnarök awaits
*/

// Validate validates a given CF template
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
