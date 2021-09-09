package commands

import "github.com/rogerwelin/cfnctl/internal/params"

func planOutput() {

}

func Plan(template, varsfile string) error {

	_, err := params.CheckParams(template)
	if err != nil {
		return err
	}

	/*
		svc, err := aws.NewAWS()
		if err != nil {
			return err
		}
		templateBody, err := utils.ReadFile(template)
		if err != nil {
			return err
		}

		stackName := utils.TrimFileSuffix(template)

		c := &client.Cfnctl{
			Svc:           svc,
			TemplateBody:  string(templateBody),
			StackName:     stackName,
			ChangesetName: stackName,
		}

		if varsfile != "" {
			_, err = params.MergeParameters(varsfile)
			if err != nil {
				return err
			}
		} else {

		}
	*/

	return nil

}
