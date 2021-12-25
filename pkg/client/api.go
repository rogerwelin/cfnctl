package client

import (
	"context"
	"errors"
	"io"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/rogerwelin/cfnctl/utils"
)

// Option is used to implement Option Pattern on the client
type Option func(*Cfnctl)

// WithAutoApprove provides access creating a new client
func WithAutoApprove(b bool) Option {
	return func(c *Cfnctl) {
		c.AutoApprove = b
	}
}

// WithVarsFile provides access creating a new client
func WithVarsFile(varsFile string) Option {
	return func(c *Cfnctl) {
		c.VarsFile = varsFile
	}
}

// WithStackName provides access creating a new client
func WithStackName(stackName string) Option {
	return func(c *Cfnctl) {
		c.StackName = stackName
	}
}

// WithChangesetName provides access creating a new client
func WithChangesetName(changesetName string) Option {
	return func(c *Cfnctl) {
		c.ChangesetName = changesetName
	}
}

// WithSvc provides access creating a new client
func WithSvc(svc CloudformationAPI) Option {
	return func(c *Cfnctl) {
		c.Svc = svc
	}
}

// WithTemplatePath provides access creating a new client
func WithTemplatePath(filePath string) Option {
	return func(c *Cfnctl) {
		c.TemplatePath = filePath
	}
}

// WithTemplateBody provides access creating a new client
func WithTemplateBody(file string) Option {
	return func(c *Cfnctl) {
		c.TemplateBody = file
	}
}

// WithOutput provides access creating a new client
func WithOutput(out io.Writer) Option {
	return func(c *Cfnctl) {
		c.Output = out
	}
}

// New creates a new client
func New(option ...Option) *Cfnctl {
	c := &Cfnctl{}
	for _, o := range option {
		o(c)
	}
	return c
}

// ApplyChangeSet executes a CloudFormation changeset
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

// IsStackCreated checks whether a CloudFormation stack is created or not
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

// ChangeSetExists checks whether a CloudFormation changeset exists or not
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

// ListChangeSet lists all changesets. Can be used to get status if changeset created or not
func (c *Cfnctl) ListChangeSet() (types.ChangeSetStatus, error) {
	input := &cloudformation.ListChangeSetsInput{
		StackName: &c.StackName,
	}

	output, err := c.Svc.ListChangeSets(context.TODO(), input)
	if err != nil {
		return "", err
	}

	if len(output.Summaries) == 0 {
		return "", errors.New("empty resultset when listing change sets")
	}

	return output.Summaries[0].Status, nil
}

// DescribeChangeSet describes a CloudFormation changeset
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

// CreateChangeSet creates a CloudFormation changeset
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

// DeleteChangeSet deletes a CloudFormation changeset
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

// DescribeStack describes a CloudFormation stack
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

// DescribeStackResources describes the resources from a particular CloudFormation stack
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

// ValidateCFTemplate validates a particular CloudFormation template
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

func (c *Cfnctl) ListExportValues() ([]types.Export, error) {
	input := &cloudformation.ListExportsInput{}
	out, err := c.Svc.ListExports(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	/*
		for _, item := range out.Exports {
			fmt.Printf("%s : %s : %s\n", *item.ExportingStackId, *item.Name, *item.Value)
		}
	*/

	return out.Exports, nil

}
