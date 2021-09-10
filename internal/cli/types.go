package cli

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/rogerwelin/cfnctl/internal/commands"
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
	err := commands.Plan(p.TemplatePath, p.ParamFile)
	if err != nil {
		return err
	}

	return nil
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
	fmt.Println("applying")
	return nil
}

func (d *Destroy) Run() error {
	fmt.Println("destroying")
	return nil
}

func (v *Version) Run() error {
	err := commands.OutputVersion(v.Version)
	return err
}
