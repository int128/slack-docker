FROM alpine:3.11
EXPOSE 3000
ARG TARGETPLATFORM
ARG VERSION
RUN wget -O slack-docker.zip "https://github.com/int128/slack-docker/releases/download/${VERSION}/slack-docker_$(echo "${TARGETPLATFORM}" | cut -d'/' -f1)_$(echo "${TARGETPLATFORM}" | cut -d'/' -f2).zip" && \
    unzip slack-docker.zip
CMD ["/bin/sh"]