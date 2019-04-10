# Setting up HSM

We're intending to use hardware security modules (HSM) as the means to sign a
given piece of information with a non-extractable key.

At the moment, we're committed to using CloudHSM provided by AWS.


## Signing key

We'd like to use another HSM-like piece of hardware, to generate an
RSA key for signing the certificates used by HSM for trust links.

The generation of that key should happen on the YubiKey itself.

You can either use the [YubiKey Manager GUI](https://www.yubico.com/products/services-software/download/yubikey-manager/) or the `yubico-piv-tool` CLI:

```sh
yubico-piv-tool -s 9c -A RSA2048 -H SHA256 -a generate
```

## Creating HSM

We use our [HSM terraform module](https://github.com/alphagov/gsp-teams/blob/master/terraform/modules/hsm/main.tf)
to create the instance next to our GSP cluster.

## Initialising

We should follow the [AWS CloudHSM - Initialise the Cluster](https://docs.aws.amazon.com/cloudhsm/latest/userguide/initialize-cluster.html)
documentation for the most part.

## Signing the HSM's CSR

It's particularly tricky to sign with the use of the YubiKey. In the end, it
worked on the Ubuntu 18.04 (bionic) with the following set:

```sh
sudo apt install build-essential

sudo apt install pcscd yubikey-manager yubico-piv-tool yubikey-personalization libpcsclite-dev

sudo apt install libengine-pkcs11-openssl opensc opensc-pkcs11

## The following paths (SO_PATH and MODULE_PATH) are very fragile.
openssl << EOF
engine dynamic -pre SO_PATH:/usr/lib/x86_64-linux-gnu/engines-1.1/pkcs11.so -pre ID:pkcs11 -pre NO_VCHECK:1 -pre LIST_ADD:1 -pre LOAD -pre MODULE_PATH:/usr/lib/x86_64-linux-gnu/opensc-pkcs11.so -pre VERBOSE
x509 -engine pkcs11 -CAkeyform engine -CAkey slot_0-id_2 -CAcreateserial -sha256 -CA customerCA.crt -req -in cluster-rn3c7izgvya_ClusterCsr.csr -out cluster-rn3c7izgvya_CustomerHsmCertificate.crt
EOF
```

The `openssl` command mentions several files:
- `customerCA.crt` is basically the CA certificate exported from the key.
  Easily obtainable from GUI or the `yubico-piv-tool`
- `cluster-rn3c7izgvya_ClusterCsr.csr` is the file we obtain from HSM itself,
  following [these instructions](https://docs.aws.amazon.com/cloudhsm/latest/userguide/initialize-cluster.html#get-csr)
- `cluster-rn3c7izgvya_CustomerHsmCertificate.crt` will be generated upon
  successful signing

The file names are bit random. They represent the HSM's ID.

In general, this post was better than any other documentation I found:
https://dennis.silvrback.com/openssl-ca-with-yubikey-neo

After we have signed the certificate, we can continue with the initialisation
and activation of the cluster. We should be following the AWS guide on any of
this, except the certificate signing which we do our way.

### Push the customerCA.crt to AWS SSM

This signed customer certificate will need to be shared with clients
interacting with the HSM.

We decided to store it along the passwords in HSM for easier access.

```sh
aws ssm put-parameter --type SecureString --name /hsm/customerCA --value "$(cat customerCA.crt)"
```

Obtainable with:

```sh
aws ssm get-parameter --query Parameter.Value --output text --with-decryption --name /hsm/customerCA
```

## Activating the HSM

The HSM cluster has ingress security groups, meaning we can only
get onto it from specific boxes; k8s EC2 Nodes in our case.

```sh
kubectl run ubuntu -it --rm --image ubuntu:trusty -- bash

## On the box follow AWS guide to install client: https://docs.aws.amazon.com/cloudhsm/latest/userguide/install-and-configure-client-linux.html

## Follow the Activate Cluster AWS docs: https://docs.aws.amazon.com/cloudhsm/latest/userguide/activate-cluster.html
```

The activation of the cluster involves password changing and user creation.
Once again, [AWS docs](https://docs.aws.amazon.com/cloudhsm/latest/userguide/activate-cluster.html)
and client `help` should be sufficient to accomplish that.

The passwords are generated on Programme's AWS SSM, with the use of:

```sh
aws ssm put-parameter --type SecureString --name /hsm/users/co/password --value "$(pwgen --capitalize --numerals --secure --symbols -1 --ambiguous 32 1 | tr -d '\n')" ## For Crypto Officer

aws ssm put-parameter --type SecureString --name /hsm/users/cu/1/password --value "$(pwgen --capitalize --numerals --secure --symbols -1 --ambiguous 32 1 | tr -d '\n')" ## For Crypto User - Increased by 1 for any next user.
```

These passwords can be retrieved with:

```sh
aws ssm get-parameter --query Parameter.Value --output text --with-decryption --name /hsm/users/cu/1/password | tr -d '\n' | pbcopy
```

## Push the HSM IP to SSM

Like wise everything else, it would be nice to have the HSM IP accessible from
the same place.

```sh
aws ssm put-parameter --type SecureString --name /hsm/ip --value "10.0.12.205"
```

Obtain with:

```sh
aws ssm get-parameter --query Parameter.Value --output text --with-decryption --name /hsm/ip
```
