package formatter

import (
	"fmt"
	"github.com/int128/slack"
)

// Error returns a message for the error.
func Error(err error) *slack.Message {
	return &slack.Message{
		Username:  "docker",
		IconEmoji: ":whale:",
		Text:      fmt.Sprintf("Error while receiving events from Docker server: %s", err),
	}
}
