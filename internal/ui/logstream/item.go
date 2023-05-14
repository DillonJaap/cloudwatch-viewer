package logstream

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/charmbracelet/bubbles/list"
)

const maxDescriptionLength = 90

// Item type
type Item struct {
	timestamp string
	name      string
}

func (i Item) FilterValue() string { return i.timestamp + " | " + i.name }

func GetLogStreamsAsItemList(streams []types.LogStream) []list.Item {
	var items []list.Item
	for k := range streams {
		name := aws.ToString(streams[k].LogStreamName)

		timeInt := aws.ToInt64(streams[k].FirstEventTimestamp)
		timestamp := time.Unix(0, timeInt*int64(time.Millisecond)).String()

		items = append(items, Item{
			timestamp: timestamp,
			name:      name,
		})
	}
	return items
}

func (i Item) getTruncatedDescription(maxLength int) string {
	if maxLength < 10 {
		maxLength = 10
	}
	if len(i.FilterValue()) > maxLength {
		return i.FilterValue()[0:maxLength-3] + "..."
	}
	return i.FilterValue()
}
