# ADR016: Code Verification

## Status

Accepted

## Context

All of our deployment processes start with code in git repositories hosted on
Github and continuous delivery pipelines automate the build and deployment of
applications out to production. Git supports GPG signing of commits.

Code changes should not be able to make it out to a production environment
without being exposed to at least two authorised code owners.

Our aim is to improve on process-based approval methods by being able to
digitally verify that code has been exposed to at least two authorised code
owners without adversely affecting developer workflow.

Some potential solutions involve:

1. Enforcing that two unique developers each an add empty signed commit at the
   tip of the branch for PRs before being merged, and verifying the existence
   of these empty signed commits as part of the delivery pipeline.
2. Enforcing that two unique developers each sign the commmit at the tip of the
   branch for PRs using
   [git-signatures](https://github.com/hashbang/git-signatures) and verifying
   that at least two trusted developer keys are present as part of the delivery
   pipeline.
3. Verify that at least one (non-PR-authoring) trusted Github user has approved
   the commit from the delivery pipeline using the Github API (rather than the
   web UI, which is vulnerable to manipulation by a single Github "owner").

## Decision

Solutions (1) and (2) both introduce significant changes to existing developer
workflows. The benefit of either of these solutions is that when combined with
an external hardware security device (Yubikey) they can provide greater
protection against developer machine compromise than Github's session duration
alone.  However we don't feel this benefit warrants the reduced usability for
developers across GDS.

We implement solution (3).

## Consequences

* Developer workflow is largely unaffected across GDS by taking advantage of
  Github features team are already familiar
* Compromise of multiple developer laptops could mitigate "two-eyes" process
