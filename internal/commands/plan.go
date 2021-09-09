package commands

import (
	"github.com/rogerwelin/cfnctl/internal/aws"
	"github.com/rogerwelin/cfnctl/internal/params"
	"github.com/rogerwelin/cfnctl/internal/utils"
	"github.com/rogerwelin/cfnctl/pkg/client"
)

func planOutput() {

}

func Plan(template, varsfile string) error {

	svc, err := aws.NewAWS()
	if err != nil {
		return err
	}

	templateBody, err := utils.ReadFile(template)
	if err != nil {
		return err
	}

	stackName := utils.TrimFileSuffix(template)

	ctl := &client.Cfnctl{
		Svc:           svc,
		TemplateBody:  string(templateBody),
		StackName:     stackName,
		ChangesetName: stackName,
	}

	// if vars file is supplied
	if varsfile != "" {
		outParams, err := params.MergeFileParams(varsfile)
		if err != nil {
			return err
		}
		err = ctl.CreateChangeSet(ctl.TemplateBody, ctl.StackName, ctl.ChangesetName, outParams)
		if err != nil {
			return err
		}
	} else {
		// no vars file. check if tempalte contains params
		ok, outParams, err := params.CheckInputParams(template)
		if err != nil {
			return err
		}
		// no input params or default value set
		if !ok {
			// create change set
		} else {
			// get user input
			params.BuildInputParams(outParams)
		}
	}

	return nil

}
