class SlackDocker < Formula
  desc "Slack integration to notify Docker events"
  homepage "https://github.com/int128/slack-docker"
  url "https://github.com/int128/slack-docker/releases/download/{{ env "VERSION" }}/slack-docker_darwin_amd64.zip"
  version "{{ env "VERSION" }}"
  sha256 "{{ sha256 .darwin_amd64_archive }}"

  def install
    bin.install "slack-docker"
  end

  test do
    system "#{bin}/slack-docker --help"
  end
end
