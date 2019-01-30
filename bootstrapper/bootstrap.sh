#!/usr/bin/env bash

set -euo pipefail

aws_vault_profile="${1}"

read -r -p "Remote state bucket name: " remote_state_bucket_name
read -r -p "Remote state key: " remote_state_key

aws-vault exec "${aws_vault_profile}" -- terraform init --upgrade=true
aws-vault exec "${aws_vault_profile}" -- terraform apply -auto-approve -var "remote_state_bucket_name=${remote_state_bucket_name}" -var "remote_state_key=${remote_state_key}"

echo "Waiting for kubernetes..."

until aws-vault exec "${aws_vault_profile}" -- kubectl --kubeconfig kubeconfig -n kube-system get pods | grep kube-apiserver | grep -v bootstrap
do
    echo -n "."
    sleep 10s
done

aws-vault exec "${aws_vault_profile}" -- kubectl --kubeconfig kubeconfig drain --ignore-daemonsets --force --grace-period=60 -l 'node-role.kubernetes.io/bootstrapper'

# the terraform destroy will delete the kubeconfig file but it is needed to delete the bootstrapper node
cp kubeconfig user-kubeconfig
aws-vault exec "${aws_vault_profile}" -- terraform destroy -auto-approve -var "remote_state_bucket_name=${remote_state_bucket_name}" -var "remote_state_key=${remote_state_key}"
aws-vault exec "${aws_vault_profile}" -- kubectl --kubeconfig user-kubeconfig delete node -l 'node-role.kubernetes.io/bootstrapper'
rm user-kubeconfig
