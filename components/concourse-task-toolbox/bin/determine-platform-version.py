#!/usr/bin/env python3
import collections
import os
import subprocess

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
        ['git', 'rev-list', '--count', 'HEAD'],
        env={'GIT_DIR': f'{partial_repo}/.git'},
        stdout=subprocess.PIPE
    )
    stdoutdata, _ = proc.communicate()
    repo_map[partial_repo] = int(stdoutdata)

repo, _ = repo_map.most_common()[0]
with open(f"{repo}/.git/ref") as f:
    commit = f.read()

print(f"Picked {commit}")
with open('platform-version/ref', 'w') as f:
    f.write(commit)