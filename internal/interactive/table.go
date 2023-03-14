package interactive

import (
	"io"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/buger/goterm"
	"github.com/olekukonko/tablewriter"
)

type StackResourceEvents struct {
	Events []types.StackResource
}

func TableOutputter(events []types.StackResource, writer io.Writer) {
	if events == nil {
		return
	}

	if len(events) == 0 {
		return
	}

	tableData := [][]string{}
	table := tablewriter.NewWriter(writer)
	table.SetHeader([]string{"Logical ID", "Physical ID", "Type", "Status", "Status Reason"})
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
	)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_RIGHT)
	table.SetAlignment(tablewriter.ALIGN_RIGHT)

	goterm.MoveCursor(1, 1)

	for _, item := range events {
		var physicalID string
		var statusReason string
		var logicalResourceId string
		var ResourceType string

		if item.PhysicalResourceId != nil {
			physicalID = *item.PhysicalResourceId
		} else {
			physicalID = "-"
		}

		if item.ResourceStatusReason != nil {
			statusReason = *item.ResourceStatusReason
		} else {
			statusReason = "-"
		}

		if item.LogicalResourceId != nil {
			logicalResourceId = *item.LogicalResourceId
		} else {
			logicalResourceId = "-"
		}

		if item.ResourceType != nil {
			ResourceType = *item.ResourceType
		} else {
			ResourceType = "-"
		}

		arr := []string{
			logicalResourceId,
			physicalID,
			ResourceType,
			string(item.ResourceStatus),
			statusReason,
		}
		tableData = append(tableData, arr)
	}

	for i := range tableData {
		switch tableData[i][3] {
		case "CREATE_COMPLETE":
			table.Rich(tableData[i], []tablewriter.Colors{{}, {}, {}, {tablewriter.Normal, tablewriter.FgHiGreenColor}, {}})
		case "CREATE_IN_PROGRESS":
			table.Rich(tableData[i], []tablewriter.Colors{{}, {}, {}, {tablewriter.Normal, tablewriter.FgHiBlueColor}, {}})
		case "CREATE_FAILED":
			table.Rich(tableData[i], []tablewriter.Colors{{tablewriter.Normal, tablewriter.FgHiRedColor}})
		default:
			table.Append(tableData[i])
		}
	}

	table.Render()
	goterm.Flush()
}
