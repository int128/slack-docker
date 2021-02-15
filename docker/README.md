## Building a multi-platform Docker image

- leverages an experimental Docker CLI feature - `buildx`
- assumes that Docker images are built _after_ the binaries are released as a Github Release; this is currently taken care of by [a CircleCI pipeline that runs on every tag](../.circleci/config.yml).


https://www.docker.com/blog/multi-platform-docker-builds/
https://docs.docker.com/buildx/working-with-buildx/

```
export VERSION=$(git describe --tags)
docker buildx build -t int128/slack-docker:${VERSION} --push --platform linux/amd64,linux/arm64 --build-arg VERSION=${VERSION} .
```