variable "domain_name" {
  type        = string
  description = "CloudFront Domain Name"
  default     = "tts.matheusrizzi.com"
}

variable "ttsapi_domain_name" {
  type        = string
  description = "API Gateway Domain Name"
  default     = "ttsapi.matheusrizzi.com"
}

variable "cert_domain_name" {
  type        = string
  description = "The wildcard domain name"
  default     = "*.matheusrizzi.com"
}
