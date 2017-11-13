const {Seq} = require('immutable');
const Docker = require('dockerode');
const Slack = require('./slack');
const JSONStream = require('JSONStream');

const docker = new Docker();
const slack = new Slack({
  username: 'docker',
  iconEmoji: ':whale:',
});
const imageRegExp = new RegExp(process.env.image_regexp);

function formatAttachmentForEvent(e) {
  switch (e.status) {
    case 'start':
      return {
        color: 'good',
        text: 'Started container',
        fields: [
          {title: 'Image', value: e.from, short: true},
          {title: 'Container Name', value: e.Actor.Attributes.name, short: true},
          {title: 'Container ID', value: e.id},
        ]
      };

    case 'kill':
      return {
        color: 'warning',
        text: 'Container is stopped',
        fields: [
          {title: 'Image', value: e.from, short: true},
          {title: 'Container Name', value: e.Actor.Attributes.name, short: true},
        ]
      };

    case 'die':
      return {
        color: 'warning',
        text: 'Container is stopped',
        fields: [
          {title: 'Image', value: e.from, short: true},
          {title: 'Container Name', value: e.Actor.Attributes.name, short: true},
        ]
      };

    case 'destroy':
      return {
        color: 'warning',
        text: 'Container has been removed',
        fields: [
          {title: 'Image', value: e.from, short: true},
          {title: 'Container Name', value: e.Actor.Attributes.name, short: true},
        ]
      };
  }
}

async function sendEvent(event) {
  console.info(event);
  if (imageRegExp.test(event.from)) {
    const attachment = formatAttachmentForEvent(event);
    if (attachment) {
      await slack.sendAttachment(attachment);
    }
  }
}

async function sendEventStream() {
  const eventStream = await docker.getEvents();
  eventStream.pipe(JSONStream.parse())
    .on('data', event => sendEvent(event).catch(handleError))
    .on('error', handleError);
}

async function sendVersion() {
  const version = await docker.version();
  await slack.sendAttachment({
    text: 'Docker is running',
    color: 'good',
    fields: Seq(version).map((value, title) => ({title, value, short: true})).toArray(),
  });
}

async function main() {
  await sendVersion();
  await sendEventStream();
}

function handleError(e) {
  console.error(e);
  slack.sendError(e).catch(console.error);
}

main().catch(handleError);
