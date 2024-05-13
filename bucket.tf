resource "random_string" "s3_random" {
  length  = 16
  special = false
  upper   = false
  numeric = true
}

module "s3_storage" {
  source  = "terraform-aws-modules/s3-bucket/aws"
  version = "~> 4.0"

  bucket                  = "tts-website-storage-${random_string.s3_random.result}"
  acl                     = "public-read"
  block_public_acls       = false
  ignore_public_acls      = false
  restrict_public_buckets = false

  control_object_ownership = true
  object_ownership         = "ObjectWriter"

  # Allow deletion of non-empty bucket
  force_destroy = true

  versioning = {
    enabled = false
  }
}

module "s3_website" {
  source = "terraform-aws-modules/s3-bucket/aws"
  bucket = var.domain_name

  # Allow deletion of non-empty bucket
  force_destroy = true

  versioning = {
    enabled = false
  }

  website = {
    index_document = "index.html"
    error_document = "error.html"
  }

  attach_policy           = true
  policy                  = data.aws_iam_policy_document.website_bucket_policy.json
  block_public_policy     = false
  block_public_acls       = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}

data "aws_iam_policy_document" "website_bucket_policy" {
  statement {
    principals {
      type        = "*"
      identifiers = ["*"]
    }

    actions = [
      "s3:GetObject",
    ]

    resources = [
      "${module.s3_website.s3_bucket_arn}/*",
    ]
  }
}
