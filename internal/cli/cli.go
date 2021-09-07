package cli

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

var (
	version = "dev"
)

func RunCLI(args []string) {
	app := cli.NewApp()
	app.Name = "cfnctl"
	app.Usage = "-"
	app.HelpName = "cfnctl"
	app.EnableBashCompletion = true
	app.UsageText = "cfntl [global options] <subcommand> [args]"
	app.Version = version
	app.HideVersion = true

	app.Commands = []*cli.Command{
		{
			Name:  "apply",
			Usage: "Create or update infrastructure",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "auto-approve",
					Usage: "Skip interactive approval of plan before applying.",
					Value: false,
				},
				&cli.StringFlag{
					Name:  "template-file",
					Usage: "The path of the Cloudformation template you wish to apply",
				},
				&cli.StringFlag{
					Name:  "param-file",
					Usage: "filename. Load parameters from the given file",
				},
				&cli.StringFlag{
					Name:  "param",
					Usage: "foo=bar. Set a value for one of the parameters. Use this option more than once to set more than one parameter",
				},
			},
		},
		{
			Name:  "plan",
			Usage: "Show changes required by the current configuration",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "param-file",
					Usage: "filename",
				},
			},
		},
		{
			Name:  "destroy",
			Usage: "Destroy previously-created infrastructure",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "auto-approve",
					Usage: "Skip interactive approval of plan before applying.",
					Value: false,
				},
			},
		},
		{
			Name:  "validate",
			Usage: "Check whether the configuration is valid",
		},
		{
			Name:  "version",
			Usage: "Show the current Cfnctl version",
		},
	}

	err := app.Run(args)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
