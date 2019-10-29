# Terraform 0.12.12
FROM ljfranklin/terraform-resource@sha256:15eee04112da38c0fcbdb9edb86a6b5acff4a800f21cb29b4e30dc58b27e5d0d

# we need the aws tools and git in the box for some of the local-exec scripts
RUN apk add --update jq python3 py3-pip git terraform zip && \
    pip3 install --upgrade pip && \
    pip3 install awscli && \
    rm /var/cache/apk/* && \
    git config --system credential.helper '!aws codecommit credential-helper $@' && \
    git config --system credential.UseHttpPath true
