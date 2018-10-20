package formatter

import (
	"fmt"
	"github.com/docker/docker/api/types/events"
	"github.com/int128/slack-docker/slack"
)

// Event returns a message for the event.
// It returns a message only for some events and may return nil.
func Event(e events.Message) *slack.Message {
	switch {
	case e.Type == "container" && e.Action == "start":
		return &slack.Message{
			Username:  username,
			IconEmoji: iconEmoji,
			Attachments: slack.Attachments{
				slack.Attachment{
					Text:     fmt.Sprintf("Started container `%s`", e.Actor.Attributes["name"]),
					Color:    "good",
					MrkdwnIn: []string{"text"},
					Fields: slack.AttachmentFields{
						slack.AttachmentField{
							Short: true,
							Title: "Image",
							Value: e.Actor.Attributes["image"],
						},
						slack.AttachmentField{
							Title: "Container ID",
							Value: e.Actor.ID,
						},
					},
				},
			},
		}
	case e.Type == "container" && e.Action == "kill":
		fallthrough
	case e.Type == "container" && e.Action == "die":
		var color string
		if e.Actor.Attributes["exitCode"] == "0" {
			color = "good"
		} else {
			color = "danger"
		}
		return &slack.Message{
			Username:  username,
			IconEmoji: iconEmoji,
			Attachments: slack.Attachments{
				slack.Attachment{
					Text:     fmt.Sprintf("Stopped container `%s`", e.Actor.Attributes["name"]),
					Color:    color,
					MrkdwnIn: []string{"text"},
					Fields: slack.AttachmentFields{
						slack.AttachmentField{
							Short: true,
							Title: "Image",
							Value: e.Actor.Attributes["image"],
						},
						slack.AttachmentField{
							Short: true,
							Title: "Exit Code",
							Value: e.Actor.Attributes["exitCode"],
						},
						slack.AttachmentField{
							Title: "Container ID",
							Value: e.Actor.ID,
						},
					},
				},
			},
		}
	case e.Type == "container" && e.Action == "destroy":
		return &slack.Message{
			Username:  username,
			IconEmoji: iconEmoji,
			Attachments: slack.Attachments{
				slack.Attachment{
					Text:     fmt.Sprintf("Destroyed container `%s`", e.Actor.Attributes["name"]),
					Color:    "warning",
					MrkdwnIn: []string{"text"},
					Fields: slack.AttachmentFields{
						slack.AttachmentField{
							Short: true,
							Title: "Image",
							Value: e.Actor.Attributes["image"],
						},
						slack.AttachmentField{
							Title: "Container ID",
							Value: e.Actor.ID,
						},
					},
				},
			},
		}
	case e.Type == "network" && e.Action == "create":
		return &slack.Message{
			Username:  username,
			IconEmoji: iconEmoji,
			Attachments: slack.Attachments{
				slack.Attachment{
					Text:     fmt.Sprintf("Created network `%s`", e.Actor.Attributes["name"]),
					Color:    "good",
					MrkdwnIn: []string{"text"},
				},
			},
		}
	case e.Type == "network" && e.Action == "destroy":
		return &slack.Message{
			Username:  username,
			IconEmoji: iconEmoji,
			Attachments: slack.Attachments{
				slack.Attachment{
					Text:     fmt.Sprintf("Destroyed network `%s`", e.Actor.Attributes["name"]),
					Color:    "warning",
					MrkdwnIn: []string{"text"},
				},
			},
		}
	}
	return nil
}
