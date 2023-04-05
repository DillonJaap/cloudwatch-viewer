package cloudwatch

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

func GetEvents(ctx context.Context) []types.OutputLogEvent {
	// Load the Shared AWS Configuration(~/.aws/config)
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// get cloudwatch client
	cwClient := cloudwatchlogs.NewFromConfig(cfg)

	// get log groups
	logGroupsOutput, err := cwClient.DescribeLogGroups(ctx, &cloudwatchlogs.DescribeLogGroupsInput{
		LogGroupNamePattern: aws.String("/aws/lambda/dev-djaap-event-handlers-batch-processor"),
	})
	if err != nil {
		log.Fatal(err)
	}

	// get log streams
	logGroupName := logGroupsOutput.LogGroups[0].LogGroupName
	logStreamsOutput, err := cwClient.DescribeLogStreams(ctx, &cloudwatchlogs.DescribeLogStreamsInput{
		Limit:              aws.Int32(1),
		LogGroupIdentifier: logGroupName,
		Descending:         aws.Bool(true),
	})
	if err != nil {
		log.Fatal(err)
	}

	// get log events
	logStreamName := logStreamsOutput.LogStreams[0].LogStreamName
	logEventsOutput, err := cwClient.GetLogEvents(ctx, &cloudwatchlogs.GetLogEventsInput{
		LogStreamName: logStreamName,
		LogGroupName:  logGroupName,
		Limit:         aws.Int32(45),
		StartFromHead: aws.Bool(true),
	})
	if err != nil {
		log.Fatal(err)
	}

	return logEventsOutput.Events

}

func GetLogGroups(ctx context.Context, in cloudwatchlogs.DescribeLogGroupsInput) []types.LogGroup {
	// Load the Shared AWS Configuration(~/.aws/config)
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// get cloudwatch client
	cwClient := cloudwatchlogs.NewFromConfig(cfg)

	// get log groups
	logGroupsOutput, err := cwClient.DescribeLogGroups(ctx, &in)
	if err != nil {
		log.Fatal(err)
	}
	return logGroupsOutput.LogGroups
}
