#!/usr/bin/env bash

set -euo pipefail

aws_vault_profile="${1}"

aws-vault exec "${aws_vault_profile}" -- terraform init --upgrade=true
aws-vault exec "${aws_vault_profile}" -- terraform apply -auto-approve

echo "Waiting for kubernetes..."

until [ -n "$(aws-vault exec "${aws_vault_profile}" -- kubectl --kubeconfig kubeconfig -n kube-system get pods | grep kube-apiserver | grep -v bootstrap)" ]
do
    echo -n "."
    sleep 10s
done

aws-vault exec "${aws_vault_profile}" -- terraform destroy -auto-approve
