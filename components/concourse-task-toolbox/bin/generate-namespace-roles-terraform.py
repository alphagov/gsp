#!/usr/bin/env python3
import argparse
import collections
import json
import yaml

parser = argparse.ArgumentParser()
parser.add_argument("--account-id")
parser.add_argument("--config-file")
args = parser.parse_args()

with open(args.config_file) as f:
    in_data = yaml.safe_load(f)

out_data = collections.defaultdict(lambda: collections.defaultdict(dict))
out_data["variable"] = {
    "aws_account_role_arn": {
        "type": "string"
    },
    "cluster_name": {
        "type": "string"
    }
}
out_data["provider"] = {
    "aws": {
        "region": "eu-west-2",
        "assume_role": {
            "role_arn": "${var.aws_account_role_arn}"
        }
    }
}

def genResourceID(*labels, delimiter='-'):
    return delimiter.join(labels)

for namespace in in_data.get('namespaces', []):
    if len(namespace.get('roles', [])):
        out_data['module'][namespace['name']] = {
            "source": "../platform/modules/gsp-namespace-policies",
            "namespace_name": namespace['name'],
            "account_id": args.account_id,
            "cluster_name": "${var.cluster_name}"
        }
    for role in namespace.get('roles', []):
        out_data['resource']['aws_iam_role'][genResourceID(namespace['name'], role['name'])] = {
            "name": genResourceID("${var.cluster_name}", "namespace", namespace['name'], role['name']),
            "assume_role_policy": json.dumps({
                "Version": "2012-10-17",
                "Statement": {
                    "Effect": "Allow",
                    "Action": "sts:AssumeRole",
                    "Principal": {
                        "AWS": genResourceID("arn", "aws", "iam", "", args.account_id, "role/${var.cluster_name}_kiam_server", delimiter=":")
                    }
                }
            }),
            "path": "/gsp/${var.cluster_name}/namespaceroles/"
        }
        for policy in role['policies']:
            policy_name = genResourceID("${var.cluster_name}", "namespace", namespace['name'], policy)
            out_data['resource']['aws_iam_role_policy_attachment'][genResourceID(namespace['name'], role['name'], policy)] = {
                "role": genResourceID("${var.cluster_name}", "namespace", namespace['name']),
                "policy_arn": genResourceID("arn", "aws", "iam", "", args.account_id, "policy/" + policy_name, delimiter=":")
            }

print(json.dumps(out_data))
