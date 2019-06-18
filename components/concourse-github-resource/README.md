# concourse-github-resource

GitHub resource for Concourse that enforces a minimum number of GitHub approvals. This relies heavily on the [Concourse `git-resource`](https://github.com/concourse/git-resource).

## Source configuration

All the required configuration for the Concourse `git-resource` will be required along with:

* `organization` *Required.* The GitHub organization the repo is in.
* `repository` *Required.* The repository name.
* `github_api_token` *Required.* A GitHub API token.
* `approvers`*Required.* A list GitHub usernames of approvers.
* `required_approval_count` *Required.* The minimum number of approvals required to proceed.

## Run tests
```
rm -rf tmp && cat test.json | docker run -v $PWD/tmp:/mnt/myapp -i $(docker build -q .) /opt/resource/in /mnt/myapp
```
