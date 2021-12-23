package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
)

// NewAWS returns a new cloudformation client
func NewAWS() (*cloudformation.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	// if we pass in a profile
	// apa, err := config.LoadSharedConfigProfile()

	svc := cloudformation.NewFromConfig(cfg)

	return svc, nil
}
