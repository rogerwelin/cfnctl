package mock

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/smithy-go/middleware"
	"github.com/rogerwelin/cfnctl/pkg/client"
)

type mockAPI struct{}

// NewMockAPI returns a new instance of mockAPI
func NewMockAPI() client.CloudformationAPI {
	return mockAPI{}
}

// ExecuteChangeSet returns a mocked response
func (m mockAPI) ExecuteChangeSet(ctx context.Context, params *cloudformation.ExecuteChangeSetInput, optFns ...func(*cloudformation.Options)) (*cloudformation.ExecuteChangeSetOutput, error) {
	res := middleware.Metadata{}
	res.Set("result", "ok")

	return &cloudformation.ExecuteChangeSetOutput{
		ResultMetadata: res,
	}, nil
}

// ExecuteChangeSet returns a mocked response
func (m mockAPI) CreateChangeSet(ctx context.Context, params *cloudformation.CreateChangeSetInput, optFns ...func(*cloudformation.Options)) (*cloudformation.CreateChangeSetOutput, error) {

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
func (m mockAPI) DescribeChangeSet(ctx context.Context, params *cloudformation.DescribeChangeSetInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DescribeChangeSetOutput, error) {

	return &cloudformation.DescribeChangeSetOutput{
		ChangeSetName: params.ChangeSetName,
		StackName:     params.StackName,
		Status:        "CREATE_COMPLETE",
		// Changes: ,
	}, nil
}

// DeleteChangeSet returns a mocked response
func (m mockAPI) DeleteChangeSet(ctx context.Context, params *cloudformation.DeleteChangeSetInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DeleteChangeSetOutput, error) {
	res := middleware.Metadata{}
	res.Set("result", "ok")
	return &cloudformation.DeleteChangeSetOutput{
		ResultMetadata: res,
	}, nil
}

// DescribeStacks returns a mocked response
func (m mockAPI) DescribeStacks(ctx context.Context, params *cloudformation.DescribeStacksInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DescribeStacksOutput, error) {
	return &cloudformation.DescribeStacksOutput{
		// Stacks: ,
	}, nil
}

// DescribeStackResources returns a mocked response
func (m mockAPI) DescribeStackResources(ctx context.Context, params *cloudformation.DescribeStackResourcesInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DescribeStackResourcesOutput, error) {
	return &cloudformation.DescribeStackResourcesOutput{}, nil
}

// ListChangeSets returns a mocked response
func (m mockAPI) ListChangeSets(ctx context.Context, params *cloudformation.ListChangeSetsInput, optFns ...func(*cloudformation.Options)) (*cloudformation.ListChangeSetsOutput, error) {
	return &cloudformation.ListChangeSetsOutput{}, nil
}

// ListStacks returns a mocked response
func (m mockAPI) ListStacks(ctx context.Context, params *cloudformation.ListStacksInput, optFns ...func(*cloudformation.Options)) (*cloudformation.ListStacksOutput, error) {
	return &cloudformation.ListStacksOutput{}, nil
}

// ValidateTemplate returns a mocked response
func (m mockAPI) ValidateTemplate(ctx context.Context, params *cloudformation.ValidateTemplateInput, optFns ...func(*cloudformation.Options)) (*cloudformation.ValidateTemplateOutput, error) {
	return &cloudformation.ValidateTemplateOutput{}, nil
}
