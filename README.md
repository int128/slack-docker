# slack-docker

A Slack integration to notify Docker events.

![slack-docker-0.1.1](https://cloud.githubusercontent.com/assets/321266/4773935/0141e0e8-5ba8-11e4-8b35-601e898c58be.png)

## How to Run

Set up [an incoming WebHook integration](https://my.slack.com/services/new/incoming-webhook) and get the Webhook URL.

Run a container as follows:

```sh
docker run -d -e webhook=URL -v /var/run/docker.sock:/var/run/docker.sock int128/slack-docker
```

Replace `URL` with your Webhook URL like `https://hooks.slack.com/services/...`


## Contribution

Please let me know an issue or pull request.
