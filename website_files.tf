resource "aws_s3_object" "index_html" {
  bucket       = module.s3_website.s3_bucket_id
  key          = "index.html"
  source       = "html/index.html"
  content_type = "text/html"
  etag         = filemd5("html/index.html")
}

resource "aws_s3_object" "error_html" {
  bucket       = module.s3_website.s3_bucket_id
  key          = "error.html"
  source       = "html/error.html"
  content_type = "text/html"
  etag         = filemd5("html/error.html")
}

resource "aws_s3_object" "styles_css" {
  bucket       = module.s3_website.s3_bucket_id
  key          = "styles.css"
  source       = "html/styles.css"
  content_type = "text/css"
  etag         = filemd5("html/styles.css")
}

resource "aws_s3_object" "scripts_js" {
  bucket       = module.s3_website.s3_bucket_id
  key          = "scripts.js"
  source       = "html/scripts.js"
  content_type = "text/javascript"
  etag         = filemd5("html/scripts.js")
}
