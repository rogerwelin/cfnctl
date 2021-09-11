package commands

import (
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/gosuri/uilive"
	"github.com/olekukonko/tablewriter"
)

/*
	There comes Fenris' twin
	His jaws are open wide
	The serpent rises from the waves
	Jormungandr twists and turns
	Mighty in his wrath
	The eyes are full of primal hate
*/

type stackResourceEvents struct {
	events []types.StackResource
}

func tableOutputter(events []types.StackResource, writer io.Writer) {
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

	if len(events) == 0 {
		return
	}

	for _, item := range events {
		var physicalID string
		var statusReason string

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

		arr := []string{
			*item.LogicalResourceId,
			physicalID,
			*item.ResourceType,
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
}

func streamStackResources(ch <-chan stackResourceEvents, done <-chan bool) {
	writer := uilive.New()
	writer.Start()

	for {
		select {
		case <-done:
			fmt.Println("we've ended")
			writer.Stop()
			return
		case item := <-ch:
			tableOutputter(item.events, writer)
		}
	}
}

func Apply(template, varsfile string) error {

	return nil
}
