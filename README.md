# slack-docker [![CircleCI](https://circleci.com/gh/int128/slack-docker.svg?style=shield)](https://circleci.com/gh/int128/slack-docker)

A Slack integration to notify [Docker events](https://docs.docker.com/engine/reference/commandline/events/) written in Go.

<img width="596" alt="slack-docker-screenshot" src="https://user-images.githubusercontent.com/321266/47410763-c7682d80-d7a1-11e8-8f05-c80786152604.png">


## Getting Started

Setup [an Incoming WebHook](https://my.slack.com/services/new/incoming-webhook) on your Slack workspace and get the WebHook URL.

Run slack-docker as follows:

```sh
# Standalone
./slack-docker --webhook=https://hooks.slack.com/services/...

# Docker
docker run -d -e webhook=https://hooks.slack.com/services/... -v /var/run/docker.sock:/var/run/docker.sock int128/slack-docker

# Docker Compose
curl -O https://raw.githubusercontent.com/int128/slack-docker/master/docker-compose.yml
docker-compose up -d
```


## Configuration

It supports the following options and environment variables:

```
Application Options:
      --webhook=      Slack Incoming WebHook URL [$webhook]
      --image-regexp= Filter events by image name (default to all) [$image_regexp]

Help Options:
  -h, --help          Show this help message
```


### Filter events by image name

```sh
webhook=https://hooks.slack.com/services/... image_regexp='^alpine$' ./slack-docker
```


## Contribution

This is an open source software licensed under Apache-2.0.
Feel free to open issues or pull requests.
