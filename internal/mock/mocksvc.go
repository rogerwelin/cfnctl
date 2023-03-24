package mock

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/smithy-go/middleware"
)

type MockAPI struct{}

// NewMockAPI returns a new instance of mockAPI
func NewMockAPI() MockAPI {
	return MockAPI{}
}

// ExecuteChangeSet returns a mocked response
func (m MockAPI) ExecuteChangeSet(ctx context.Context, params *cloudformation.ExecuteChangeSetInput, optFns ...func(*cloudformation.Options)) (*cloudformation.ExecuteChangeSetOutput, error) {

	res := middleware.Metadata{}
	res.Set("result", "ok")

	return &cloudformation.ExecuteChangeSetOutput{
		ResultMetadata: res,
	}, nil
}

// ExecuteChangeSet returns a mocked response
func (m MockAPI) CreateChangeSet(ctx context.Context, params *cloudformation.CreateChangeSetInput, optFns ...func(*cloudformation.Options)) (*cloudformation.CreateChangeSetOutput, error) {

	id := "apa"
	stackID := "123456"
	res := middleware.Metadata{}
	res.Set("Status", "change set created")

	return &cloudformation.CreateChangeSetOutput{
		Id:             &id,
		StackId:        &stackID,
		ResultMetadata: res,
	}, nil
}

// DescribeChangeSet returns a mocked response
func (m MockAPI) DescribeChangeSet(ctx context.Context, params *cloudformation.DescribeChangeSetInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DescribeChangeSetOutput, error) {

	return &cloudformation.DescribeChangeSetOutput{
		ChangeSetName: params.ChangeSetName,
		StackName:     params.StackName,
		Status:        "CREATE_COMPLETE",
		// Changes: ,
	}, nil
}

// DeleteChangeSet returns a mocked response
func (m MockAPI) DeleteChangeSet(ctx context.Context, params *cloudformation.DeleteChangeSetInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DeleteChangeSetOutput, error) {
	res := middleware.Metadata{}
	res.Set("result", "ok")
	return &cloudformation.DeleteChangeSetOutput{
		ResultMetadata: res,
	}, nil
}

// DescribeStacks returns a mocked response
func (m MockAPI) DescribeStacks(ctx context.Context, params *cloudformation.DescribeStacksInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DescribeStacksOutput, error) {
	return &cloudformation.DescribeStacksOutput{
		// Stacks: ,
	}, nil
}

// DescribeStackResources returns a mocked response
func (m MockAPI) DescribeStackResources(ctx context.Context, params *cloudformation.DescribeStackResourcesInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DescribeStackResourcesOutput, error) {
	return &cloudformation.DescribeStackResourcesOutput{}, nil
}

// ListChangeSets returns a mocked response
func (m MockAPI) ListChangeSets(ctx context.Context, params *cloudformation.ListChangeSetsInput, optFns ...func(*cloudformation.Options)) (*cloudformation.ListChangeSetsOutput, error) {
	status := types.ChangeSetSummary{Status: "CREATE_COMPLETE"}
	sum := []types.ChangeSetSummary{status}

	return &cloudformation.ListChangeSetsOutput{Summaries: sum}, nil
}

// ListStacks returns a mocked response
func (m MockAPI) ListStacks(ctx context.Context, params *cloudformation.ListStacksInput, optFns ...func(*cloudformation.Options)) (*cloudformation.ListStacksOutput, error) {
	return &cloudformation.ListStacksOutput{}, nil
}

// ValidateTemplate returns a mocked response
func (m MockAPI) ValidateTemplate(ctx context.Context, params *cloudformation.ValidateTemplateInput, optFns ...func(*cloudformation.Options)) (*cloudformation.ValidateTemplateOutput, error) {
	return &cloudformation.ValidateTemplateOutput{}, nil
}

// ListExports returns a mocked response
func (m MockAPI) ListExports(ctx context.Context, params *cloudformation.ListExportsInput, optFns ...func(*cloudformation.Options)) (*cloudformation.ListExportsOutput, error) {
	return &cloudformation.ListExportsOutput{
		Exports: []types.Export{
			{
				ExportingStackId: aws.String("template"),
				Name:             aws.String("Bucket"),
				Value:            aws.String("TestBucket"),
			},
		},
	}, nil
}

// DeleteStack returns a mocked response
func (m MockAPI) DeleteStack(ctx context.Context, params *cloudformation.DeleteStackInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DeleteStackOutput, error) {
	return &cloudformation.DeleteStackOutput{}, nil
}

// DetectStackDrift returns a mocked response
func (m MockAPI) DetectStackDrift(ctx context.Context, params *cloudformation.DetectStackDriftInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DetectStackDriftOutput, error) {
	return nil, nil
}

// DescribeStackDriftDetectionStatus returns a mocked response
func (m MockAPI) DescribeStackDriftDetectionStatus(ctx context.Context, params *cloudformation.DescribeStackDriftDetectionStatusInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DescribeStackDriftDetectionStatusOutput, error) {
	return nil, nil
}
