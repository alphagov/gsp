# gsp-terraform-ignition
Terraform for the various parts of the GSP platform using Ignition.

## Intention

Use Container Linux Ignition to terraform a kubernetes cluster, including both the infrastructure and basic software components. The bootstrap phase of the initialisation of a cluster will be kept separate, using state imported from the main cluster's terraform state file. So the typical workflow will be:

1. `terraform apply` a cluster's infrastructure and basic components
1. `terraform apply` a bootstrapper and add it to the cluster
1. after the kubernetes control plane is operational, `terraform destroy` the bootstrap components

## Usage

The main entrypoint of this repo is the `gsp-cluster` terraform module. An example might look like:

```
module "gsp-cluster" {
    source = "git::https://github.com/alphagov/gsp-terraform-ignition//modules/gsp-cluster"
    cluster_name = "production"
    dns_zone = "govuk.aws.ext.govsvc.uk"
    user_data_bucket_name = "gds-govuk-production-tfstate"
    user_data_bucket_region = "eu-west-2"
    k8s_tag = "v1.12.2"
    admin_role_arns = ["arn:aws:iam::111111111111:role/admin"]

    ...
}

```

In order to use this cluster with the bootstrapper several values will need to be output by the root into the state file (see `bootstrapper/main.tf`).
