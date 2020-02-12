#!/usr/bin/env bash

AWS_ACCOUNT_ID="$(aws sts get-caller-identity | jq -r .Account)"
AWS_RDS_SECURITY_GROUP_ID=$(aws ec2 describe-security-groups | jq -r '.SecurityGroups[] | select(.GroupName == "sandbox_rds_from_worker") | .GroupId')

docker build \
	--network host \
	--build-arg AWS_INTEGRATION=true \
	--build-arg AWS_ACCESS_KEY_ID \
	--build-arg AWS_SECRET_ACCESS_KEY \
	--build-arg AWS_SESSION_TOKEN \
	--build-arg AWS_RDS_SECURITY_GROUP_ID=$AWS_RDS_SECURITY_GROUP_ID \
	--build-arg AWS_RDS_SUBNET_GROUP_NAME=sandbox-private \
	--build-arg AWS_PRINCIPAL_PERMISSIONS_BOUNDARY_ARN=arn:aws:iam::${AWS_ACCOUNT_ID}:policy/sandbox-service-operator-managed-role-permissions-boundary \
	--build-arg AWS_PRINCIPAL_SERVER_ROLE_ARN=arn:aws:iam::${AWS_ACCOUNT_ID}:role/sandbox_kiam_server \
	--build-arg AWS_ROLE_ARN=arn:aws:iam::${AWS_ACCOUNT_ID}:role/admin \
	--build-arg AWS_OIDC_PROVIDER_ARN=arn:aws:iam::${AWS_ACCOUNT_ID}:oidc-provider/oidc.eks.eu-west-2.amazonaws.com/id/D4AF693862F6BE27DFD2FCA407D8990D \
	--build-arg AWS_OIDC_PROVIDER_URL=oidc.eks.eu-west-2.amazonaws.com/id/D4AF693862F6BE27DFD2FCA407D8990D \
	.
