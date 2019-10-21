#!/usr/bin/env python3
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
commits_known_by = {open(f"{r}/.git/ref").read(): [] for r in partial_repos}

for partial_repo in partial_repos:
    proc = subprocess.Popen(
        ['git', 'log', '--format=%H'],
        env={'GIT_DIR': f'{partial_repo}/.git'},
        stdout=subprocess.PIPE
    )
    seen_commits = set()
    while len(seen_commits) < len(list(commits_known_by.keys())):
        line = proc.stdout.readline().decode('ascii')
        if not line:
            break
        if line in commits_known_by:
            seen_commits.update([line])
            commits_known_by[line].append(partial_repo)

# Pick the only commit which has one repo knowing about it
obscure_commits = list(filter(lambda i: len(i[1]) == 1, commits_known_by.items()))
if len(obscure_commits) != 1:
    sys.stderr.write(f"More or less than one commit was found in a single repository\n")
    sys.exit(1)

(commit, _), = obscure_commits

print(f"Picked {commit}")
with open('platform-version/ref', 'w') as f:
    f.write(commit)