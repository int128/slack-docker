# slack-docker

A Slack integration to notify Docker events.


## Run

Set up [an incoming WebHook integration](https://my.slack.com/services/new/incoming-webhook) and get the token.

Run a container

```sh
docker run -d -v /var/run/docker.sock:/var/run/docker.sock -e domain=DOMAIN -e token=TOKEN -e channel=infra int128/slack-docker
```

with following environment variables.

* `domain` is the first part of your .slack.com.
* `token` is the token provided on the integration page.
* `channel` is the channel without hash prefix. Default is `general`.


## Contribution

Please let me know an issue or pull request.
