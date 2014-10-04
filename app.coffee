if not process.env.domain or not process.env.token
  throw 'Run with mandatory environment variables: domain=DOMAIN token=TOKEN'

Slack       = require 'node-slack'
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

slack       = new Slack process.env.domain, process.env.token
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
        notify container.Name, "Started #{event.id.substring 0, 8} from #{event.from} at #{container.NetworkSettings?.IPAddress}"
    when 'die', 'kill'
      containers.get event.id, (container) ->
        notify container.Name, "Stopped #{event.id.substring 0, 8}"
    when 'destroy'
      containers.getAndRemove event.id, (container) ->
        notify container.Name, "Removed #{event.id.substring 0, 8}"

notify = (name, text) ->
  slack.send
    username: "docker#{name}"
    icon_emoji: ':whale:'
    channel: '#' + process.env.channel || 'general'
    text: text
  .then (->), ((error) -> console.error error)
