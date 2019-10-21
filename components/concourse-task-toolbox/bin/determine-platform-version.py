#!/usr/bin/env python3
import os
import subprocess
import sys

os.mkdirs("platform-version", exist_ok=True)

print("Picking platform version...")
partial_repos = [
    "platform",
    "service-operator-source",
    "concourse-task-toolbox-source",
    "concourse-operator-source",
    "concourse-github-resource-source",
    "concourse-harbor-resource-source"
]
hashes = set(map(lambda r: open(f"{r}/.git/ref").read(), partial_repos))

proc = subprocess.Popen(
    ['git', 'log', '--format=%H'],
    env={'GIT_DIR': 'full-gsp/.git'},
    stdout=subprocess.PIPE
)
while True:
    line = proc.stdout.readline().decode('ascii')
    if not line:
        sys.stderr.write("Did not get line when running git-log through full-gsp\n")
        sys.exit(1)
    if line in hashes:
        proc.terminate()
        break

line = line.strip()
print(f"Picked {line}")
with open('platform-version/ref', 'w') as f:
    f.write(line)