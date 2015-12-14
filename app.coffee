if not process.env.webhook
  throw 'Run with environment variable: webhook=https://hooks.slack.com/services/...'

Docker      = require 'dockerode'
JSONStream  = require 'JSONStream'
EventStream = require 'event-stream'
Q           = require 'q'

class NamedMap
  constructor: ->
    @map = {}
  get: (k, f) ->
    f @map[k] if @map[k]
  getAndRemove: (k, f) ->
    f @map[k] if @map[k]
    delete @map[k]
  put: (k, v) ->
    @map[k] = v

docker      = new Docker socketPath: '/var/run/docker.sock'
containers  = new NamedMap


Slack =
  _service: (require 'slack-notify')(process.env.webhook)
  send: (name, text, fields) ->
    Q.ninvoke @_service, 'send',
      username: "docker#{name}"
      icon_emoji: ':whale:'
      channel: ''
      text: text
      fields: fields
    .fail (e) -> console.error e


EventProcessor =
  handle: (event) ->
    switch event.status
      when 'start'
        Q.ninvoke docker.getContainer(event.id), 'inspect'
          .then (container) ->
            containers.put event.id, container
            Slack.send container.Name, "Started #{container.Config.Hostname}",
              'Image': event.from
              'IP Address': container.NetworkSettings.IPAddress
              'Path': container.Path
              'Arguments': container.Args
              'Started at': container.State.StartedAt
      when 'die', 'kill'
        containers.get event.id, (container) ->
          Slack.send container.Name, "Stopped #{container.Config.Hostname}"
      when 'destroy'
        containers.getAndRemove event.id, (container) ->
          Slack.send container.Name, "Removed #{container.Config.Hostname}"


Q.ninvoke docker, 'version'
  .then (version) ->
    console.info version
    Q.ninvoke docker, 'getEvents', {}

  .then (stream) ->
    stream.pipe(JSONStream.parse()).pipe EventStream.map (event) ->
      console.info "#{event.time}: #{event.status}: #{event.id} from #{event.from}"
      EventProcessor.handle event

  .fail (e) ->
    console.error e

