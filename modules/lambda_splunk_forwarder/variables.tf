variable "enabled" {
  default = 1
}

variable "name" {
  description = "A unique (within the cluster) name"
  type = "string"
}

variable "cloudwatch_log_group_arn" {
  description = "The ARN of the cloudwatch log group to ship to Splunk"
  type        = "string"
}

variable "cloudwatch_log_group_name" {
  description = "The name of the cloudwatch log group to ship to Splunk"
  type        = "string"
}

variable "cluster_name" {
  type = "string"
}

variable "splunk_hec_token" {
  description = "Splunk HTTP event collector token for authentication"
  type        = "string"
}

variable "splunk_hec_url" {
  description = "Splunk HTTP event collector URL to send logs to"
  type        = "string"
}

variable "splunk_index" {
  description = "Name of index to be added as metadata to logs for use in splunk"
  type        = "string"
}
