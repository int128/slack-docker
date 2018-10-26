package formatter

import (
	"fmt"
	"github.com/docker/docker/api/types/events"
	"github.com/int128/slack"
	"regexp"
)

// EventFilter represents a filter for events.
type EventFilter struct {
	ImageRegexp *regexp.Regexp
}

// Match returns true if the event satisfies the filter.
func (filter *EventFilter) Match(e events.Message) bool {
	if filter.ImageRegexp != nil {
		switch e.Type {
		case "container":
			return filter.ImageRegexp.MatchString(e.Actor.Attributes["image"])
		}
	}
	return true
}

// Event returns a message for the event.
// It returns a message only for some events and may return nil.
func Event(e events.Message, filter EventFilter) *slack.Message {
	if !filter.Match(e) {
		return nil
	}
	switch e.Type {
	case "container":
		return containerEvent(e, filter)
	case "image":
		return imageEvent(e)
	case "volume":
		return volumeEvent(e)
	case "network":
		return networkEvent(e)
	}
	return nil
}

func containerEvent(e events.Message, filter EventFilter) *slack.Message {
	switch {
	case e.Action == "start":
		return &slack.Message{
			Username:  username,
			IconEmoji: iconEmoji,
			Attachments: []slack.Attachment{
				slack.Attachment{
					Title:     fmt.Sprintf(":package: Started %s", e.Actor.Attributes["name"]),
					Text:      fmt.Sprintf("Container `%s` has been started with the image `%s`.", e.Actor.Attributes["name"], e.Actor.Attributes["image"]),
					Footer:    e.Actor.ID,
					Color:     "good",
					Timestamp: e.Time,
					MrkdwnIn:  []string{"text"},
				},
			},
		}
	case e.Action == "kill":
		fallthrough
	case e.Action == "die":
		var color string
		switch e.Actor.Attributes["exitCode"] {
		case "0":
			color = "good"
		default:
			color = "danger"
		}
		return &slack.Message{
			Username:  username,
			IconEmoji: iconEmoji,
			Attachments: []slack.Attachment{
				slack.Attachment{
					Title:     fmt.Sprintf(":package: Stopped %s", e.Actor.Attributes["name"]),
					Text:      fmt.Sprintf("Container `%s` has been terminated with exit code `%s`.", e.Actor.Attributes["name"], e.Actor.Attributes["exitCode"]),
					Footer:    e.Actor.ID,
					Color:     color,
					Timestamp: e.Time,
					MrkdwnIn:  []string{"text"},
				},
			},
		}
	case e.Action == "destroy":
		return &slack.Message{
			Username:  username,
			IconEmoji: iconEmoji,
			Attachments: []slack.Attachment{
				slack.Attachment{
					Title:     fmt.Sprintf(":package: Removed %s", e.Actor.Attributes["name"]),
					Text:      fmt.Sprintf("Container `%s` has been removed.", e.Actor.Attributes["name"]),
					Footer:    e.Actor.ID,
					Color:     "warning",
					Timestamp: e.Time,
					MrkdwnIn:  []string{"text"},
				},
			},
		}
	}
	return nil
}

func imageEvent(e events.Message) *slack.Message {
	switch {
	case e.Action == "tag":
		return &slack.Message{
			Username:  username,
			IconEmoji: iconEmoji,
			Attachments: []slack.Attachment{
				slack.Attachment{
					Title:     fmt.Sprintf(":label: Tagged %s", e.Actor.Attributes["name"]),
					Text:      fmt.Sprintf("Image `%s` has been tagged.", e.Actor.Attributes["name"]),
					Footer:    e.Actor.ID,
					Color:     "good",
					Timestamp: e.Time,
					MrkdwnIn:  []string{"text"},
				},
			},
		}
	case e.Action == "delete":
		return &slack.Message{
			Username:  username,
			IconEmoji: iconEmoji,
			Attachments: []slack.Attachment{
				slack.Attachment{
					Title:     fmt.Sprintf(":label: Deleted %s", e.Actor.Attributes["name"]),
					Text:      fmt.Sprintf("Image `%s` has been deleted.", e.Actor.Attributes["name"]),
					Footer:    e.Actor.ID,
					Color:     "good",
					Timestamp: e.Time,
					MrkdwnIn:  []string{"text"},
				},
			},
		}
	}
	return nil
}

func volumeEvent(e events.Message) *slack.Message {
	switch {
	case e.Action == "create":
		return &slack.Message{
			Username:  username,
			IconEmoji: iconEmoji,
			Attachments: []slack.Attachment{
				slack.Attachment{
					Title:     fmt.Sprintf(":oil_drum: Started %s", e.Actor.Attributes["name"]),
					Text:      fmt.Sprintf("Volume `%s` has been created.", e.Actor.Attributes["name"]),
					Footer:    e.Actor.ID,
					Color:     "good",
					Timestamp: e.Time,
					MrkdwnIn:  []string{"text"},
				},
			},
		}
	case e.Action == "destroy":
		return &slack.Message{
			Username:  username,
			IconEmoji: iconEmoji,
			Attachments: []slack.Attachment{
				slack.Attachment{
					Title:     fmt.Sprintf(":oil_drum: Removed %s", e.Actor.Attributes["name"]),
					Text:      fmt.Sprintf("Volume `%s` has been removed.", e.Actor.Attributes["name"]),
					Footer:    e.Actor.ID,
					Color:     "warning",
					Timestamp: e.Time,
					MrkdwnIn:  []string{"text"},
				},
			},
		}
	}
	return nil
}

func networkEvent(e events.Message) *slack.Message {
	switch {
	case e.Action == "create":
		return &slack.Message{
			Username:  username,
			IconEmoji: iconEmoji,
			Attachments: []slack.Attachment{
				slack.Attachment{
					Title:     fmt.Sprintf(":link: Started %s", e.Actor.Attributes["name"]),
					Text:      fmt.Sprintf("Network `%s` has been created.", e.Actor.Attributes["name"]),
					Footer:    e.Actor.ID,
					Color:     "good",
					Timestamp: e.Time,
					MrkdwnIn:  []string{"text"},
				},
			},
		}
	case e.Action == "destroy":
		return &slack.Message{
			Username:  username,
			IconEmoji: iconEmoji,
			Attachments: []slack.Attachment{
				slack.Attachment{
					Title:     fmt.Sprintf(":link: Removed %s", e.Actor.Attributes["name"]),
					Text:      fmt.Sprintf("Network `%s` has been removed.", e.Actor.Attributes["name"]),
					Footer:    e.Actor.ID,
					Color:     "warning",
					Timestamp: e.Time,
					MrkdwnIn:  []string{"text"},
				},
			},
		}
	}
	return nil
}
