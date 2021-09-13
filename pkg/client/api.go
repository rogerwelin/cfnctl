package client

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/rogerwelin/cfnctl/internal/utils"
)

type Option func(*Cfnctl)

func WithAutoApprove(b bool) Option {
	return func(c *Cfnctl) {
		c.AutoApprove = b
	}
}

func WithVarsFile(varsFile string) Option {
	return func(c *Cfnctl) {
		c.VarsFile = varsFile
	}
}

func WithStackName(stackName string) Option {
	return func(c *Cfnctl) {
		c.StackName = stackName
	}
}

func WithChangesetName(changesetName string) Option {
	return func(c *Cfnctl) {
		c.ChangesetName = changesetName
	}
}
func WithSvc(svc CloudformationAPI) Option {
	return func(c *Cfnctl) {
		c.Svc = svc
	}
}

func WithTemplatePath(filePath string) Option {
	return func(c *Cfnctl) {
		c.TemplatePath = filePath
	}
}

func WithTemplateBody(file string) Option {
	return func(c *Cfnctl) {
		c.TemplateBody = file
	}
}

func WithOutput(out io.Writer) Option {
	return func(c *Cfnctl) {
		c.Output = out
	}
}

func New(option ...Option) *Cfnctl {
	c := &Cfnctl{}
	for _, o := range option {
		o(c)
	}
	return c
}

/*
func New(autoApprove bool, varsFile, stackName, changesetName string, svc CloudformationAPI) *Cfnctl {
	return &Cfnctl{
		AutoApprove:   autoApprove,
		VarsFile:      varsFile,
		StackName:     stackName,
		ChangesetName: changesetName,
		Svc:           svc,
	}
}
*/

func (c *Cfnctl) ApplyChangeSet() error {
	input := &cloudformation.ExecuteChangeSetInput{
		ChangeSetName: &c.StackName,
		StackName:     &c.StackName,
	}

	_, err := c.Svc.ExecuteChangeSet(context.TODO(), input)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cfnctl) IsStackCreated() (bool, error) {
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
		if *val.StackName == c.StackName {
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
func (c *Cfnctl) ListChangeSet() (types.ChangeSetStatus, error) {
	input := &cloudformation.ListChangeSetsInput{
		StackName: &c.StackName,
	}

	output, err := c.Svc.ListChangeSets(context.TODO(), input)
	if err != nil {
		return "", err
	}

	return output.Summaries[0].Status, nil
}

func (c *Cfnctl) DescribeChangeSet() ([]types.Change, error) {
	input := &cloudformation.DescribeChangeSetInput{
		ChangeSetName: &c.ChangesetName,
		StackName:     &c.StackName,
	}

	out, err := c.Svc.DescribeChangeSet(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	return out.Changes, nil
}

func (c *Cfnctl) CreateChangeSet() error {

	capabilities := []types.Capability{"CAPABILITY_NAMED_IAM"}
	var changeSetType types.ChangeSetType

	// 1. check of stack already exists. if so choose UPDATE. if not choose CREATE
	// 2. if stack already exists choose new change set name
	created, err := c.IsStackCreated()
	if err != nil {
		return err
	}

	if created {
		changeSetType = "UPDATE"
		found, err := c.ChangeSetExists(c.ChangesetName, c.StackName)
		if err != nil {
			return err
		}
		if found {
			suffix := utils.ReturnRandom(5)
			c.ChangesetName = c.ChangesetName + "-" + suffix
			//changesetName = changesetName + "-" + suffix
		}
	} else {
		changeSetType = "CREATE"
	}

	input := &cloudformation.CreateChangeSetInput{
		ChangeSetName: &c.ChangesetName,
		StackName:     &c.StackName,
		ChangeSetType: changeSetType,
		TemplateBody:  &c.TemplateBody,
		Capabilities:  capabilities,
	}

	if c.Parameters != nil {
		input.Parameters = c.Parameters
	}

	_, err = c.Svc.CreateChangeSet(context.TODO(), input)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cfnctl) DeleteChangeSet() error {
	input := &cloudformation.DeleteChangeSetInput{
		ChangeSetName: &c.ChangesetName,
		StackName:     &c.StackName,
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

func (c *Cfnctl) DescribeStack() (string, error) {
	input := &cloudformation.DescribeStacksInput{
		StackName: &c.StackName,
	}
	out, err := c.Svc.DescribeStacks(context.TODO(), input)
	if err != nil {
		return "", err
	}

	return string(out.Stacks[0].StackStatus), nil
}

func (c *Cfnctl) DescribeStackEvents() error {
	input := &cloudformation.DescribeStackEventsInput{
		StackName: &c.StackName,
	}

	out, err := c.Svc.DescribeStackEvents(context.TODO(), input)
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

func (c *Cfnctl) DescribeStackResources() ([]types.StackResource, error) {
	input := &cloudformation.DescribeStackResourcesInput{
		StackName: &c.StackName,
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
