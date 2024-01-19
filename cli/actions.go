package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/rogerwelin/cfnctl/aws"
	"github.com/rogerwelin/cfnctl/commands"
	"github.com/rogerwelin/cfnctl/pkg/client"
	"github.com/rogerwelin/cfnctl/utils"
)

func (p plan) run() error {
	ctl, err := commands.CommandBuilder(p.templatePath, p.paramFile, false)
	if err != nil {
		return err
	}

	_, err = commands.Plan(ctl, true)

	return err
}

func (v validate) run() error {
	greenBold := color.New(color.Bold, color.FgHiGreen).SprintFunc()
	err := commands.Validate(v.templatePath)
	if err != nil {
		return err
	}
	fmt.Printf("%s The configuration is valid.\n", greenBold("Success!"))
	return nil
}

func (a apply) run() error {
	ctl, err := commands.CommandBuilder(a.templatePath, a.paramFile, a.autoApprove)
	if err != nil {
		return err
	}
	err = commands.Apply(ctl)
	return err
}

func (d destroy) run() error {
	ctl, err := commands.CommandBuilder(d.templatePath, "", d.autoApprove)
	if err != nil {
		return err
	}

	err = commands.Destroy(ctl)
	return err
}

func (v version) run() error {
	err := commands.OutputVersion(v.version, os.Stdout)
	return err
}

func (o output) run() error {
	svc, err := aws.NewAWS()
	if err != nil {
		return err
	}

	ctl := client.New(
		client.WithSvc(svc),
		client.WithOutput(os.Stdout),
	)

	err = commands.Output(ctl)
	return err
}

func (d drift) run() error {
	svc, err := aws.NewAWS()
	if err != nil {
		return err
	}
	stackName := utils.TrimFileSuffix(d.templatePath)

	ctl := client.New(
		client.WithSvc(svc),
		client.WithStackName(stackName),
		client.WithOutput(os.Stdout),
	)
	// get drift id
	id, err := ctl.StackDriftInit()

	if err != nil {
		return err
	}

	// poll for completion
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		status, err := ctl.GetDriftStatus(id)
		if err != nil {
			return err
		}
		if status == "DETECTION_COMPLETE" {
			break
		}
	}

	status, err := ctl.GetStackDriftInfo()
	if err != nil {
		return err
	}

	for _, item := range status {
		if item.StackResourceDriftStatus != "IN_SYNC" {
			err := utils.JsonDiff(*item.ExpectedProperties, *item.ActualProperties, ctl.Output)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
