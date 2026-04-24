package drift

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
)

// CloudFormationDetector checks a CloudFormation stack for drift.
type CloudFormationDetector struct {
	client *cloudformation.Client
}

// NewCloudFormationDetector creates a new detector using the given region and profile.
func NewCloudFormationDetector(region, profile string) (*CloudFormationDetector, error) {
	opts := []func(*config.LoadOptions) error{
		config.WithRegion(region),
	}
	if profile != "" {
		opts = append(opts, config.WithSharedConfigProfile(profile))
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), opts...)
	if err != nil {
		return nil, fmt.Errorf("loading AWS config: %w", err)
	}

	return &CloudFormationDetector{
		client: cloudformation.NewFromConfig(cfg),
	}, nil
}

// Check initiates drift detection for the given stack and prints results.
func (d *CloudFormationDetector) Check(ctx context.Context, stackName string) error {
	detectOut, err := d.client.DetectStackDrift(ctx, &cloudformation.DetectStackDriftInput{
		StackName: aws.String(stackName),
	})
	if err != nil {
		return fmt.Errorf("initiating drift detection: %w", err)
	}

	detectionID := detectOut.StackDriftDetectionId
	status, err := d.waitForDetection(ctx, stackName, detectionID)
	if err != nil {
		return err
	}

	timestamp := time.Now().Format(time.RFC3339)
	fmt.Printf("[%s] Stack %q drift status: %s\n", timestamp, stackName, status)
	return nil
}

func (d *CloudFormationDetector) waitForDetection(ctx context.Context, stackName string, detectionID *string) (types.StackDriftStatus, error) {
	for {
		out, err := d.client.DescribeStackDriftDetectionStatus(ctx, &cloudformation.DescribeStackDriftDetectionStatusInput{
			StackDriftDetectionId: detectionID,
		})
		if err != nil {
			return "", fmt.Errorf("describing drift detection status: %w", err)
		}

		switch out.DetectionStatus {
		case types.StackDriftDetectionStatusDetectionComplete:
			return out.StackDriftStatus, nil
		case types.StackDriftDetectionStatusDetectionFailed:
			return "", fmt.Errorf("drift detection failed for stack %q", stackName)
		}

		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(5 * time.Second):
		}
	}
}
