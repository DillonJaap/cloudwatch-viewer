package loggrouplist

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/charmbracelet/bubbles/list"

	cw "clviewer/internal/cloudwatch"
)

const maxDescriptionLength = 90

// Item type
type Item string

func (i Item) FilterValue() string { return string(i) }

func GetLogGroupsAsItemList(pattern string) []list.Item {
	logGroups := cw.GetLogGroups(context.Background(), cloudwatchlogs.DescribeLogGroupsInput{
		LogGroupNamePrefix: aws.String(pattern),
	})

	var groups []list.Item

	for k := range logGroups {
		name := aws.ToString(logGroups[k].LogGroupName)
		groups = append(groups, Item(name))
	}

	return groups
}

func (i Item) getTruncatedDescription(maxLength int) string {
	if len(i) > maxLength {
		return string(i[0:maxLength-3] + "...")
	}
	return string(i)
}
