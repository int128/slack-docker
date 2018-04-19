# slack-docker

A Slack integration to notify Docker events.

<img width="623" alt="slack-docker-screenshot" src="https://user-images.githubusercontent.com/321266/32720381-bc847924-c8a6-11e7-8348-fa7e03e82939.png">

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

### Include docker hostname in the messages

If you run standalone docker hosts you'll probably find usefull include the `HOSTNAME` in the 
messages posted to slack. Run the container as follows:

```sh
# Docker
docker run --net=host --env include_hostname=1 -d -e webhook=URL -v /var/run/docker.sock:/var/run/docker.sock int128/slack-docker
```

or set the variable explicitly

```sh
# Docker
docker run --env HOSTNAME=${HOSTNAME} --env include_hostname=1 -d -e webhook=URL -v /var/run/docker.sock:/var/run/docker.sock int128/slack-docker
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
