package logstream

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/charmbracelet/bubbles/list"
)

const maxDescriptionLength = 90

// Item type
type Item string

func (i Item) FilterValue() string { return string(i) }

func GetLogStreamsAsItemList(streams []types.LogStream) []list.Item {
	var items []list.Item
	for k := range streams {
		msg := aws.ToString(streams[k].LogStreamName)
		items = append(items, Item(msg))
	}
	return items
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
