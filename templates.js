module.exports.container_start = e => ({
  color: 'good',
  text: `Started container \`${e.Actor.Attributes.name}\``,
  mrkdwn_in: ['text'],
  fields: [
    {short: true, title: 'Image', value: e.Actor.Attributes.image},
    {title: 'Container ID', value: e.Actor.ID},
  ]
});

module.exports.container_kill = e => ({
  color: e.Actor.Attributes.exitCode == 0 ? 'good' : 'danger',
  text: `Stopped container \`${e.Actor.Attributes.name}\``,
  mrkdwn_in: ['text'],
  fields: [
    {short: true, title: 'Image', value: e.Actor.Attributes.image},
    {short: true, title: 'Exit Code', value: e.Actor.Attributes.exitCode},
    {title: 'Container ID', value: e.Actor.ID},
  ]
});

module.exports.container_die = module.exports.container_kill;

module.exports.container_destroy = e => ({
  color: 'warning',
  text: `Destroyed container \`${e.Actor.Attributes.name}\``,
  mrkdwn_in: ['text'],
  fields: [
    {short: true, title: 'Image', value: e.Actor.Attributes.image},
    {title: 'Container ID', value: e.Actor.ID},
  ]
});

module.exports.network_create = e => ({
  color: 'warning',
  text: `Created network \`${e.Actor.Attributes.name}\``,
  mrkdwn_in: ['text'],
});

module.exports.network_destroy = e => ({
  color: 'warning',
  text: `Destroyed network \`${e.Actor.Attributes.name}\``,
  mrkdwn_in: ['text'],
});
