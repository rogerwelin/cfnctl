package commands

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/rogerwelin/cfnctl/pkg/client"
)

func destroytOutput(input []types.StackResource, writer io.Writer) int {

	tableData := [][]string{}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Action", "Logical ID", "Physical ID", "Resource type"})
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
	)

	for _, v := range input {
		arr := []string{
			"Destroy",
			*v.LogicalResourceId,
			*v.PhysicalResourceId,
			*v.ResourceType,
		}
		tableData = append(tableData, arr)
	}

	for i := range tableData {
		switch tableData[i][0] {
		case "Destroy":
			table.Rich(tableData[i], []tablewriter.Colors{{tablewriter.Normal, tablewriter.FgHiRedColor}})
		default:
			table.Append(tableData[i])
		}
	}

	fmt.Fprintf(writer, "\nCfnctl will perform the following actions:\n\n")

	table.Render()

	return len(tableData)
}

// Destroy destroys all the resources in a given stack
func Destroy(ctl *client.Cfnctl) error {
	whiteBold := color.New(color.Bold).SprintfFunc()
	greenBold := color.New(color.Bold, color.FgHiGreen).SprintFunc()

	// check if stack exists
	ok, err := ctl.IsStackCreated()
	if err != nil {
		return err
	}

	if !ok {
		fmt.Printf("\n%s %s\n\n", greenBold("No changes."), whiteBold("No objects need to be destroyed"))
		fmt.Printf("Either you have not created any objects yet, there is no Stack named %s or the existing objects were already deleted outside of Cfnctl\n\n", ctl.StackName)
		fmt.Printf("%s", greenBold("Destroy complete! Resources: 0 destroyed\n"))
		return nil
	}

	out, err := ctl.DescribeStackResources()
	if err != nil {
		return err
	}

	noChanges := destroytOutput(out, ctl.Output)

	if !ctl.AutoApprove {
		fmt.Fprintf(ctl.Output, "%s\n"+
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
			fmt.Fprintf(ctl.Output, "\nDestroy cancelled.\n")
			return nil
		}
	}

	fmt.Fprintf(ctl.Output, "\n%s %s %d %s\n", greenBold("Destroy complete!"), greenBold("Resources:"), noChanges, greenBold("destroyed"))

	return nil
}
