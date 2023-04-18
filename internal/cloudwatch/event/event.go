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

func New(ctx context.Context, logGroupPattern, logStreamPrefix string) Paginator {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}

	cw := cloudwatchlogs.NewFromConfig(cfg)

	// get log groups
	logGroupsOutput, err := cw.DescribeLogGroups(ctx, &cloudwatchlogs.DescribeLogGroupsInput{
		LogGroupNamePattern: aws.String(logGroupPattern),
	})
	if err != nil {
		log.Fatal(err)
	}

	// get log streams
	logGroupName := logGroupsOutput.LogGroups[0].LogGroupName
	logStreamsOutput, err := cw.DescribeLogStreams(ctx, &cloudwatchlogs.DescribeLogStreamsInput{
		Limit:               aws.Int32(1),
		LogGroupIdentifier:  logGroupName,
		LogStreamNamePrefix: aws.String(logStreamPrefix),
		Descending:          aws.Bool(true),
	})
	if err != nil {
		log.Fatal(err)
	}

	// get log events paginator
	logStreamName := logStreamsOutput.LogStreams[0].LogStreamName
	eventsPaginator := cloudwatchlogs.NewGetLogEventsPaginator(
		cw,
		&cloudwatchlogs.GetLogEventsInput{
			LogStreamName: logStreamName,
			LogGroupName:  logGroupName,
			Limit:         aws.Int32(45),
			StartFromHead: aws.Bool(true),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	return Paginator{
		logGroup:        logGroupPattern,
		logStream:       logStreamPrefix,
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
