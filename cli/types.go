package cli

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/rogerwelin/cfnctl/aws"
	"github.com/rogerwelin/cfnctl/commands"
	"github.com/rogerwelin/cfnctl/pkg/client"
)

type validate struct {
	templatePath string
}

type plan struct {
	templatePath string
	paramFile    string
}

type apply struct {
	autoApprove  bool
	templatePath string
	paramFile    string
}

type destroy struct {
	autoApprove  bool
	templatePath string
}

type output struct{}

type version struct {
	version string
}

// Runner interface simplifies command interaction
type Runner interface {
	Run() error
}

// Run executes the function receives command
func (p *plan) run() error {
	ctl, err := commands.CommandBuilder(p.templatePath, p.paramFile, false)
	if err != nil {
		return err
	}

	_, err = commands.Plan(ctl, true)

	return err
}

// Run executes the function receives command
func (v *validate) run() error {
	greenBold := color.New(color.Bold, color.FgHiGreen).SprintFunc()
	err := commands.Validate(v.templatePath)
	if err != nil {
		return err
	}
	fmt.Printf("%s The configuration is valid.\n", greenBold("Success!"))
	return nil
}

// Run executes the function receives command
func (a *apply) run() error {
	ctl, err := commands.CommandBuilder(a.templatePath, a.paramFile, a.autoApprove)
	if err != nil {
		return err
	}
	err = commands.Apply(ctl)
	return err
}

// Run executes the function receives command
func (d *destroy) run() error {
	ctl, err := commands.CommandBuilder(d.templatePath, "", d.autoApprove)
	if err != nil {
		return err
	}

	err = commands.Destroy(ctl)
	return err
}

// Run executes the function receives command
func (v *version) run() error {
	err := commands.OutputVersion(v.version, os.Stdout)
	return err
}

// Run executes the function receives command
func (o *output) run() error {
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
