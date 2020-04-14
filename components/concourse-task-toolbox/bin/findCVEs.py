#!/usr/bin/python3

import subprocess
import json
import sys

from kubernetes import client, config
# whitelists against images that are problematic to pull/scan
GLOBAL_IMAGE_WHITELIST = [
    'istio/mixer:1.3.5', # error in image scan: scan failed: failed to apply layers: unknown OS - no shell, no ls - possibly scratch
    'jaegertracing/all-in-one:1.14', # error in image scan: scan failed: failed to apply layers: unknown OS - no shell, no ls - possibly scratch
    'k8s.gcr.io/kubernetes-dashboard-amd64:v1.10.1', # error in image scan: scan failed: failed to apply layers: unknown OS - no shell, no ls - possibly scratch
    'k8s.gcr.io/metrics-server-amd64:v0.3.0', # error in image scan: scan failed: failed to apply layers: unknown OS
]
GLOBAL_IMAGE_SOURCE_WHITELIST = [
    '.dkr.ecr.eu-west-2.amazonaws.com/', # ECR
    '.dkr.ecr.us-west-2.amazonaws.com/', # ECR - for EKS upstream
    'registry.london.verify.govsvc.uk', # TODO: why do we have these old references in sandbox-proxy-node-dev?
    'quay.io/', # https://github.com/aquasecurity/trivy/issues/401
]

# whitelists against vulnerabilities we've considered for various reasons
ISTIO_1_3_5_VULNERABILITIES_WHITELIST = [
    'CVE-2018-20961',
    'CVE-2019-14287',
    'CVE-2019-14896',
    'CVE-2019-14901',
    'CVE-2019-15505',
    'CVE-2019-15926'
]
FLUENTD_1_7_3_VULNERABILITIES_WHITELIST = [
    'CVE-2019-10220',
    'CVE-2019-14896',
    'CVE-2019-14901',
    'CVE-2019-15505',
    'CVE-2019-20636',
]


def whitelisted(vulnerability):
    if vulnerability['image_name'].startswith('istio/') and \
       vulnerability['image_name'].endswith(':1.3.5') and \
       vulnerability['vulnerability']['VulnerabilityID'] in ISTIO_1_3_5_VULNERABILITIES_WHITELIST:
        # these should be fixed in :1.5.1.
        return True
    if vulnerability['image_name'] == 'fluent/fluentd-kubernetes-daemonset:v1.7.3-debian-cloudwatch-1.0' and \
       vulnerability['vulnerability']['VulnerabilityID'] in FLUENTD_1_7_3_VULNERABILITIES_WHITELIST:
        # these should be be fixed in:
        # fluent/fluentd-kubernetes-daemonset:v1.9.3-debian-cloudwatch-1.0
        return True
    if vulnerability['image_name'].startswith('fluent/fluentd-kubernetes-daemonset:v1.') and \
       vulnerability['vulnerability']['VulnerabilityID'] == 'CVE-2020-8130':
        # this shows up in usr/local/bundle/gems/async-http-0.50.0/examples/fetch/Gemfile.lock -
        # which is just an example in one of the libraries, and also in
        # usr/local/bundle/gems/http_parser.rb-0.6.0/Gemfile.lock
        # The second one is slightly more concerning but the nature of the vulnerability appears
        # to be unwanted behaviour from some internal functions of a build library, which seems
        # unlikely to pose a real problem for us.
        # In https://hackerone.com/reports/651518 it was written:
        # "the attack surface was limited because if It's difficult to inject malicious input to
        # Rake::FileList by attackers with the current usage of Rake in the world."
        return True
    return False

trivy_cache = {}
config.load_kube_config()
vulnerabilities = []
for pod in client.CoreV1Api().list_pod_for_all_namespaces(watch=False).items:
    for container in pod.spec.containers:
        image_name = container.image.replace('docker.io/', '')
        if image_name in GLOBAL_IMAGE_WHITELIST:
            continue
        if any(source in image_name for source in GLOBAL_IMAGE_SOURCE_WHITELIST):
            continue
        if image_name not in trivy_cache:
            trivy_cache[image_name] = []
            data = json.loads(subprocess.check_output([
                'trivy',
                '--format', 'json',
                '--quiet',
                '--ignore-unfixed', # remove this if you want to learn about CVE-2005-2541
                '-s', 'CRITICAL',
                image_name
            ]))
            for target in data:
                trivy_cache[image_name] += target.get('Vulnerabilities') or []
        for trivy_vulnerability_obj in trivy_cache[image_name]:
            vulnerability = {
                'namespace': pod.metadata.namespace,
                'container_name': container.name,
                'image_name': image_name,
                'vulnerability': trivy_vulnerability_obj,
            }
            # de-duplicate multiple pods belonging to the same ReplicaSet/StatefulSet/DaemonSet etc. by attributing to their owning object
            if len(pod.metadata.owner_references) > 0:
                assert len(pod.metadata.owner_references) == 1
                vulnerability['kind'] = pod.metadata.owner_references[0].kind
                vulnerability['name'] = pod.metadata.owner_references[0].name
            else:
                vulnerability['kind'] = 'Pod'
                vulnerability['name'] = pod.metadata.name
            if whitelisted(vulnerability):
                continue
            if vulnerability not in vulnerabilities:
                vulnerabilities.append(vulnerability)
                print(json.dumps(vulnerability, indent=4))

if len(vulnerabilities) > 0:
    sys.exit(1)
