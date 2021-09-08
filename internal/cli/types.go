package cli

import (
	"fmt"

	"github.com/rogerwelin/cfnctl/internal/commands"
)

type Validate struct {
	TemplateName string
}

type Plan struct {
	TemplateName string
}

type Apply struct {
	AutoApprove  bool
	TemplateName string
}

type Destroy struct {
	AutoApprove  bool
	TemplateName string
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
	fmt.Println("validating")
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
