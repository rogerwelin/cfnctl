package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/rogerwelin/cfnctl/pkg/client"
)

func Destroy(ctl *client.Cfnctl) error {
	whiteBold := color.New(color.Bold).SprintfFunc()

	fmt.Printf("%s\n"+
		"  Cfnctl will destroy all your managed infrastructure, as shown above\n"+
		"  There is no undo. Only 'yes' will be accepted to approve.\n\n"+
		"  %s", whiteBold("Do you really want to destroy all resources?"), whiteBold("Enter a value: "))

	reader := bufio.NewReader(os.Stdin)

	choice, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	choice = strings.TrimSuffix(choice, "\n")

	if choice != "yes" {
		fmt.Println("\nDestroy cancelled.")
		return nil
	}

	return nil
}
