package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"regexp"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/int128/slack"
	"github.com/int128/slack-docker/formatter"
	"github.com/jessevdk/go-flags"
)

// Run parses the command line arguments and run the corresponding task.
// You need to pass os.Args to the 2nd argument.
func Run(ctx context.Context, osArgs []string) int {
	var o options
	args, err := flags.NewParser(&o, flags.HelpFlag).ParseArgs(osArgs[1:])
	if err != nil {
		log.Printf("invalid argument: %s", err)
		return 1
	}
	if len(args) > 0 {
		log.Printf("too many arguments")
		return 1
	}
	if err := o.Run(ctx); err != nil {
		log.Printf("error: %s", err)
		return 1
	}
	return 0
}

type options struct {
	Webhook           string `long:"webhook" env:"webhook" description:"Slack Incoming WebHook URL" required:"1"`
	DockerImageRegexp string `long:"image-regexp" env:"image_regexp" description:"Filter events by image name (default to all)"`
}

func (o *options) Run(ctx context.Context) error {
	var eventFilter formatter.EventFilter
	if o.DockerImageRegexp != "" {
		r, err := regexp.Compile(o.DockerImageRegexp)
		if err != nil {
			return fmt.Errorf("Invalid image-regexp: %s", err)
		}
		eventFilter.ImageRegexp = r
	}
	docker, err := client.NewEnvClient()
	if err != nil {
		return fmt.Errorf("Could not create a Docker client: %s", err)
	}
	if err := o.showVersion(ctx, docker); err != nil {
		return err
	}
	if err := o.showEvents(ctx, docker, eventFilter); err != nil {
		return err
	}
	return nil
}

func (o *options) showVersion(ctx context.Context, docker *client.Client) error {
	v, err := docker.ServerVersion(ctx)
	if err != nil {
		return fmt.Errorf("Could not get version from the Docker server: %s", err)
	}
	log.Printf("Connected to Docker server: %+v", v)
	if err := slack.Send(o.Webhook, formatter.Version(v)); err != nil {
		return fmt.Errorf("could not send a message to Slack: %s", err)
	}
	return nil
}

func (o *options) showEvents(ctx context.Context, docker *client.Client, filter formatter.EventFilter) error {
	msgCh, errCh := docker.Events(ctx, types.EventsOptions{})
	for {
		select {
		case msg := <-msgCh:
			log.Printf("Event %+v", msg)
			m := formatter.Event(msg, filter)
			if m != nil {
				if err := slack.Send(o.Webhook, m); err != nil {
					log.Printf("Error while sending a message to Slack: %s", err)
				}
			}
		case err := <-errCh:
			if err == io.EOF {
				break
			}
			log.Printf("Error while receiving events from Docker server: %s", err)
			if err := slack.Send(o.Webhook, formatter.Error(err)); err != nil {
				log.Printf("Error while sending a message to Slack: %s", err)
			}
		case <-ctx.Done():
			break
		}
	}
}
