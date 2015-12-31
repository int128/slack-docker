"use strict";

const Promise     = require('bluebird');
const JSONStream  = require('JSONStream');
const EventStream = require('event-stream');
const Dockerode   = require('dockerode');
const Slack       = require('slack-notify');

class EventFilter {
  constructor(image) {
    this._imageRegExp = new RegExp(image);
  }
  filter() {
    return EventStream.map((event, callback) => {
      if (this._imageRegExp.test(event.from)) {
        callback(null, event);
      } else {
        callback();
      }
    });
  }
}

class EventInspector {
  constructor(docker) {
    this._docker = docker;
    this._containers = {};
  }
  map() {
    return EventStream.map((event, callback) => {
      if (this._containers[event.id]) {
        event.container = this._containers[event.id];
        callback(null, event);
      } else {
        Promise.promisifyAll(this._docker.getContainer(event.id))
        .inspectAsync()
        .then((container) => {
          event.container = this._containers[event.id] = container;
          callback(null, event);
        });
      }
      if (event.status == 'destroy') {
        delete this._containers[event.id];
      }
    });
  }
}

class EventNotifier {
  constructor(slack) {
    this._slack = slack;
  }
  _send(event, text, fields) {
    return this._slack.sendAsync({
      username: `docker${event.container.Name}`,
      icon_emoji: ':whale:',
      channel: '',
      text: text,
      fields: fields
    });
  }
  _map_start(e) {
		return this._send(e, `Started ${e.container.Config.Hostname}`, {
      'Image': e.from,
      'IP Address': e.container.NetworkSettings.IPAddress,
      'Path': e.container.Path,
      'Arguments': e.container.Args,
      'Started at': e.container.State.StartedAt
    });
  }
  _map_kill(e) {
		return this._send(e, `Stopped ${e.container.Config.Hostname}`);
	}
  _map_die(e) {
		return this._send(e, `Stopped ${e.container.Config.Hostname}`);
	}
  _map_destroy(e) {
		return this._send(e, `Removed ${e.container.Config.Hostname}`);
	}
  map() {
    return EventStream.map((event, callback) => {
      if (this[`_map_${event.status}`]) {
        this[`_map_${event.status}`](event).then((sent) => callback(null, sent));
      }
    });
  }
}

const docker = Promise.promisifyAll(new Dockerode());
const slack  = Promise.promisifyAll(Slack(process.env.webhook));

docker.versionAsync()
.then((version) => console.info(version))
.then(() => docker.getEventsAsync())
.then((stream) => stream
  .pipe(JSONStream.parse())
  .pipe(new EventFilter(process.env.image_regexp).filter())
  .pipe(new EventInspector(docker).map())
  .pipe(new EventNotifier(slack).map())
).catch((e) => console.error(e));
