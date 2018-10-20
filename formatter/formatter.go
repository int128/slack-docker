package formatter

import (
	"fmt"
	"log"
	"os"
)

const iconEmoji = ":whale:"

var username = "docker"

func init() {
	hostname, err := os.Hostname()
	if err != nil {
		log.Printf("Skip to set Slack username with hostname: %s", err)
	}
	username = fmt.Sprintf("docker@%s", hostname)
}
