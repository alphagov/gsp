# pinned to a beta-backend version of ljfranklin/terraform-resource
# FIXME: once beta-backend is promoted to master, use propper version labels
FROM ljfranklin/terraform-resource@sha256:688b879952a041c176087973bdefb6d40bc14a0028f69d4576764d8b2b67501d

# we need the aws tools and git in the box for some of the local-exec scripts
RUN apk add --update jq python3 py3-pip git terraform zip && \
    pip3 install --upgrade pip && \
    pip3 install awscli && \
    rm /var/cache/apk/* && \
    git config --system credential.helper '!aws codecommit credential-helper $@' && \
    git config --system credential.UseHttpPath true
