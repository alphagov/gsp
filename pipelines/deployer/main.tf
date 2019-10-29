terraform {
  backend "s3" {
  }
}

variable "aws_account_role_arn" {
  type = string
}

provider "aws" {
  region = "eu-west-2"

  assume_role {
    role_arn = var.aws_account_role_arn
  }
}

