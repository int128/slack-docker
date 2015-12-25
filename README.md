# slack-docker

A Slack integration to notify Docker events.

<img width="485" src="https://cloud.githubusercontent.com/assets/321266/11998879/e63bac0a-aaf4-11e5-935f-7cb840c3c501.png">

## How to Run

Set up [an incoming WebHook integration](https://my.slack.com/services/new/incoming-webhook) and get the Webhook URL.

Run a container as follows:

```sh
# Docker
docker run -d -e webhook=URL -v /var/run/docker.sock:/var/run/docker.sock int128/slack-docker

# Docker Compose
curl -O https://raw.githubusercontent.com/int128/slack-docker/master/docker-compose.yml
docker-compose up -d
```

### Filter events by image name

By default all events are sent to Slack, but events can be filtered by the environment variable `image_regexp` as follows:

```sh
# show events only from node
-e image_regexp='^node:'

# show events but exclude from node
-e image_regexp='^(?!node:)'
```


## Contribution

Please let me know an issue or pull request.
