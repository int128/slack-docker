const {promisify} = require('util');
const Slack = require('@slack/client');

class IncomingWebhook {
  constructor(defaults) {
    this.slack = new Slack.IncomingWebhook(process.env.webhook, defaults);
  }

  send(message) {
    return promisify(cb => this.slack.send(message, cb))();
  }

  async sendAttachment(attachment) {
    await this.send({attachments: [attachment]});
  }

  async sendError(e) {
    await this.sendAttachment({
      color: 'danger',
      fields: [{title: 'Error', value: `${e}`}],
      fallback: `${e}`,
    });
  }
}

module.exports = IncomingWebhook;
