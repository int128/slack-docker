package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/int128/slack-docker/formatter"
	"github.com/int128/slack-docker/slack"
	"github.com/jessevdk/go-flags"
	"io"
	"log"
)

type options struct {
	Webhook string `long:"webhook" env:"webhook" description:"Slack Webhook URL" required:"1"`
}

func (o *options) Run(ctx context.Context) error {
	docker, err := client.NewEnvClient()
	if err != nil {
		return fmt.Errorf("Could not create a Docker client: %s", err)
	}
	if err := o.showVersion(ctx, docker); err != nil {
		return err
	}
	if err := o.showEvents(ctx, docker); err != nil {
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
		return fmt.Errorf("Could not send a message to Slack: %s", err)
	}
	return nil
}

func (o *options) showEvents(ctx context.Context, docker *client.Client) error {
	msgCh, errCh := docker.Events(ctx, types.EventsOptions{})
	for {
		select {
		case msg := <-msgCh:
			log.Printf("Event %+v", msg)
			m := formatter.Event(msg)
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

func main() {
	var o options
	args, err := flags.NewParser(&o, flags.HelpFlag).Parse()
	if err != nil {
		log.Fatal(err)
	}
	if len(args) > 0 {
		log.Fatalf("Too many arguments")
	}
	o.Run(context.Background())
}
