# Adding a custom domain to your GSP cluster

You may wish to host your GSP apps on a custom domain.  This guide
explains how it works.

## Assumptions

We assume that:

 - you want GSP to automatically provision TLS certificates (rather
   than providing your own)
 
## Telling the cluster about the domain

To tell a cluster that you want it to manage a particular domain, add
the name to the `extra-zones` value in the cluster-config repository.
For example:

    extra-zones: ["my-domain.example.com"]

## Delegating the domain to the cluster

Once you have told the cluster about your custom domain, you need to
delegate the domain to the cluster.  Currently, the way you do this is
to go into the AWS console to discover the appropriate NS records for
the zone, then set these NS records in your DNS provider to delegate
that domain to the cluster.

## Caveats

CNAMEs are not supported.  This is because we use the
[`tls-dns-01`][acme-challenges] challenge type to verify domains for
issuing certificates, and to do this we need to be able to create TXT
records on the domain.

[acme-challenges]: https://letsencrypt.org/docs/challenge-types/
