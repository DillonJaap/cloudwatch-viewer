package timestamp

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/charmbracelet/bubbles/list"
)

var (
	_ list.Item         = Item{}          // Item implements list.Item
	_ list.ItemDelegate = &ItemDelegate{} // ItemDelegate implements list.ItemDelegate
)

// List.Item that contains cloudwatch events as its content
type Item struct {
	TimeStamp string
	Message   string
}

func (i Item) Title() string       { return i.TimeStamp }
func (i Item) Description() string { return "" }
func (i Item) FilterValue() string { return i.Message }

// TODO combine these two functions
func (i Item) getTruncatedTimeStamp(maxLength int) string {
	if maxLength < 10 {
		maxLength = 10
	}

	// TODO add error handling
	timeInt, _ := strconv.ParseInt(i.TimeStamp, 10, 64)
	time := time.Unix(timeInt, 0).String()
	if len(time) > maxLength {
		return time[0:maxLength-3] + "..."
	}
	return time
}

func (i Item) getTruncatedDescription(maxLength int) string {
	if maxLength < 10 {
		maxLength = 10
	}
	msg := strings.ReplaceAll(i.Message, "\t", " ")
	msg = strings.ReplaceAll(msg, "\n", " ")
	if len(msg) > maxLength {
		return msg[0:maxLength-3] + "..."
	}
	return msg
}

func logEventsToItemList(logEvents []types.OutputLogEvent) []list.Item {
	var items []list.Item
	for k := range logEvents {
		msg := aws.ToString(logEvents[k].Message)
		timeStamp := logEvents[k].Timestamp

		items = append(
			items,
			Item{
				Message:   msg,
				TimeStamp: fmt.Sprintf("%v", *timeStamp),
			},
		)
	}

	return items
}
