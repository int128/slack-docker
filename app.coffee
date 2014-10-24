if not process.env.webhook
  throw 'Run with environment variable: webhook=https://hooks.slack.com/services/...'

Slack       = require 'slack-notify'
Docker      = require 'dockerode'
JSONStream  = require 'JSONStream'

class NamedMap
  constructor: ->
    @map = {}
  get: (k, f) ->
    f @map[k]
  getAndRemove: (k, f) ->
    f @map[k]
    delete @map[k]
  put: (k, v) ->
    @map[k] = v

slack       = Slack process.env.webhook
docker      = new Docker socketPath: '/var/run/docker.sock'
containers  = new NamedMap

docker.version (error, version) ->
  throw error if error
  console.info version
  docker.getEvents {}, (error, stream) ->
    throw error if error
    stream?.pipe JSONStream.parse().on 'root', handle

handle = (event) ->
  console.info "#{event.time}: #{event.status}: #{event.id} from #{event.from}"
  switch event.status
    when 'start'
      docker.getContainer(event.id).inspect (error, container) ->
        throw error if error
        containers.put event.id, container
        notify container.Name, "Started #{container.Config.Hostname} from #{event.from} at #{container.NetworkSettings.IPAddress}"
    when 'die', 'kill'
      containers.get event.id, (container) ->
        notify container.Name, "Stopped #{container.Config.Hostname}"
    when 'destroy'
      containers.getAndRemove event.id, (container) ->
        notify container.Name, "Removed #{container.Config.Hostname}"

notify = (name, text) ->
  slack.send
    username: "docker#{name}"
    icon_emoji: ':whale:'
    channel: '#' + process.env.channel || 'general'
    text: text

slack.onError = (e) -> console.error e
