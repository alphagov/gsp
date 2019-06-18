FROM concourse/git-resource:1.1.1

RUN apk update

RUN curl https://github.com/web-flow.gpg -o /tmp/web-flow.gpg && \
    gpg --import /tmp/web-flow.gpg && \
    rm /tmp/web-flow.gpg

RUN mkdir -p /opt/resource/origin && \
    mv /opt/resource/check /opt/resource/check.origin && \
    mv /opt/resource/in /opt/resource/in.origin && \
    mv /opt/resource/out /opt/resource/out.origin

COPY ./assets/check /opt/resource/
COPY ./assets/in /opt/resource/
COPY ./assets/out /opt/resource/
