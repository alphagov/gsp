FROM amazonlinux:2

RUN yum update -y && \
    yum install -y systemd curl tar sudo && \
    yum install -y https://s3.amazonaws.com/ec2-downloads-windows/SSMAgent/latest/linux_amd64/amazon-ssm-agent.rpm

WORKDIR /opt/amazon/ssm/

CMD ["amazon-ssm-agent", "start"]
