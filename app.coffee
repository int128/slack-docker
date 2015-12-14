if not process.env.webhook
  throw 'Run with environment variable: webhook=https://hooks.slack.com/services/...'

JSONStream  = require 'JSONStream'
EventStream = require 'event-stream'
Q           = require 'q'

Docker =
  _service: new (require 'dockerode')(socketPath: '/var/run/docker.sock')
  getVersion: ->
    Q.ninvoke @_service, 'version'
  getContainer: (id) ->
    Q.ninvoke @_service.getContainer(id), 'inspect'
  events: ->
    Q.ninvoke(@_service, 'getEvents', {}).then (stream) ->
      stream.pipe JSONStream.parse()

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

EventFilter =
  _imageRegExp: new RegExp(process.env.image_regexp)
  satisfy: (event) ->
    @_imageRegExp.test event.from

EventProcessor =
  _containers: {}
  _start: (event) ->
    Docker.getContainer(event.id).then (container) =>
      @_containers[event.id] = container
      Slack.send container.Name, "Started #{container.Config.Hostname}",
        'Image': event.from
        'IP Address': container.NetworkSettings.IPAddress
        'Path': container.Path
        'Arguments': container.Args
        'Started at': container.State.StartedAt
  _kill: (event) ->
    if container = @_containers[event.id]
      Slack.send container.Name, "Stopped #{container.Config.Hostname}"
  _die: (event) ->
    @_kill event
  _destroy: (event) ->
    if container = @_containers[event.id]
      delete @_containers[event.id]
      Slack.send container.Name, "Removed #{container.Config.Hostname}"
  handle: (event) ->
    @["_#{event.status}"]?.call(@, event)


# main
Docker.getVersion().then (version) ->
  console.info version
  Docker.events().then (stream) ->
    stream.pipe EventStream.map (event) ->
      if EventFilter.satisfy event
        EventProcessor.handle event
.fail (e) ->
  console.error e

