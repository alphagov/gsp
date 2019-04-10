# Incident Report - February 27th, 2019

* This is primarily for the GSP team, and will also be available to other GDS teams
* Keep report factual - avoid sentiment/blame, and no names except to record the meeting
* All decisions, events and actions should be recorded in the timeline

| Date | February 27th, 2019 |
| --- | ---: |
| Time | 2019-02-27 06:45 - 2019-02-28 13:00 (UTC) |
| Application/process: | Verify tools cluster |

## Overview

### How users were affected 

* No application builds could proceed.
* No application deployments (from either the eIDAS side or from GSP side, where there was a dependency on the tools cluster).

### Timeline

_All times in UTC, unless otherwise stated._

**Types:** `automated`, `scheduled`, `app event`, `notification`, `human notification`, `operator`, `comms`

| Time | Event Type | Description |
| ---- | ---------- | ----------- |
| 2019-02-27 06:45 | `automated` | AWS initiated a reboot of the EC2 instance running the tools cluster control plane.|
| 10:49 | `operator` | Noticed the Verify CD pipeline was red in Concourse. |
| 13:30 | `operator` | Decided to treat it as an incident and investigate. |
| 13:45 | `operator` | Discovered the control plane was dead. |
| 14:25 | `operator` | Attempted to re-bootstrap tools cluster. |
| 14:45 | `operator` | Control plane came up. Problems with worker networking. |
| 16:15 | `operator` | Destroyed & rebuilt tools cluster with more control plane nodes. |
| 17:00 | `operator` | Discovered problems with kiam agents and harbor components. |
| 17:15 - 23:00 | `operator` | Discovered coherent versions of Harbor components. |
| 2019-02-28 10:15 | `operator` | Successful tests on tools cluster. |
| 11:15 | `operator` | Destroyed and rebuilt tools cluster. |
| 11:45 - 13:00 | `operator` | Manual fixing of keys and Notary state. |
| 13:00 | `automated` | Builds succeeded. |

## Post-mortem

_To be filled in at the time of incident retro_

| Date | 2019-03-13 |
| --- | ---: |
| Time | 11:00 |
| Attendance: | Chris Farms, Rafal Proszowski, Dave Povey, Joshua Kuforiji, Iain Baars-Gordon, Daniel Blair, David Pye, Paul Dougan |

### Root cause

#### What caused the original problem?

Several cascading factors:
* Having a single-node control plane
* Having a single worker
* Not pinning versions of 3rd party dependencies (Harbor, in this case)

#### Why wasn’t this caught before it reached prod?

N/A

#### Were we alerted quickly?

We were not alerted at all.

#### Were we able to diagnose and fix the immediate issue quickly?

Diagnosis was quick. It will not be this quick in future as SSH access to the control plane was required.

We were unable to fix the issue. We abandoned trying to fix it after 3 hours and destroyed and respun the cluster from scratch.

#### Was the process followed well, were comms effective?

The incident was declared late. But once we followed the process a slack thread was kept up to date with latest news.

### Possible Improvements:

#### What we could do to prevent this happening again

* multi-node control plane
* chaos engineering
* pin dependencies to specific versions
* externalise control plane management (e.g. EKS)
* verify control plane replicas are spread across multiple nodes

#### What we could do to improve our response

* monitoring the api availability
* Game Days
* Run books


#### What we could do to improve comms/process

* publish "status page" for clusters

### Actions

* Move incident reports to team manual [DB]
* Add basic chronitor / pingdom healthcheck [CF to write story] 
* Make our incident handling process official [IBG]
* Create a dedicated incident slack channel (#gsp-incident) [IBG] [done]
* Investigate a statuspage account [IBG]

### Recommendations:

* move to some managed kubernetes service
* prioritise monitoring & alerting

### Comments and questions (can remove after discussion):
