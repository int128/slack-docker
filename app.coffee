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

class Containers
  constructor: -> @_map = {}
  find: (k, f) -> f @_map[k] if @_map[k]
  put: (k, v) -> @_map[k] = v
  remove: (k) -> delete @_map[k]

EventProcessor =
  _containers: new Containers()
  handle: (event) ->
    switch event.status
      when 'start'
        Docker.getContainer(event.id).then (container) =>
          @_containers.put event.id, container
          Slack.send container.Name, "Started #{container.Config.Hostname}",
            'Image': event.from
            'IP Address': container.NetworkSettings.IPAddress
            'Path': container.Path
            'Arguments': container.Args
            'Started at': container.State.StartedAt
      when 'die', 'kill'
        @_containers.find event.id, (container) =>
          Slack.send container.Name, "Stopped #{container.Config.Hostname}"
      when 'destroy'
        @_containers.find event.id, (container) =>
          @_containers.remove event.id
          Slack.send container.Name, "Removed #{container.Config.Hostname}"


# main
Docker.getVersion().then (version) ->
  console.info version
  Docker.events().then (stream) ->
    stream.pipe EventStream.map (event) ->
      console.info "#{event.time}: #{event.status}: #{event.id} from #{event.from}"
      EventProcessor.handle event
.fail (e) ->
  console.error e

