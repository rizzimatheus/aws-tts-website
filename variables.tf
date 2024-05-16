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

variable "create_golang_lambda" {
  type        = bool
  description = "Controls whether lambda functions in golang should be created"
  default     = true
}

variable "create_python_lambda" {
  type        = bool
  description = "Controls whether lambda functions in python should be created"
  default     = false
}
