const {Seq} = require('immutable');
const Docker = require('dockerode');
const Slack = require('./slack');
const JSONStream = require('JSONStream');

const docker = new Docker();
const slack = new Slack({
  username: 'docker',
  iconEmoji: ':whale:',
});

async function sendEvent(e) {
  console.info(e);
  switch (e.status) {
    case 'start':
      await slack.sendAttachment({
        color: 'good',
        text: 'Container is running',
        fields: [
          {title: 'Image', value: e.from},
          {title: 'Container Name', value: e.Actor.Attributes.name},
          {title: 'Container ID', value: e.id},
        ]
      });
      break;

    // TODO
    case 'kill':
    case 'die':
    case 'destroy':
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
