FROM golang:alpine

USER root:root
ENV HOME /root

RUN apk add --update python py-pip git bash && \
    pip install --upgrade pip && \
    pip install awscli && \
	  apk -v --purge del py-pip && \
	  rm /var/cache/apk/*


# install kubectl
ENV KUBECTL_VERSION=1.14.0
RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/v${KUBECTL_VERSION}/bin/linux/amd64/kubectl \
	&& chmod +x ./kubectl \
	&& mv ./kubectl /usr/local/bin/kubectl


RUN git clone git@github.com:alphagov/gsp-acceptance-tests.git

USER nobody:nobody
ENV HOME /tmp
