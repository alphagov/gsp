#!/usr/bin/env python
import argparse
import json
import yaml

parser = argparse.ArgumentParser()
parser.add_argument("--account-id")
parser.add_argument("--config-file")
args = parser.parse_args()

with open(args.config_file) as f:
    in_data = yaml.safe_load(f)

out_data = {
    "module": {},
    "variable": {
        "aws_account_role_arn": {
            "type": "string"
        },
        "cluster_name": {
            "type": "string"
        }
    },
    "provider": {
        "aws": {
            "region": "eu-west-2",
            "assume_role": {
                "role_arn": "${var.aws_account_role_arn}"
            }
        }
    }
}

for namespace in in_data.get('namespaces', []):
    if len(namespace.get('roles', [])):
        out_data['module']['-'.join([namespace['name'], 'policies'])] = {
            "source": "../platform/modules/gsp-namespace-policies",
            "namespace_name": namespace['name'],
            "account_id": args.account_id,
            "cluster_name": "${var.cluster_name}"
        }
    for role in namespace.get('roles', []):
        out_data['module']['-'.join([namespace['name'], role['name']])] = {
            "source": "../platform/modules/gsp-namespace-role",
            "namespace_name": namespace['name'],
            "role_name": role['name'],
            "account_id": args.account_id,
            "cluster_name": "${var.cluster_name}"
        }
        for policy in role['policies']:
            out_data['module']['-'.join([namespace['name'], role['name'], policy])] = {
                "source": "../platform/modules/gsp-namespace-role-policy-attachment",
                "namespace_name": namespace['name'],
                "role_name": role['name'],
                "policy_id": policy,
                "account_id": args.account_id,
                "cluster_name": "${var.cluster_name}"
            }

print(json.dumps(out_data))
