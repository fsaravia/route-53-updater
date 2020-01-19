package route53

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

const (
	defaultTTL    = 300
	recordSetType = "A"
	upsertAction  = "UPSERT"
)

func CreateSession() *route53.Route53 {
	mySession := session.Must(session.NewSession())

	return route53.New(mySession)
}

func UpsertZone(session *route53.Route53, hostedZoneId string, recordSet string, recordSetValue string) (string, error) {
	change := &route53.Change{
		Action: aws.String(upsertAction),
		ResourceRecordSet: &route53.ResourceRecordSet{
			Name: aws.String(recordSet),
			Type: aws.String(recordSetType),
			TTL:  aws.Int64(defaultTTL),
			ResourceRecords: []*route53.ResourceRecord{
				{
					Value: aws.String(recordSetValue),
				},
			},
		},
	}

	input := &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: aws.String(hostedZoneId),
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{change},
		},
	}

	output, err := session.ChangeResourceRecordSets(input)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%#v", output), nil
}
