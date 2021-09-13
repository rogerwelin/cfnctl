package commands

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/rogerwelin/cfnctl/internal/aws"
	"github.com/rogerwelin/cfnctl/internal/params"
	"github.com/rogerwelin/cfnctl/internal/utils"
	"github.com/rogerwelin/cfnctl/pkg/client"
)

/*
	Vingtor rise to face
	The snake with hammer high
	At the edge of the world
	Bolts of lightning fills the air
	As Mjölner does its work
	The dreadful serpent roars in pain
*/

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

func Plan(template, varsfile string, deleteChangeSet bool) error {

	var deleteStack bool

	svc, err := aws.NewAWS()
	if err != nil {
		return err
	}

	templateBody, err := utils.ReadFile(template)
	if err != nil {
		return err
	}

	stackName := utils.TrimFileSuffix(template)

	ctl := &client.Cfnctl{
		Svc:           svc,
		TemplateBody:  string(templateBody),
		StackName:     stackName,
		ChangesetName: stackName,
		Output:        os.Stdout,
	}

	// if vars file is supplied
	if varsfile != "" {
		out, err := params.MergeFileParams(varsfile)
		ctl.Parameters = out
		if err != nil {
			return err
		}
		err = ctl.CreateChangeSet(ctl.TemplateBody, ctl.StackName, ctl.ChangesetName)
		if err != nil {
			return err
		}
	} else {
		// no vars file. check if tempalte contains params
		ok, outParams, err := params.CheckInputParams(template)
		if err != nil {
			return err
		}
		// no input params or default value set
		if !ok {
			// create change set
			err = ctl.CreateChangeSet(ctl.TemplateBody, ctl.StackName, ctl.ChangesetName)
			if err != nil {
				return err
			}
		} else {
			// get user input
			out, err := params.BuildInputParams(outParams)
			if err != nil {
				return err
			}
			ctl.Parameters = out
			err = ctl.CreateChangeSet(ctl.TemplateBody, ctl.StackName, ctl.ChangesetName)
			if err != nil {
				return err
			}
		}
	}

	// needs to be improved
	count := 15
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

	_ = planOutput(createEvents, ctl.Output)

	// clean up changeset
	if deleteChangeSet {

		err = ctl.DeleteChangeSet(ctl.StackName, ctl.ChangesetName)
		if err != nil {
			return err
		}

		if deleteStack {
			err := ctl.DeleteStack(ctl.StackName)
			if err != nil {
				return err
			}
		}

	}
	return nil
}
