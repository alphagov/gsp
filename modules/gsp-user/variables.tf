variable "role_prefix" {
  description = "prefix string given to role"
  default     = "user"
}

variable "user_name" {
  description = "unique name for the user"
}

variable "user_arn" {
  description = "IAM user arn that will assume this role"
}

variable "cluster_name" {
  description = "cluster name to scope this role to"
}

variable "source_cidrs" {
  description = "Source CIDRs that are allowed to perform the assume role"
  type        = list(string)

  default = [
    "213.86.153.212/32",
    "213.86.153.213/32",
    "213.86.153.214/32",
    "213.86.153.235/32",
    "213.86.153.236/32",
    "213.86.153.237/32",
    "85.133.67.244/32",
  ]
}

