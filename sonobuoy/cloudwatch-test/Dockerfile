FROM gcr.io/google_containers/git-sync:v2.0.6

USER root:root
ENV HOME /root

RUN apk add --update python py-pip && \
    pip install --upgrade pip && \
    pip install awscli && \
    rm /var/cache/apk/*

RUN git clone git@github.com:alphagov/gsp-acceptance-tests.git

USER nobody:nobody
ENV HOME /tmp
