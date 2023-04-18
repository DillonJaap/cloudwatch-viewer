package cloudwatch

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

func GetLogGroups(ctx context.Context, in cloudwatchlogs.DescribeLogGroupsInput) []types.LogGroup {
	// Load the Shared AWS Configuration(~/.aws/config)
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// get cloudwatch client
	cwClient := cloudwatchlogs.NewFromConfig(cfg)
	cwPaginator := cloudwatchlogs.NewDescribeLogGroupsPaginator(cwClient, &in)

	// get all the log groups via paginator
	var logGroups []types.LogGroup
	for cwPaginator.HasMorePages() {
		output, err := cwPaginator.NextPage(ctx)
		logGroups = append(logGroups, output.LogGroups...)
		if err != nil {
			log.Fatal(err)
		}
	}

	return logGroups
}

func GetLogStreams(ctx context.Context, in cloudwatchlogs.DescribeLogStreamsInput) []types.LogStream {
	// Load the Shared AWS Configuration(~/.aws/config)
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// get cloudwatch client
	cwClient := cloudwatchlogs.NewFromConfig(cfg)

	// get log groups
	logStreamsOutput, err := cwClient.DescribeLogStreams(ctx, &in)
	if err != nil {
		log.Fatal(err)
	}
	return logStreamsOutput.LogStreams
}