package cli

import (
	"fmt"

	"github.com/rogerwelin/cfnctl/internal/commands"
)

type Validate struct {
	TemplatePath string
}

type Plan struct {
	TemplatePath string
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
	fmt.Println("planning")
	return nil
}

func (v *Validate) Run() error {
	err := commands.Validate(v.TemplatePath)
	if err != nil {
		return err
	}
	fmt.Println("Success! The configuration is valid.")
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
