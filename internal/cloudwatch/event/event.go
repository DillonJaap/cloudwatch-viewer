package event

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

type Paginator struct {
	logGroup        string
	logStream       string
	eventsPaginator *cloudwatchlogs.GetLogEventsPaginator
}

func New(
	ctx context.Context,
	logGroupName, logStreamName, filterPattern string,
) Paginator {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}

	cw := cloudwatchlogs.NewFromConfig(cfg)

	// get log events paginator
	eventsPaginator := cloudwatchlogs.NewGetLogEventsPaginator(
		cw,
		&cloudwatchlogs.GetLogEventsInput{
			Limit:         aws.Int32(200),
			LogStreamName: aws.String(logStreamName),
			LogGroupName:  aws.String(logGroupName),
			StartFromHead: aws.Bool(true),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	return Paginator{
		logGroup:        logGroupName,
		logStream:       logStreamName,
		eventsPaginator: eventsPaginator,
	}
}

// Get next page of events, return nil if no pages remain
func (ep Paginator) NextPage(ctx context.Context) []types.OutputLogEvent {
	if !ep.eventsPaginator.HasMorePages() {
		return nil
	}
	eventsOutput, err := ep.eventsPaginator.NextPage(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return eventsOutput.Events
}
