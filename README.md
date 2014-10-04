# Slack integration for Docker events

## Setup Slack Incoming WebHook

Set up [an incoming WebHook](https://my.slack.com/services/new/incoming-webhook) on your team.

## Run

```sh
domain=DOMAIN token=TOKEN channel=#infra npm start
```

### Environment variables

* `domain` is the first part of your .slack.com. (Mandatory)
* `token` is the token provided on the integration setup page. (Mandatory)
* `channel` is the channel. (Default is `#general`)

