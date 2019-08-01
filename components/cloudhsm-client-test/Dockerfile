FROM ubuntu:xenial

RUN apt-get update && \
    apt-get install -y curl \
                netcat \
                python3 \
                redis-tools \
                telnet \
                vim \
                wget && \
    wget https://s3.amazonaws.com/cloudhsmv2-software/CloudHsmClient/Xenial/cloudhsm-client_latest_amd64.deb && \
    dpkg -i cloudhsm-client_latest_amd64.deb; apt-get install -f -y && \
    apt-get clean
