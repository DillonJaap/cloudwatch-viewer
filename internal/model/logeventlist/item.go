package logeventlist

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/TylerBrock/colorjson"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/charmbracelet/bubbles/list"

	cw "clviewer/internal/cloudwatch"
)

var (
	_ list.Item         = Item{}          // Item implements list.Item
	_ list.ItemDelegate = &ItemDelegate{} // ItemDelegate implements list.ItemDelegate
)

const maxDescriptionLength = 50

// List.Item that contains cloudwatch events as its content
type Item struct {
	TimeStamp string
	Message   string
}

func (i Item) Title() string       { return i.TimeStamp }
func (i Item) Description() string { return i.getTruncatedDescription() }
func (i Item) FilterValue() string { return i.Message }

func (i Item) getTruncatedDescription() string {
	if len(i.Message) > maxDescriptionLength {
		return i.Message[0:maxDescriptionLength-3] + "..."
	}
	return i.Message
}

func GetLogEventsAsItemList() []list.Item {
	logEvents := cw.GetEvents(context.Background())

	var events []list.Item

	for k := range logEvents {
		msg := aws.ToString(logEvents[k].Message)
		timeStamp := logEvents[k].Timestamp

		events = append(
			events,
			Item{
				Message:   msg,
				TimeStamp: fmt.Sprintf("%v", *timeStamp),
			},
		)
	}

	return events
}

func formatList(itemList []list.Item, formatAsJson bool) []list.Item {
	var formattedList []list.Item
	for _, item := range itemList {
		event, ok := item.(Item)
		if !ok {
			formattedList = append(formattedList, item)
			continue
		}

		formattedList = append(
			formattedList,
			Item{
				Message:   FormatMessage(event.Message, formatAsJson),
				TimeStamp: event.TimeStamp,
			},
		)
	}
	return formattedList
}

func FormatMessage(in string, formatAsJson bool) string {
	in = strings.ReplaceAll(in, "\t", " ")

	regx, _ := regexp.Compile(`(.*)(?P<json>{.*})(.*)`)
	submatches := regx.FindStringSubmatch(in)

	if len(submatches) == 0 || !formatAsJson {
		return in
	}

	return submatches[1] +
		"\n" + formatJson(submatches[2]) +
		"\n" + submatches[3]
}

func formatJson(in string) string {
	var obj map[string]interface{}
	json.Unmarshal([]byte(in), &obj)

	f := colorjson.NewFormatter()
	f.Indent = 2

	s, _ := f.Marshal(obj)
	return string(s)
}
