package route53

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	awsroute53 "github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

const (
	defaultTTL    int64 = 300
	recordSetType       = "A"
)

func NewClient(ctx context.Context) (*awsroute53.Client, error) {
	cfg, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("load AWS config: %w", err)
	}

	log.Printf("route53: AWS config loaded")
	return awsroute53.NewFromConfig(cfg), nil
}

func UpsertZone(ctx context.Context, client *awsroute53.Client, hostedZoneID string, recordSet string, recordSetValue string) (string, error) {
	log.Printf("route53: building UPSERT change set")

	input := &awsroute53.ChangeResourceRecordSetsInput{
		HostedZoneId: aws.String(hostedZoneID),
		ChangeBatch: &types.ChangeBatch{
			Changes: []types.Change{
				{
					Action: types.ChangeActionUpsert,
					ResourceRecordSet: &types.ResourceRecordSet{
						Name: aws.String(recordSet),
						Type: types.RRType(recordSetType),
						TTL:  aws.Int64(defaultTTL),
						ResourceRecords: []types.ResourceRecord{
							{Value: aws.String(recordSetValue)},
						},
					},
				},
			},
		},
	}

	output, err := client.ChangeResourceRecordSets(ctx, input)
	if err != nil {
		return "", fmt.Errorf("change resource record sets: %w", err)
	}

	if output.ChangeInfo == nil {
		return "route53 change submitted; no change info returned", nil
	}

	submittedAt := ""
	if output.ChangeInfo.SubmittedAt != nil {
		submittedAt = output.ChangeInfo.SubmittedAt.UTC().Format(time.RFC3339)
	}

	return fmt.Sprintf(
		"change_id=%s status=%s submitted_at=%s",
		aws.ToString(output.ChangeInfo.Id),
		output.ChangeInfo.Status,
		submittedAt,
	), nil
}
