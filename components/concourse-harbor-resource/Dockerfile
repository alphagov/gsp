FROM concourse/docker-image-resource:1.3.0

RUN apk update && apk add curl openssl

RUN curl -L -o /usr/local/bin/notary https://github.com/theupdateframework/notary/releases/download/v0.6.1/notary-Linux-amd64 && \
    chmod +x /usr/local/bin/notary && \
    echo -n '73353b2b4b85604c738a6800465133cb3a828dff0aa26f3c0926dd9a73e19879  /usr/local/bin/notary' | sha256sum -c -s

RUN mkdir -p /opt/resource/origin && \
    mv /opt/resource/check /opt/resource/check.origin && \
    mv /opt/resource/in /opt/resource/in.origin && \
    mv /opt/resource/out /opt/resource/out.origin

RUN sed -i 's/docker build/DOCKER_CONTENT_TRUST=0 docker build/g' /opt/resource/out.origin && \
    sed -i 's/tail -1 push-output/grep digest push-output/g' /opt/resource/out.origin

COPY ./assets/check /opt/resource/
COPY ./assets/in /opt/resource/
COPY ./assets/out /opt/resource/
