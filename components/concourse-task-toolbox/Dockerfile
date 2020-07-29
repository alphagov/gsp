FROM hashicorp/terraform:0.12.29 AS task-toolbox

RUN apk add --update \
	curl \
	git \
	wget \
	unzip \
	jq \
	openssh \
	ruby \
	bash \
	openssl \
	file \
	tar \
	netcat-openbsd \
	groff \
	less \
	python3 \
	py3-pip \
	mailcap \
	ncurses \
	gnupg \
	rpm \
	&& pip3 install awscli s3cmd yq PyYAML kubernetes \
	&& rm /var/cache/apk/*

# install kubeseal
ENV KUBESEAL_VERSION=0.9.0
RUN wget https://github.com/bitnami-labs/sealed-secrets/releases/download/v${KUBESEAL_VERSION}/kubeseal-linux-amd64 \
	&& chmod +x /kubeseal-linux-amd64 \
	&& mv kubeseal-linux-amd64 /usr/local/bin/kubeseal

# install kubectl
ENV KUBECTL_VERSION=1.14.0
RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/v${KUBECTL_VERSION}/bin/linux/amd64/kubectl \
	&& chmod +x ./kubectl \
	&& mv ./kubectl /usr/local/bin/kubectl

# install fly 
ENV FLY_VERSION=4.2.1
RUN wget https://github.com/concourse/concourse/releases/download/v${FLY_VERSION}/fly_linux_amd64 \
	&& chmod +x /fly_linux_amd64 \
	&& mv fly_linux_amd64 /usr/local/bin/_fly

# install helm
ENV HELM_VERSION=2.13.1
RUN wget https://storage.googleapis.com/kubernetes-helm/helm-v${HELM_VERSION}-linux-amd64.tar.gz \
	&& tar -zxvf helm-v${HELM_VERSION}-linux-amd64.tar.gz \
	&& mv linux-amd64/helm /bin/helm \
	&& rm -rf linux-amd64

# install sonobuoy
ENV SONOBUOY_VERSION=0.14.3
RUN wget https://github.com/heptio/sonobuoy/releases/download/v${SONOBUOY_VERSION}/sonobuoy_${SONOBUOY_VERSION}_linux_amd64.tar.gz \
	&& tar -zxvf sonobuoy_${SONOBUOY_VERSION}_linux_amd64.tar.gz \
	&& mv sonobuoy /usr/local/bin/sonobuoy \
	&& rm -rf LICENSE

# install kapp
ENV KAPP_VERSION=0.12.0
RUN wget https://github.com/k14s/kapp/releases/download/v${KAPP_VERSION}/kapp-linux-amd64 \
	&& mv kapp-linux-amd64 /bin/kapp \
	&& chmod +x /bin/kapp

# install spruce
ENV SPRUCE_VERSION=1.20.0
RUN wget https://github.com/geofffranks/spruce/releases/download/v${SPRUCE_VERSION}/spruce-linux-amd64 \
	&& mv spruce-linux-amd64 /bin/spruce \
	&& chmod +x /bin/spruce

# install aws-iam-authenticator
RUN curl -LO https://github.com/kubernetes-sigs/aws-iam-authenticator/releases/download/0.4.0-alpha.1/aws-iam-authenticator_0.4.0-alpha.1_linux_amd64 \
	&& chmod +x aws-iam-authenticator_0.4.0-alpha.1_linux_amd64 \
	&& mv aws-iam-authenticator_0.4.0-alpha.1_linux_amd64 /usr/local/bin/aws-iam-authenticator

COPY bin/aws-assume-role bin/setup-kube-deployer bin/determine-platform-version.py bin/findCVEs.py /usr/local/bin/

# --------------------output------------------------

FROM task-toolbox
COPY --from=hairyhenderson/gomplate:v3.5.0-slim /gomplate /bin/gomplate
COPY --from=aquasec/trivy@sha256:3cd0f92134bcdac0b00598f275b1d112312341f5a3d19cc492cecfef39e67c42 /usr/local/bin/trivy /usr/local/bin/trivy
ENTRYPOINT ["/bin/bash"]
CMD []

