package main

import (
	"context"
	"os"

	"github.com/int128/slack-docker/cmd"
)

func main() {
	os.Exit(cmd.Run(context.Background(), os.Args))
}
