#!/usr/bin/env bash

AWS_ACCOUNT_ID="$(aws sts get-caller-identity | jq -r .Account)"

docker build \
	--network host \
	--build-arg AWS_INTEGRATION=true \
	--build-arg AWS_ACCESS_KEY_ID \
	--build-arg AWS_SECRET_ACCESS_KEY \
	--build-arg AWS_SESSION_TOKEN \
	--build-arg AWS_RDS_SECURITY_GROUP_ID=sg-04521d05ba3d9edb5 \
	--build-arg AWS_RDS_SUBNET_GROUP_NAME=sandbox-private \
	--build-arg AWS_PRINCIPAL_PERMISSIONS_BOUNDARY_ARN=arn:aws:iam::${AWS_ACCOUNT_ID}:policy/sandbox-service-operator-managed-role-permissions-boundary \
	--build-arg AWS_PRINCIPAL_SERVER_ROLE_ARN=arn:aws:iam::${AWS_ACCOUNT_ID}:role/sandbox_kiam_server \
	.
