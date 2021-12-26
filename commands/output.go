package commands

import (
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/olekukonko/tablewriter"
	"github.com/rogerwelin/cfnctl/pkg/client"
)

func outputExportTable(values []types.Export, writer io.Writer) {
	tableData := [][]string{}
	table := tablewriter.NewWriter(writer)
	table.SetHeader([]string{"Name", "Value"})
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
	)

	for _, item := range values {
		arr := []string{
			*item.Name,
			*item.Value,
		}
		tableData = append(tableData, arr)
	}

	for _, item := range tableData {
		table.Append(item)
	}

	fmt.Printf("\n")
	table.Render()
	fmt.Printf("\n")
}

func Output(ctl *client.Cfnctl) error {
	out, err := ctl.ListExportValues()

	if err != nil {
		return nil
	}

	if len(out) == 0 {
		fmt.Fprintf(ctl.Output, "No exported values in the selected region")
	}

	outputExportTable(out, ctl.Output)

	return nil
}
