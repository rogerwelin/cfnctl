package client

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
)

// CloudformationAPI provides access to AWS CloudFormation API
type CloudformationAPI interface {
	ExecuteChangeSet(ctx context.Context, params *cloudformation.ExecuteChangeSetInput, optFns ...func(*cloudformation.Options)) (*cloudformation.ExecuteChangeSetOutput, error)
	CreateChangeSet(ctx context.Context, params *cloudformation.CreateChangeSetInput, optFns ...func(*cloudformation.Options)) (*cloudformation.CreateChangeSetOutput, error)
	DescribeChangeSet(ctx context.Context, params *cloudformation.DescribeChangeSetInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DescribeChangeSetOutput, error)
	DeleteChangeSet(ctx context.Context, params *cloudformation.DeleteChangeSetInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DeleteChangeSetOutput, error)
	DescribeStacks(ctx context.Context, params *cloudformation.DescribeStacksInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DescribeStacksOutput, error)
	DescribeStackResources(ctx context.Context, params *cloudformation.DescribeStackResourcesInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DescribeStackResourcesOutput, error)
	ListChangeSets(ctx context.Context, params *cloudformation.ListChangeSetsInput, optFns ...func(*cloudformation.Options)) (*cloudformation.ListChangeSetsOutput, error)
	ListStacks(ctx context.Context, params *cloudformation.ListStacksInput, optFns ...func(*cloudformation.Options)) (*cloudformation.ListStacksOutput, error)
	ValidateTemplate(ctx context.Context, params *cloudformation.ValidateTemplateInput, optFns ...func(*cloudformation.Options)) (*cloudformation.ValidateTemplateOutput, error)
	ListExports(ctx context.Context, params *cloudformation.ListExportsInput, optFns ...func(*cloudformation.Options)) (*cloudformation.ListExportsOutput, error)
	DeleteStack(ctx context.Context, params *cloudformation.DeleteStackInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DeleteStackOutput, error)
	DetectStackDrift(ctx context.Context, params *cloudformation.DetectStackDriftInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DetectStackDriftOutput, error)
	DescribeStackDriftDetectionStatus(ctx context.Context, params *cloudformation.DescribeStackDriftDetectionStatusInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DescribeStackDriftDetectionStatusOutput, error)
	DescribeStackResourceDrifts(ctx context.Context, params *cloudformation.DescribeStackResourceDriftsInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DescribeStackResourceDriftsOutput, error)
}

// Cfnctl provides access to all actions in the programs lifecycle
type Cfnctl struct {
	AutoApprove       bool
	VarsFile          string
	StackName         string
	ChangesetName     string
	TemplateBody      string
	TemplatePath      string
	Parameters        []types.Parameter
	Output            io.Writer
	Svc               CloudformationAPI
	DoDeleteChangeset bool
}
