package cli

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/rogerwelin/cfnctl/internal/aws"
	"github.com/rogerwelin/cfnctl/internal/commands"
	"github.com/rogerwelin/cfnctl/pkg/client"
)

type Validate struct {
	TemplatePath string
}

type Plan struct {
	TemplatePath string
	ParamFile    string
}

type Apply struct {
	AutoApprove  bool
	TemplatePath string
	ParamFile    string
}

type Destroy struct {
	AutoApprove  bool
	TemplatePath string
}

type Version struct {
	Version string
}

type CLIRunner interface {
	Run() error
}

func (p *Plan) Run() error {
	ctl, err := commands.CommandBuilder(p.TemplatePath, p.ParamFile, false)
	if err != nil {
		return err
	}

	_, err = commands.Plan(ctl, true)

	return err
}

func (v *Validate) Run() error {
	greenBold := color.New(color.Bold, color.FgHiGreen).SprintFunc()
	err := commands.Validate(v.TemplatePath)
	if err != nil {
		return err
	}
	fmt.Printf("%s The configuration is valid.\n", greenBold("Success!"))
	return nil
}

func (a *Apply) Run() error {
	ctl, err := commands.CommandBuilder(a.TemplatePath, a.ParamFile, a.AutoApprove)
	if err != nil {
		return err
	}
	err = commands.Apply(ctl)
	return err
}

func (d *Destroy) Run() error {
	svc, err := aws.NewAWS()
	if err != nil {
		return err
	}
	ctl := &client.Cfnctl{
		Svc:          svc,
		AutoApprove:  d.AutoApprove,
		TemplatePath: d.TemplatePath,
		StackName:    "dynamolambda",
	}
	err = commands.Destroy(ctl)
	return err
}

func (v *Version) Run() error {
	err := commands.OutputVersion(v.Version)
	return err
}
