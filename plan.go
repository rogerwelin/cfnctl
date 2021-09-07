package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/awslabs/goformation/v4"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/rogerwelin/cfctl/pkg/client"
)

type planChanges struct {
	containsChanges bool
	changes         map[string]int
}

func planOutput(changes []types.Change, writer io.Writer) planChanges {
	tableData := [][]string{}
	table := tablewriter.NewWriter(writer)
	table.SetHeader([]string{"Action", "Logical ID", "Physical ID", "Resource type", "Replacement"})
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
	)

	actionMap := make(map[string]int)

	for _, v := range changes {
		var physicalID string
		var replacement string

		if v.ResourceChange.PhysicalResourceId != nil {
			physicalID = *v.ResourceChange.PhysicalResourceId
		} else {
			physicalID = "-"
		}

		if v.ResourceChange.Replacement != "" {
			replacement = string(v.ResourceChange.Replacement)
		} else {
			replacement = "-"
		}

		if v.ResourceChange.Action == "Add" {
			actionMap["add"] += 1
		} else if v.ResourceChange.Action == "Remove" {
			actionMap["destroy"] += 1
		} else if v.ResourceChange.Action == "Modify" {
			actionMap["change"] += 1
		} else {
		}

		arr := []string{
			string(v.ResourceChange.Action),
			*v.ResourceChange.LogicalResourceId,
			physicalID,
			*v.ResourceChange.ResourceType,
			replacement,
		}
		tableData = append(tableData, arr)
	}

	for i := range tableData {
		switch tableData[i][0] {
		case "Add":
			table.Rich(tableData[i], []tablewriter.Colors{{tablewriter.Normal, tablewriter.FgHiGreenColor}})
		case "Delete":
			table.Rich(tableData[i], []tablewriter.Colors{{tablewriter.Normal, tablewriter.FgHiRedColor}})
		case "Modify":
			table.Rich(tableData[i], []tablewriter.Colors{{tablewriter.Normal, tablewriter.FgHiYellowColor}})
		default:
			table.Append(tableData[i])
		}
	}

	whiteBold := color.New(color.Bold).SprintFunc()
	fmt.Print("\nCfnctl will perform the following actions:\n\n")

	var modifications bool

	if len(changes) != 0 {
		modifications = true
		table.Render()
	} else {
		modifications = false
	}

	fmt.Printf("\n%s %d to add, %d to change, %d to destroy\n\n", whiteBold("Plan:"), actionMap["add"], actionMap["change"], actionMap["destroy"])

	pc := planChanges{
		containsChanges: modifications,
		changes:         actionMap,
	}

	return pc
}

func cfnCtlPlan(ctl *client.Cfnctl) (planChanges, error) {

	// set output to stdout
	ctl.Output = os.Stdout

	// only if we pass in vars file
	template, err := goformation.Open("dynamolambda.yaml")
	if err != nil {
		log.Fatalf("There was an error processing the template: %s", err)
	}

	dat, err := ioutil.ReadFile("dynamolambda.yaml")
	if err != nil {
		panic(err)
	}

	templateToString := string(dat)

	outParams, err := mergeParameters(template.Parameters)
	if err != nil {
		panic(err)
	}

	err = ctl.CreateChangeSet(templateToString, ctl.StackName, ctl.ChangesetName, outParams)
	if err != nil {
		panic(err)
	}

	// needs to be improved
	count := 8
	for i := 0; i < count; i++ {
		time.Sleep(1 * time.Second)
		status, err := ctl.ListChangeSet(ctl.StackName)
		if err != nil {
			panic(err)
		}
		if status == "CREATE_COMPLETE" {
			break
		}
	}

	createEvents, err := ctl.DescribeChangeSet(ctl.StackName, ctl.ChangesetName)
	if err != nil {
		panic(err)
	}

	pc := planOutput(createEvents, ctl.Output)

	// clean up changeset
	/*
		err = ctl.DeleteChangeSet(ctl.StackName, ctl.ChangesetName)
		if err != nil {
			return err
		}
	*/
	return pc, nil
}
