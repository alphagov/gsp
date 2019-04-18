variable "existing_zone" {
  description = "the FQDN of the existing root zone to delegate a subdomain from"
  type        = "string"
}

variable "delegated_zone" {
  description = "the FQDN of the new zone delegated from the existing_zone"
  type        = "string"
}
