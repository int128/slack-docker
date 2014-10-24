# slack-docker

A Slack integration to notify Docker events.

![slack-docker-0.1.0](https://cloud.githubusercontent.com/assets/321266/4516390/0a61ecb8-4bf4-11e4-9928-5992c2d0a395.png)

## Run

Set up [an incoming WebHook integration](https://my.slack.com/services/new/incoming-webhook) and get the token.

Run a container

```sh
docker run -d -v /var/run/docker.sock:/var/run/docker.sock -e webhook=URL -e channel=infra int128/slack-docker
```

with following environment variables.

* `webhook` is the Webhook URL like `https://hooks.slack.com/services/...`
* `channel` is the channel without hash prefix. Default is `general`.


## Contribution

Please let me know an issue or pull request.
