# How to upgrade

In general copy the [instructions provided by AWS](https://docs.aws.amazon.com/eks/latest/userguide/metrics-server.html).

It all nails down to:

```sh
DOWNLOAD_URL=$(curl --silent "https://api.github.com/repos/kubernetes-incubator/metrics-server/releases/latest" | jq -r .tarball_url)
DOWNLOAD_VERSION=$(grep -o '[^/v]*$' <<< $DOWNLOAD_URL)
curl -Ls $DOWNLOAD_URL -o metrics-server-$DOWNLOAD_VERSION.tar.gz
mkdir metrics-server-$DOWNLOAD_VERSION
tar -xzf metrics-server-$DOWNLOAD_VERSION.tar.gz --directory metrics-server-$DOWNLOAD_VERSION --strip-components 1
```

With additional command of:

```sh
mv metrics-server-$DOWNLOAD_VERSION/deploy/1.8+ ./metrics-server
```
