package logevent

import (
	"clviewer/internal/events"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/TylerBrock/colorjson"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

const maxDescriptionLength = 50

type Event struct {
	TimeStamp string
	Message   string
}

func (e Event) Title() string       { return e.TimeStamp }
func (e Event) Description() string { return e.getTruncatedDescription() }
func (e Event) FilterValue() string { return e.Message }

func (e Event) getTruncatedDescription() string {
	// TODO make const
	if len(e.Message) > maxDescriptionLength {
		return e.Message[0:maxDescriptionLength-3] + "..."
	}
	return e.Message
}

type eventDelegate struct{}

func (d *eventDelegate) Height() int { return 5 }

func (d *eventDelegate) Spacing() int { return 0 }

func (d *eventDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	m.VisibleItems()
	return nil
}

func (d *eventDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Event)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.Message)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + s[0])
		}
	}

	fmt.Fprint(w, fn(str))
}

func GetLogEventsAsItemList() []list.Item {
	logEvents := events.GetEvents(context.Background())

	var events []list.Item

	for k := range logEvents {
		msg := aws.ToString(logEvents[k].Message)
		timeStamp := logEvents[k].Timestamp

		events = append(
			events,
			Event{
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
		event, ok := item.(Event)
		if !ok {
			formattedList = append(formattedList, item)
			continue
		}

		formattedList = append(
			formattedList,
			Event{
				Message:   formatMessage(event.Message, formatAsJson),
				TimeStamp: event.TimeStamp,
			},
		)
	}
	return formattedList
}

func formatMessage(in string, formatAsJson bool) string {
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
