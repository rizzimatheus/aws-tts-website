module "cdn" {
  source  = "terraform-aws-modules/cloudfront/aws"
  version = "~> 3.0"

  aliases = [var.domain_name]

  enabled             = true
  wait_for_deployment = false

  origin = {
    s3_origin = {
      domain_name = module.s3_website.s3_bucket_bucket_regional_domain_name
    }
  }

  default_root_object = "index.html"

  default_cache_behavior = {
    target_origin_id       = "s3_origin"
    viewer_protocol_policy = "allow-all"

    allowed_methods = ["GET", "HEAD", "OPTIONS"]
    cached_methods  = ["GET", "HEAD"]
    compress        = true
    query_string    = true
  }

  ordered_cache_behavior = [
    {
      path_pattern           = "*"
      target_origin_id       = "s3_origin"
      viewer_protocol_policy = "redirect-to-https"

      allowed_methods = ["GET", "HEAD", "OPTIONS"]
      cached_methods  = ["GET", "HEAD"]
      compress        = true
      query_string    = true
    }
  ]

  viewer_certificate = {
    acm_certificate_arn = aws_acm_certificate.cert.arn
    ssl_support_method  = "sni-only"
  }

  # If you're not using Rout53, you have to manually validate the domain, 
  # use the depends_on to explicit wait for the validation
  depends_on = [
    aws_acm_certificate_validation.cert
  ]
}

output "cdn_domain_name" {
  value = module.cdn.cloudfront_distribution_domain_name
}
