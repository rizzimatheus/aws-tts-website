module "api_gateway" {
  source  = "terraform-aws-modules/apigateway-v2/aws"
  version = "~> 4.0"

  name          = "tts-website-api"
  description   = "TTS Website HTTP API Gateway"
  protocol_type = "HTTP"

  cors_configuration = {
    allow_headers = ["content-type", "x-amz-date", "authorization", "x-api-key", "x-amz-security-token", "x-amz-user-agent"]
    allow_methods = ["*"]
    allow_origins = ["*"]
  }

  # Custom domain
  domain_name                 = var.ttsapi_domain_name
  domain_name_certificate_arn = aws_acm_certificate.cert.arn

  # If you're not using Rout53, you have to manually validate the domain, 
  # use the depends_on to explicit wait for the validation
  depends_on = [
    aws_acm_certificate_validation.cert
  ]

  # Routes and integrations
  integrations = {
    "POST /" = {
      lambda_arn             = module.lambda_newpost.lambda_function_arn
      payload_format_version = "2.0"
      timeout_milliseconds   = 12000
      throttling_rate_limit  = 5
      throttling_burst_limit = 2
    }

    "GET /" = {
      lambda_arn             = module.lambda_getposts.lambda_function_arn
      payload_format_version = "2.0"
      timeout_milliseconds   = 12000
      throttling_rate_limit  = 5
      throttling_burst_limit = 2
    }
  }
}

output "api_gateway_domain_name" {
  value = module.api_gateway.apigatewayv2_domain_name_target_domain_name
}
