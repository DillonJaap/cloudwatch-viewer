package logstream

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/charmbracelet/bubbles/list"

	cw "clviewer/internal/cloudwatch"
)

const maxDescriptionLength = 90

// Item type
type Item string

func (i Item) FilterValue() string { return string(i) }

func GetLogStreamsAsItemList(pattern string) []list.Item {
	logStreams := cw.GetLogStreams(context.Background(), cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName: aws.String(pattern),
		Limit:        aws.Int32(30),
		OrderBy:      types.OrderByLastEventTime,
		Descending:   aws.Bool(true),
	})

	var streams []list.Item

	for k := range logStreams {
		name := aws.ToString(logStreams[k].LogStreamName)
		streams = append(streams, Item(name))
	}

	return streams
}

func (i Item) getTruncatedDescription(maxLength int) string {
	if maxLength < 10 {
		maxLength = 10
	}
	if len(i) > maxLength {
		return string(i[0:maxLength-3] + "...")
	}
	return string(i)
}
