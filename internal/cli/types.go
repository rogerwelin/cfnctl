package cli

import (
	"fmt"
	"os"

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

// Runner interface simplifies command interaction
type Runner interface {
	Run() error
}

// Run executes the function receives command
func (p *Plan) Run() error {
	ctl, err := commands.CommandBuilder(p.TemplatePath, p.ParamFile, false)
	if err != nil {
		return err
	}

	_, err = commands.Plan(ctl, true)

	return err
}

// Run executes the function receives command
func (v *Validate) Run() error {
	greenBold := color.New(color.Bold, color.FgHiGreen).SprintFunc()
	err := commands.Validate(v.TemplatePath)
	if err != nil {
		return err
	}
	fmt.Printf("%s The configuration is valid.\n", greenBold("Success!"))
	return nil
}

// Run executes the function receives command
func (a *Apply) Run() error {
	ctl, err := commands.CommandBuilder(a.TemplatePath, a.ParamFile, a.AutoApprove)
	if err != nil {
		return err
	}
	err = commands.Apply(ctl)
	return err
}

// Run executes the function receives command
func (d *Destroy) Run() error {
	svc, err := aws.NewAWS()
	if err != nil {
		return err
	}

	ctl := client.New(
		client.WithSvc(svc),
		client.WithAutoApprove(d.AutoApprove),
		client.WithTemplatePath(d.TemplatePath),
		client.WithStackName("dynamolambda"),
		client.WithOutput(os.Stdout),
	)

	err = commands.Destroy(ctl)
	return err
}

// Run executes the function receives command
func (v *Version) Run() error {
	err := commands.OutputVersion(v.Version, os.Stdout)
	return err
}
