package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"

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
		LogGroupNamePrefix: aws.String("/aws/lambda"),
	})
	if err != nil {
		log.Fatal(err)
	}

	// get log streams
	logStreamsOutput, err := cwClient.DescribeLogStreams(ctx, &cloudwatchlogs.DescribeLogStreamsInput{
		Limit:        aws.Int32(1),
		LogGroupName: logGroupsOutput.LogGroups[0].LogGroupName,
	})
	if err != nil {
		log.Fatal(err)
	}

	// get log events
	eventsOutput, err := cwClient.GetLogEvents(ctx, &cloudwatchlogs.GetLogEventsInput{
		Limit:         aws.Int32(5),
		LogGroupName:  logGroupsOutput.LogGroups[0].LogGroupName,
		LogStreamName: logStreamsOutput.LogStreams[0].LogStreamName,
	})
	if err != nil {
		log.Fatal(err)
	}

	regx, _ := regexp.Compile(`.*(?P<json>{.*}).*`)

	var formattedEvents []list.Item

	for k := range eventsOutput.Events {
		submatches := regx.FindStringSubmatch(
			aws.ToString(eventsOutput.Events[k].Message),
		)
		timeStamp := eventsOutput.Events[k].Timestamp

		if len(submatches) == 0 {
			continue
		}
		jsonMessage := submatches[1]

		formattedEvents = append(
			formattedEvents,
			Event{Message: formatJson(jsonMessage), TimeStamp: fmt.Sprintf("%v", *timeStamp)},
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
