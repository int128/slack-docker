"use strict";

const Promise     = require('bluebird');
const JSONStream  = require('JSONStream');
const EventStream = require('event-stream');
const Dockerode   = require('dockerode');
const SlackNotify = require('slack-notify');

const Docker = new Dockerode();
Promise.promisifyAll(Docker);

const Slack = SlackNotify(process.env.webhook);
Promise.promisifyAll(Slack);

const EventFilter = {
  _imageRegExp: new RegExp(process.env.image_regexp),
  satisfy(event) {
    return this._imageRegExp.test(event.from);
  }
};

class EventProcessor {
  constructor() {
    this._containers = {};
  }
  _format(container, text, fields) {
    return {
      username: `docker${container.Name}`,
      icon_emoji: ':whale:',
      channel: '',
      text: text,
      fields: fields
    };
  }
  handle(event) {
    const method = this[`handle_${event.status}`];
    if (method) {
      return method.call(this, event);
    }
  }
  handle_start(event) {
    const containerObject = Docker.getContainer(event.id);
    Promise.promisifyAll(containerObject);
    containerObject.inspectAsync().then((container) => {
      this._containers[event.id] = container;
      return Slack.sendAsync(this._format(container, `Started ${container.Config.Hostname}`, {
        'Image': event.from,
        'IP Address': container.NetworkSettings.IPAddress,
        'Path': container.Path,
        'Arguments': container.Args,
        'Started at': container.State.StartedAt
      }));
    });
  }
  handle_kill(event) {
    const container = this._containers[event.id];
    if (container) {
      return Slack.sendAsync(this._format(container, `Stopped ${container.Config.Hostname}`));
    }
  }
  handle_die(event) {
    return this.handle_kill(event);
  }
  handle_destroy(event) {
    const container = this._containers[event.id];
    if (container) {
      delete this._containers[event.id];
      return Slack.sendAsync(this._format(container, `Removed ${container.Config.Hostname}`));
    }
  }
}

Docker.versionAsync()
.then((version) => console.info(version))
.then(() => Docker.getEventsAsync())
.then((stream) => {
  const eventProcessor = new EventProcessor();
  stream.pipe(JSONStream.parse())
    .pipe(EventStream.map((event) => console.info(event)))
    .pipe(EventStream.map((event) => eventProcessor.handle(event)))
})
.catch((e) => console.error(e));
