Slack       = require 'node-slack'
Docker      = require 'dockerode'
JSONStream  = require 'JSONStream'

slack  = new Slack process.env.domain, process.env.token
docker = new Docker socketPath: '/var/run/docker.sock'

notify = (id, text) ->
  slack.send
    username: "Docker #{id.substring 0, 8}"
    icon_emoji: ':whale:'
    channel: process.env.channel || '#general'
    text: text
  .then (->), ((error) -> console.error error)

docker.getEvents {}, (error, stream) ->
  stream?.pipe JSONStream.parse().on 'root', (event) ->
    switch event.status
      when 'start'
        docker.getContainer(event.id).inspect (error, detail) ->
          notify event.id, "Started a container #{detail?.Name} from #{event.from} at #{detail?.NetworkSettings.IPAddress}"
          console.error "#{event.time}: #{error}" if error
      when 'die', 'kill'
        notify event.id, "Stopped the container created from #{event.from}"
      when 'destroy'
        notify event.id, "Removed the container created from #{event.from}"
    console.info "#{event.time}: #{event.status}: #{event.id} from #{event.from}"
  console.error "#{event.time}: #{error}" if error
