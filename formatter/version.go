package formatter

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/int128/slack"
)

// Version returns a message for the Docker version.
func Version(v types.Version) *slack.Message {
	return &slack.Message{
		Username:  username,
		IconEmoji: iconEmoji,
		Text:      fmt.Sprintf("Docker version %s (%s)", v.Version, v.KernelVersion),
	}
}
