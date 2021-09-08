package client

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
)

func New(autoApprove bool, varsFile, stackName, changesetName string, svc CloudformationAPI) *Cfnctl {
	return &Cfnctl{
		AutoApprove:   autoApprove,
		VarsFile:      varsFile,
		StackName:     stackName,
		ChangesetName: changesetName,
		Svc:           svc,
	}
}

func returnRandom(value int) string {
	stringArr := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "x", "y", "z", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T"}
	newString := ""

	for i := 0; i <= value; i++ {
		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)
		randIndex := r1.Intn(len(stringArr))
		newString = newString + stringArr[randIndex]
	}
	return newString
}

func (c *Cfnctl) ApplyChangeSet(stackName string) error {
	input := &cloudformation.ExecuteChangeSetInput{
		ChangeSetName: &stackName,
		StackName:     &stackName,
	}

	_, err := c.Svc.ExecuteChangeSet(context.TODO(), input)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cfnctl) IsStackCreated(stackName string) (bool, error) {
	input := &cloudformation.ListStacksInput{
		StackStatusFilter: []types.StackStatus{
			"CREATE_COMPLETE",
		},
	}

	out, err := c.Svc.ListStacks(context.TODO(), input)
	if err != nil {
		return false, nil
	}

	for _, val := range out.StackSummaries {
		if *val.StackName == stackName {
			return true, nil
		}
	}
	return false, nil
}

func (c *Cfnctl) ChangeSetExists(stackName, changesetName string) (bool, error) {
	input := &cloudformation.ListChangeSetsInput{
		StackName: &stackName,
	}
	out, err := c.Svc.ListChangeSets(context.TODO(), input)
	if err != nil {
		return false, err
	}
	found := false

	for _, val := range out.Summaries {
		if *val.ChangeSetName == changesetName {
			found = true
			break
		}
	}
	return found, nil
}

// use this one to get status / if created or not
func (c *Cfnctl) ListChangeSet(stackName string) (types.ChangeSetStatus, error) {
	input := &cloudformation.ListChangeSetsInput{
		StackName: &stackName,
	}

	output, err := c.Svc.ListChangeSets(context.TODO(), input)
	if err != nil {
		return "", err
	}

	return output.Summaries[0].Status, nil
}

func (c *Cfnctl) DescribeChangeSet(stackName, changesetName string) ([]types.Change, error) {
	input := &cloudformation.DescribeChangeSetInput{
		ChangeSetName: &changesetName,
		StackName:     &stackName,
	}

	out, err := c.Svc.DescribeChangeSet(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	return out.Changes, nil
}

func (c *Cfnctl) CreateChangeSet(tBody, stackName, changesetName string, params []types.Parameter) error {

	capabilities := []types.Capability{"CAPABILITY_NAMED_IAM"}
	var changeSetType types.ChangeSetType

	// 1. check of stack already exists. if so choose UPDATE. if not choose CREATE
	// 2. if stack already exists choose new change set name
	created, err := c.IsStackCreated(stackName)
	if err != nil {
		return err
	}

	if created {
		changeSetType = "UPDATE"
		found, err := c.ChangeSetExists(changesetName, stackName)
		if err != nil {
			return err
		}
		if found {
			suffix := returnRandom(5)
			c.ChangesetName = changesetName + "-" + suffix
			changesetName = changesetName + "-" + suffix
		}
	} else {
		changeSetType = "CREATE"
	}

	input := &cloudformation.CreateChangeSetInput{
		ChangeSetName: &changesetName,
		StackName:     &stackName,
		ChangeSetType: changeSetType,
		Parameters:    params,
		TemplateBody:  &tBody,
		Capabilities:  capabilities,
	}

	_, err = c.Svc.CreateChangeSet(context.TODO(), input)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cfnctl) DeleteChangeSet(stackName, changesetName string) error {
	input := &cloudformation.DeleteChangeSetInput{
		ChangeSetName: &changesetName,
		StackName:     &stackName,
	}

	_, err := c.Svc.DeleteChangeSet(context.TODO(), input)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cfnctl) DeleteStack(stackName string) error {
	input := &cloudformation.DeleteStackInput{
		StackName: &stackName,
	}

	_, err := c.Svc.DeleteStack(context.TODO(), input)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cfnctl) DescribeStack(stackName string) (string, error) {
	input := &cloudformation.DescribeStacksInput{
		StackName: &stackName,
	}
	out, err := c.Svc.DescribeStacks(context.TODO(), input)
	if err != nil {
		return "", err
	}
	//fmt.Printf("%+v\n", out)
	//fmt.Println(out)

	return string(out.Stacks[0].StackStatus), nil
}

func (c *Cfnctl) DescribeStackEvents(stackName string) error {
	input := &cloudformation.DescribeStackEventsInput{
		StackName: &stackName,
	}

	out, err := c.Svc.DescribeStackEvents(context.TODO(), input)
	// fmt.Println(out.StackEvents)
	for _, item := range out.StackEvents {
		fmt.Print(*item.LogicalResourceId)
		fmt.Print(" : ")
		fmt.Print(*item.PhysicalResourceId)
		fmt.Print(" : ")
		fmt.Print(*item.ResourceType)
		fmt.Println("")
	}
	if err != nil {
		return err
	}
	return nil
}

func (c *Cfnctl) DescribeStackResources(stackName string) ([]types.StackResource, error) {
	input := &cloudformation.DescribeStackResourcesInput{
		StackName: &stackName,
	}

	out, err := c.Svc.DescribeStackResources(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	return out.StackResources, nil
}

func (c *Cfnctl) ValidateCFTemplate() error {
	input := &cloudformation.ValidateTemplateInput{
		TemplateBody: &c.TemplateBody,
	}
	_, err := c.Svc.ValidateTemplate(context.TODO(), input)
	if err != nil {
		return err
	}
	return nil
}
