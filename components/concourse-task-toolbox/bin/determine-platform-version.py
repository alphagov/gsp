#!/usr/bin/env python3
import collections
import os
import subprocess
import sys

os.makedirs("platform-version", exist_ok=True)

print("Picking platform version...")
partial_repos = [
    "platform",
    "service-operator-source",
    "concourse-task-toolbox-source",
    "concourse-operator-source",
    "concourse-github-resource-source",
    "concourse-harbor-resource-source"
]

repo_map = collections.Counter()
for partial_repo in partial_repos:
    proc = subprocess.Popen(
        ['git', 'log', '--format=%H'],
        env={'GIT_DIR': f'{partial_repo}/.git'},
        stdout=subprocess.PIPE
    )
    while True:
        line = proc.stdout.readline()
        if not line:
            break
        repo_map[partial_repo] += 1

repo, _ = repo_map.most_common()[0]
with open(f"{repo}/.git/ref") as f:
	commit = f.read()

print(f"Picked {commit}")
with open('platform-version/ref', 'w') as f:
    f.write(commit)