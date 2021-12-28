package cli

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/rogerwelin/cfnctl/didyoumean"
	"github.com/urfave/cli/v2"
)

var (
	version = "0.1.0"
	cmds    = []string{"apply", "destroy", "plan", "validate", "version", "output", "help"}
)

// RunCLI runs a new instance of cfnctl
func RunCLI(args []string) {
	app := cli.NewApp()
	setCustomCLITemplate(app)
	app.Name = "cfnctl"
	app.Usage = "✨ Terraform cli experience for AWS Cloudformation"
	app.HelpName = "cfnctl"
	app.EnableBashCompletion = true
	app.UsageText = "cfntl [global options] <subcommand> [args]"
	app.Version = version
	app.HideVersion = true
	app.CommandNotFound = func(c *cli.Context, command string) {
		res := didyoumean.NameSuggestion(command, cmds)
		if res == "" {
			fmt.Println("apa") // FIX
		} else {
			fmt.Println("Cfnctl has no command named: " + command + ". Did you mean: " + res + "?")
			fmt.Println("\nToo see all of Cfnctl's top-level commands, run\n\tcfnctl --help")
		}
	}
	app.Commands = []*cli.Command{
		{
			Name:  "apply",
			Usage: "Create or update infrastructure",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "auto-approve",
					Usage: "Skip interactive approval of plan before applying.",
				},
				&cli.StringFlag{
					Name:     "template-file",
					Usage:    "The path of the Cloudformation template you wish to apply",
					Required: true,
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
			Action: func(c *cli.Context) error {
				apply := Apply{
					TemplatePath: c.String("template-file"),
					ParamFile:    c.String("param-file"),
					AutoApprove:  c.Bool("auto-approve"),
				}
				err := apply.Run()
				return err
			},
		},
		{
			Name:  "plan",
			Usage: "Show changes required by the current configuration",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "template-file",
					Usage:    "filename",
					Required: true,
				},
				&cli.StringFlag{
					Name:  "param-file",
					Usage: "filename",
				},
			},
			Action: func(c *cli.Context) error {
				plan := Plan{
					TemplatePath: c.String("template-file"),
					ParamFile:    c.String("param-file"),
				}
				err := plan.Run()
				return err
			},
		},
		{
			Name:  "destroy",
			Usage: "Destroy previously-created infrastructure",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "template-file",
					Usage:    "filename",
					Required: true,
				},
				&cli.BoolFlag{
					Name:  "auto-approve",
					Usage: "Skip interactive approval of plan before applying.",
					Value: false,
				},
			},
			Action: func(c *cli.Context) error {
				destroy := Destroy{
					AutoApprove:  c.Bool("auto-approve"),
					TemplatePath: c.String("template-file"),
				}
				err := destroy.Run()
				return err
			},
		},
		{
			Name:  "output",
			Usage: "Show all exported output values of the selected account and region",
			Action: func(c *cli.Context) error {
				out := Output{}
				err := out.Run()
				return err
			},
		},
		{
			Name:  "validate",
			Usage: "Check whether the configuration is valid",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "template-file",
					Usage:    "The path of the Cloudformation template you wish to apply",
					Required: true,
				},
			},
			Action: func(c *cli.Context) error {
				v := Validate{TemplatePath: c.String("template-file")}
				err := v.Run()
				return err
			},
		},
		{
			Name:  "version",
			Usage: "Show the current Cfnctl version",
			Action: func(c *cli.Context) error {
				v := Version{Version: version}
				err := v.Run()
				return err
			},
		},
	}

	err := app.Run(args)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

func setCustomCLITemplate(c *cli.App) {
	whiteBold := color.New(color.Bold).SprintfFunc()
	whiteUnderline := color.New(color.Bold).Add(color.Underline).SprintfFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	c.CustomAppHelpTemplate = fmt.Sprintf(` %s:
		{{.Name}}{{if .Usage}} - {{.Usage}}{{end}}{{if .Description}}

	 DESCRIPTION:
		{{.Description | nindent 3 | trim}}{{end}}{{if len .Authors}}

	 AUTHOR{{with $length := len .Authors}}{{if ne 1 $length}}S{{end}}{{end}}:
		{{range $index, $author := .Authors}}{{if $index}}
		{{end}}{{$author}}{{end}}{{end}}{{if .VisibleCommands}}

 %s:{{range .VisibleCategories}}{{if .Name}}
	{{.Name}}:{{range .VisibleCommands}}
	  {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}{{else}}{{range .VisibleCommands}}
	{{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}{{end}}{{end}}{{end}}{{if .VisibleFlags}}

 %s:
	{{range $index, $option := .VisibleFlags}}{{if $index}}
	{{end}}{{$option}}{{end}}{{end}}{{if .Copyright}}

 COPYRIGHT:
	{{.Copyright}}{{end}}
		
 %s
  Apply infrastructure using the "apply" command.
    %s
`, whiteBold("NAME"),
		whiteBold("COMMANDS"),
		whiteBold("GLOBAL OPTIONS"),
		whiteUnderline("Examples"),
		cyan("$ cfnctl apply --template-file mycfntmpl.yaml --auto-approve"))
}
