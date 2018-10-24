# slack-docker

A Slack integration to notify [Docker events](https://docs.docker.com/engine/reference/commandline/events/) written in Go.

<img width="623" alt="slack-docker-screenshot" src="https://user-images.githubusercontent.com/321266/32720381-bc847924-c8a6-11e7-8348-fa7e03e82939.png">


## Getting Started

Setup [an Incoming WebHook](https://my.slack.com/services/new/incoming-webhook) on your Slack workspace and get the WebHook URL.

You can run slack-docker as follows:

```sh
# Standalone
webhook=https://hooks.slack.com/services/... ./slack-docker

# Docker
docker run -d -e webhook=https://hooks.slack.com/services/... -v /var/run/docker.sock:/var/run/docker.sock int128/slack-docker

# Docker Compose
curl -O https://raw.githubusercontent.com/int128/slack-docker/master/docker-compose.yml
docker-compose up -d
```


## Contribution

This is an open source software licensed under Apache-2.0.
Feel free to open issues or pull requests.
