# Global Accelerator Spike

## What

Exploring use (and ultimate deprecation) of AWS Global Accelerator to provide
static IPs for ingress. 

## Provisioning and wiring up Global Accelerator to ingress

We wanted to run through the process of provisioning and ensure routing works using the canary's exposed metrics endpoint.

We decided to do this in the console as the terraform provider does not support all the GlobalAccelerator resource types at this time:

* A GlobalAccelerator was provisioned in the console
* listeners setup for TCP ports 80 and 443
* endpoint groups created for eu-west-2 for both listeners
* the NLB with tag "sandbox-main" (the canary's NLB) was added as endpoints to the relevant endpoint groups.

The global accelerator provides a DNS name that resolves to two static IPs:

```

;; ANSWER SECTION:
aa310ecc8b9df7e97.awsglobalaccelerator.com. 300	IN A 13.248.152.89
aa310ecc8b9df7e97.awsglobalaccelerator.com. 300	IN A 76.223.26.39
```

The VirtualService would need configuring with another host name to answer to, or you can hit the endpoint with some curl trickery to resolve IP:

```
curl --resolve canary.london.sandbox.govsvc.uk:443:76.223.26.39 https://canary.london.sandbox.govsvc.uk/metrics
```

## Migrating away...

When the day comes that we can allocate IP addresses to the NLB Services we could do one of the following:

#### Option 1: Spend some error budget

* Wait til everyone has gone to bed
* Cause the recreation of the NLB with the new EIP allocation
* Update DNS to point at NLB IPs
* Take the downtime hit of ~10mins (mostly waiting for NLB to get recreated)
* Remove global accelerator

#### Option 2:  Traffic Shaping (aka: Divide & Conquer)

* Stop the autoscaler
* Create a second global accelerator endpoint group / endpoints that target all of EC2 worker nodes on the existing NodePort directly (ie by passing the NLB)
* Turn the traffic dial so that no traffic is going to the NLB and all direct to worker nodes
* Apply the change that will cause the NLB to have static ips (and get recreated)
* Update DNS to point at NLB IPs
* Remove global accelerator

#### And more...

There's probably some other options depending on the appetite for downtime.

## Exploring effects of Service/NLB recreation...

I have not actually found a change to the `Service` definition that actually causes a recreation (other than deleting it!). So I'm not sure under what circumstances we expect to see an NLB get destroyed... likely candidates:

* Accidental renaming of the ingress Service (ie gsp-cluster release name change or something like that)
* Explicit deletion of Service resource using kubectl
* Accidental misconfiguration of istio values.yaml (ie during an upgrade or as part of other istio config changes)

When the NLB is destroyed, it gets disconnected from the GA, and the new one comes back with a different name so requires ClickOpsing back into existence again.

#### Possible mitigation

* disallow deletion of the ingress Services to non-admin roles (using gatekeeper?)
* add some code-comments around part of istio config that might cause problems
* add a test to the release/deploy pipeline to ensure that change doesn't cause ingress downtime (extension of the check canary test maybe?)
