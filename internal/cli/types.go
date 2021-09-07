package cli

import (
	"fmt"
)

type Validate struct{}
type Plan struct{}
type Apply struct{}
type Destroy struct{}
type Version struct{}

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
	fmt.Println("version")
	return nil
}
