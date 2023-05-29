package stream

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

type Paginator struct {
	logGroup         string
	streamsPaginator *cloudwatchlogs.DescribeLogStreamsPaginator
}

func New(ctx context.Context, logGroupName string) Paginator {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}

	cw := cloudwatchlogs.NewFromConfig(cfg)

	// get log events paginator
	streamsPaginator := cloudwatchlogs.NewDescribeLogStreamsPaginator(
		cw,
		&cloudwatchlogs.DescribeLogStreamsInput{
			LogGroupName: aws.String(logGroupName),
			Limit:        aws.Int32(50),
			Descending:   aws.Bool(true),
			OrderBy:      types.OrderByLastEventTime,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	return Paginator{
		logGroup:         logGroupName,
		streamsPaginator: streamsPaginator,
	}
}

// Get next page of events, return nil if no pages remain
func (ep Paginator) NextPage(ctx context.Context) []types.LogStream {
	if !ep.streamsPaginator.HasMorePages() {
		return nil
	}
	streamsOutput, err := ep.streamsPaginator.NextPage(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return streamsOutput.LogStreams
}
