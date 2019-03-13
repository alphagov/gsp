# Incident Report - February 4th, 2019

* This is primarily for the GSP team, and will also be available to other GDS teams
* Keep report factual - avoid sentiment/blame, and no names except to record the meeting
* All decisions, events and actions should be recorded in the timeline

| Date | February 4th, 2019 |
| --- | ---: |
| Time | 0730 - 1230 (UTC) |
| Application/process: | Verify clusters |

## Overview

### How users were affected 

* Verify unable to run build pipelines
* The Pen-testers were unable to run their test against application

### Timeline

_All times in UTC, unless otherwise stated._

**Types:** `automated`, `scheduled`, `app event`, `notification`, `human notification`, `operator`, `comms`

| Time | Event Type | Description |
| ---- | ---------- | ----------- |
| 0730 | `operator` | The operator has decided to renew the metadata in order to allow Pen-Test to continue. |
| 0740 | `operator` | The operator has noticed that the containers are incapable of starting up on their own. Starts looking into it. |
| 0750 | `operator` | The operator has triggered a pipeline to see it can successfully re-provision a cluster. |
| 0809 | `operator` | The operator notices the cluster is unable to provision properly due to CodeCommit 403 error. |
| 0810 | `operator` | The operator starts to  investigate potential cause of failure. |
| 0836 | `human notification` | A pen-tester has sent an email notifying they're unable to access the application. |
| 0851 | `human notification` | A tenant prods the operator about the issue to ensure it's being handled. Is reasured the work on the fix is in progress. |
| 0853 | `operator` | The operator starts to implement changes and creates PRs to fix issues caused by semi-tested feature. |
| 0907 | `operator` | The operator prioritises getting the service back up and running and decides to directly merge needed changes. - BAD |
| 0942 | `automated` | The operator triggers `destroy` pipeline for the `tools` cluster. Hopes that it would trigger setting up new one shortly. |
| 0953 | `scheduled` | The `tools` cluster fails to spin back up due to invalid/missing arguments in terraform config. |
| 1000 | `operator` | The operator digs into number of errors in an attpemt to fix the issue by roll back. |
| 1030 | `operator` | The operator after digging deep enough into a rabbit hole, decides to stop and revert all the "hacks" that have been made thus far. |
| 1115 | `comms` | The elected `comms` person notifies the pen-testers that the issue is being worked on. |
| 1138 | `operator` | The operator decides to ditch an attempt at properly fixing the configuration and decides to build images locally and push them to public registry. |
| 1146 | `operator` | The operator logs into staging cluster as an admin and manually changes the source of docker images. |
| 1200 | `operator` | The operator executes the command to roll the containers on the cluster. |
| 1210 | `operator` | The operator witnesses pods being stabilised and runs a journey in a browser. |
| 1230 | `operator` | The operator has finished restoring to stable state and calls the end of incident. Continues investigating issue in peace. |

## Post-mortem

_To be filled in at the time of incident retro_

| Date | February 19th, 2019 |
| --- | ---: |
| Time | 1500 |
| Attendance: | Chris M, Paul D, Daniel B, Veronika K, Joshua K, David M, David P, Rafal, Chris F, Stephen F |

### Root cause

#### What caused the original problem?

A Kiam feature restricting access to AWS resources by the system, caused accessing artifacts impossible.

#### Why wasn’t this caught before it reached prod?

The pipeline is set to always read the latest changes from master branch.

* We made we a concious decision to release the Kiam features (that lock down IAM access from pods), as we considered it easier to fix the problems it causes rather than preempt them
* We have no acceptance tests for infrastructure changes to validate a "staged" change before promoting to production

#### Were we alerted quickly?

No. This presumably has been an issue for few days, but no-one needed to use the service.

We have no alerts configured at this time.

#### Were we able to diagnose and fix the immediate issue quickly?

* Identifying the issue with the registry was achieved promptly once the the issue was raised
* Applying the initial fix was hindered by not being able to decide if the issue is an "incident" and how to prioritize the fix

#### Was the process followed well, were comms effective?

This wasn't treated as an incident at the time, no comms were elected initially. Later on, the comms have been appreciated and effective.

### Possible Improvements:

#### What we could do to prevent this happening again

* Cluster acceptance tests on a separate cluster
* Pinning terraform modules to a specific version

#### What we could do to improve our response

* Formalizing our incident process

#### What we could do to improve comms/process

* As above

### Actions

* Define what support and incident process looks like : Shared
* Define acceptance tests : Chris F create story
* Pin terraform modules : Existing story in backlog
* Define promotion pipeline : Chris F create story

### Recommendations:

_These are recommendations for other teams which fix/prevent/improve this issue_

* N/A

### Comments and questions (can remove after discussion):

* N/A
