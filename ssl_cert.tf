resource "aws_acm_certificate" "cert" {
  domain_name       = var.cert_domain_name
  validation_method = "DNS"

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_acm_certificate_validation" "cert" {
  certificate_arn = aws_acm_certificate.cert.arn
}