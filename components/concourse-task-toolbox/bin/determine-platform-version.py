#!/usr/bin/env python3
import collections
import os
import subprocess

os.makedirs("platform-version", exist_ok=True)

print("Picking platform version...")
partial_repos = [
    "platform",
    "aws-node-lifecycle-hook",
    "service-operator-source",
    "concourse-task-toolbox-source",
    "concourse-operator-source",
    "concourse-github-resource-source",
    "concourse-harbor-resource-source",
    "concourse-terraform-resource-source",
    "aws-ssm-agent-source"
]

commit_map = collections.Counter()
for partial_repo in partial_repos:
    with open(f"{partial_repo}/.git/ref") as f:
        commit = f.read()
    proc = subprocess.Popen(
        # This could probably use 'HEAD' instead of reading .git/ref
        ['git', 'rev-list', '--count', commit.strip()],
        env={'GIT_DIR': f'{partial_repo}/.git'},
        stdout=subprocess.PIPE
    )
    stdoutdata, _ = proc.communicate()
    commit_map[commit] = int(stdoutdata)

commit, _ = commit_map.most_common()[0]

print(f"Picked {commit}")
with open('platform-version/ref', 'w') as f:
    f.write(commit)
