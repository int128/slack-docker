# slack-docker

Slack integration that notifies Docker events on your host.


## How to Run

Set up [an incoming WebHook](https://my.slack.com/services/new/incoming-webhook) and get the token.

Run a container,

```sh
docker run -d -v /var/run/docker.sock:/var/run/docker.sock -e domain=DOMAIN -e token=TOKEN -e channel=#infra int128/slack-docker
```

with following environment variables.

* `domain` is the first part of your .slack.com. (Mandatory)
* `token` is the token provided on the integration setup page. (Mandatory)
* `channel` is the channel. (Default is `#general`)

