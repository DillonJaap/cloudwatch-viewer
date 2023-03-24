package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/TylerBrock/colorjson"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/charmbracelet/bubbles/list"
)

func GetEvents(ctx context.Context) []list.Item {
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

	for _, v := range logGroupsOutput.LogGroups {
		fmt.Printf("%v\n", aws.ToString(v.LogGroupName))
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

	for _, v := range logStreamsOutput.LogStreams {
		fmt.Printf("%v\n", aws.ToString(v.LogStreamName))
	}

	// get log events
	logStreamName := logStreamsOutput.LogStreams[0].LogStreamName
	logEventsOutput, err := cwClient.GetLogEvents(ctx, &cloudwatchlogs.GetLogEventsInput{
		LogStreamName: logStreamName,
		LogGroupName:  logGroupName,
		Limit:         aws.Int32(20),
		StartFromHead: aws.Bool(true),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v\n", len(logEventsOutput.Events))

	regx, _ := regexp.Compile(`(.*)(?P<json>{.*})(.*)`)
	var formattedEvents []list.Item

	for k := range logEventsOutput.Events {
		msg := aws.ToString(logEventsOutput.Events[k].Message)

		submatches := regx.FindStringSubmatch(msg)

		timeStamp := logEventsOutput.Events[k].Timestamp

		if len(submatches) > 1 {
			msg = submatches[1] + "\n" + formatJson(submatches[2]) + "\n" + submatches[3]
		}

		formattedEvents = append(
			formattedEvents,
			Event{Message: strings.ReplaceAll(msg, "\t", " "), TimeStamp: fmt.Sprintf("%v", *timeStamp)},
		)
	}

	return formattedEvents
}

func formatJson(in string) string {
	var obj map[string]interface{}
	json.Unmarshal([]byte(in), &obj)

	f := colorjson.NewFormatter()
	f.Indent = 2

	s, _ := f.Marshal(obj)
	return string(s)
}
