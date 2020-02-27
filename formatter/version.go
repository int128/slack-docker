package formatter

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/int128/slack"
	"strings"
	"time"
)

// Version returns a message for the Docker version.
func Version(v types.Version, filter EventFilter) *slack.Message {
	var textArray []string
	if filter.TypeRegexp != nil {
		textArray = append(textArray, fmt.Sprintf("TypeRegexp: %s", filter.TypeRegexp.String()))
	}
	if filter.ImageRegexp != nil {
		textArray = append(textArray, fmt.Sprintf("ImageRegexp: %s", filter.ImageRegexp.String()))
	}
	if filter.ContainerNameRegexp != nil {
		textArray = append(textArray, fmt.Sprintf("ContainerNameRegexp: %s", filter.ContainerNameRegexp.String()))
	}
	text := strings.Join(textArray, ", ")
	message := &slack.Message{
		Username:  username,
		IconEmoji: iconEmoji,
		Attachments: []slack.Attachment{
			slack.Attachment{
				Title:     "Started Slack-Docker",
				Text:      text,
				Footer:    fmt.Sprintf("Docker version %s (%s)", v.Version, v.KernelVersion),
				Color:     "#00a29a",
				Timestamp: time.Now().Unix(),
				MrkdwnIn:  []string{"text"},
			},
		},
	}

	return message
}
