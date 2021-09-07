package main

import (
	"os"

	"github.com/rogerwelin/cfnctl/internal/cli"
)

func main() {
	cli.RunCLI(os.Args)
}
