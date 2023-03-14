package commands

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/buger/goterm"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/rogerwelin/cfnctl/pkg/client"
)

type stackResourceEvents struct {
	events []types.StackResource
}

func tableOutputter(events []types.StackResource, writer io.Writer) {
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

func streamStackResources(ch <-chan stackResourceEvents, done <-chan bool) {
	for {
		select {
		case <-done:
			return
		case item := <-ch:
			tableOutputter(item.events, os.Stdout)
		}
	}
}

// Apply executes a given CF template
func Apply(ctl *client.Cfnctl) error {
	greenBold := color.New(color.Bold, color.FgHiGreen).SprintFunc()
	whiteBold := color.New(color.Bold).SprintfFunc()

	eventsChan := make(chan stackResourceEvents)
	doneChan := make(chan bool)

	pc, err := Plan(ctl, false)
	if err != nil {
		return err
	}

	if !pc.containsChanges {
		fmt.Fprintf(ctl.Output, "\n%s. %s\n\n", greenBold("No changes"), whiteBold("Your infrastructure matches the configuration"))
		fmt.Fprintf(ctl.Output, "Cfnctl has compared your real infrastructure against your configuration and found no differences, so no changes are needed.\n")
		fmt.Fprintf(ctl.Output, "\n%s %d added, %d changed, %d destroyed\n", greenBold("Apply complete! Resources:"), (pc.changes["add"]), pc.changes["change"], pc.changes["destroy"])
		return nil
	}

	if !ctl.AutoApprove {
		reader := bufio.NewReader(os.Stdin)

		fmt.Fprintf(ctl.Output, "%s\n"+
			"  Cfnctl will perform the actions described above.\n"+
			"  Only 'yes' will be accepted to approve.\n\n"+
			"  %s", whiteBold("Do you want to perform the following actions?"), whiteBold("Enter a value: "))

		choice, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		choice = strings.TrimSuffix(choice, "\n")

		if choice != "yes" {
			fmt.Fprintf(ctl.Output, "\nApply cancelled.\n")
			return nil
		}
	}

	goterm.Clear()

	err = ctl.ApplyChangeSet()
	if err != nil {
		return err
	}

	go streamStackResources(eventsChan, doneChan)

	// to be improved
	for {
		time.Sleep(900 * time.Millisecond)
		status, err := ctl.DescribeStack()
		if err != nil {
			return err
		}
		if status == "UPDATE_COMPLETE" || status == "CREATE_FAILED" || status == "CREATE_COMPLETE" {
			break
		} else {
			event, err := ctl.DescribeStackResources()
			if err != nil {
				return err
			}
			eventsChan <- stackResourceEvents{events: event}
		}
	}

	close(eventsChan)
	doneChan <- true

	// this is a really dirty hack
	// insert newlines so table does not dissapear
	numAdd := pc.changes["add"]
	numChange := pc.changes["change"]
	numDestroy := pc.changes["destroy"]
	total := numAdd + numChange + numDestroy + 4 // for header and padding
	//lint:ignore SA1006 I know what i'm doing
	fmt.Printf(strings.Repeat("\n", total))
	fmt.Fprintf(ctl.Output, "\n%s %d added, %d changed, %d destroyed\n", greenBold("Apply complete! Resources:"), (pc.changes["add"]), pc.changes["change"], pc.changes["destroy"])

	return nil
}
