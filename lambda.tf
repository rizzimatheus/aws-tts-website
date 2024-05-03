data "aws_iam_policy_document" "newpost_lambda_policy" {
  statement {
    actions = [
      "dynamodb:PutItem",
    ]
    resources = [
      module.dynamodb_table.dynamodb_table_arn,
    ]
  }
  statement {
    actions = [
      "sns:Publish",
    ]
    resources = [
      module.sns_topic.topic_arn,
    ]
  }
}

data "aws_iam_policy_document" "convertaudio_lambda_policy" {
  statement {
    actions = [
      "polly:SynthesizeSpeech",
    ]
    resources = [
      "*"
    ]
  }
  statement {
    actions = [
      "dynamodb:GetItem",
      "dynamodb:UpdateItem",
    ]
    resources = [
      module.dynamodb_table.dynamodb_table_arn,
    ]
  }
  statement {
    actions = [
      "s3:PutObject",
      "s3:PutObjectAcl",
      "s3:GetBucketLocation",
    ]
    resources = [
      module.s3_storage.s3_bucket_arn,
      "${module.s3_storage.s3_bucket_arn}/*",
    ]
  }
}

data "aws_iam_policy_document" "getposts_lambda_policy" {
  statement {
    actions = [
      "dynamodb:Scan",
      "dynamodb:GetItem",
      "dynamodb:UpdateItem",
    ]
    resources = [
      module.dynamodb_table.dynamodb_table_arn,
    ]
  }
}

module "lambda_newpost" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "~> 7.0"

  function_name          = "tts-website-newpost"
  description            = "Lambda function to handle new posts"
  handler                = "bootstrap"
  runtime                = "provided.al2023"
  architectures          = ["x86_64"]
  create_package         = false
  local_existing_package = "lambda/golang/newpost/bin/newpost.zip"

  timeout            = 10
  publish            = true
  attach_policy_json = true
  policy_json        = data.aws_iam_policy_document.newpost_lambda_policy.json

  environment_variables = {
    DB_TABLE_NAME = module.dynamodb_table.dynamodb_table_id
    SNS_TOPIC     = module.sns_topic.topic_arn
  }

  allowed_triggers = {
    AllowExecutionFromAPIGateway = {
      service    = "apigateway"
      source_arn = "${module.api_gateway.apigatewayv2_api_execution_arn}/*"
    }
  }
}

module "lambda_convertaudio" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "~> 7.0"

  function_name          = "tts-website-convertaudio"
  description            = "Lambda function to convert new posts to mp3 files"
  handler                = "bootstrap"
  runtime                = "provided.al2023"
  architectures          = ["x86_64"]
  create_package         = false
  local_existing_package = "lambda/golang/convertaudio/bin/convertaudio.zip"

  timeout            = 900
  publish            = true
  attach_policy_json = true
  policy_json        = data.aws_iam_policy_document.convertaudio_lambda_policy.json

  environment_variables = {
    DB_TABLE_NAME = module.dynamodb_table.dynamodb_table_id
    BUCKET_NAME   = module.s3_storage.s3_bucket_id
  }

  allowed_triggers = {
    AllowExecutionFromSNS = {
      principal  = "sns.amazonaws.com"
      source_arn = module.sns_topic.topic_arn
    }
  }
}

module "lambda_getposts" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "~> 7.0"

  function_name          = "tts-website-getposts"
  description            = "Lambda function to get posts from database"
  handler                = "bootstrap"
  runtime                = "provided.al2023"
  architectures          = ["x86_64"]
  create_package         = false
  local_existing_package = "lambda/golang/getposts/bin/getposts.zip"

  timeout            = 10
  publish            = true
  attach_policy_json = true
  policy_json        = data.aws_iam_policy_document.getposts_lambda_policy.json

  environment_variables = {
    DB_TABLE_NAME = module.dynamodb_table.dynamodb_table_id
  }

  allowed_triggers = {
    AllowExecutionFromAPIGateway = {
      service    = "apigateway"
      source_arn = "${module.api_gateway.apigatewayv2_api_execution_arn}/*"
    }
  }
}

# module "lambda_newpost" {
#   source = "terraform-aws-modules/lambda/aws"
#   version = "~> 7.0"

#   function_name      = "TTSWebsite_NewPost"
#   description        = "Lambda function to handle new posts"
#   handler            = "newpost.lambda_handler"
#   runtime            = "python3.12"
#   source_path        = "lambda/python/newpost"
#   timeout            = 10
#   publish            = true
#   attach_policy_json = true
#   policy_json        = data.aws_iam_policy_document.polly_lambda_policy.json

#   environment_variables = {
#     DB_TABLE_NAME = module.dynamodb_table.dynamodb_table_id
#     SNS_TOPIC     = module.sns_topic.topic_arn
#   }

#   allowed_triggers = {
#     AllowExecutionFromAPIGateway = {
#       service    = "apigateway"
#       source_arn = "${module.api_gateway.apigatewayv2_api_execution_arn}/*"
#     }
#   }
# }

# module "lambda_convertaudio" {
#   source  = "terraform-aws-modules/lambda/aws"
#   version = "~> 7.0"

#   function_name      = "TTSWebsite_ConvertAudio"
#   description        = "Lambda function to convert new posts to mp3 files"
#   handler            = "convertaudio.lambda_handler"
#   runtime            = "python3.12"
#   source_path        = "lambda/python/convertaudio"
#   timeout            = 900
#   publish            = true
#   attach_policy_json = true
#   policy_json        = data.aws_iam_policy_document.polly_lambda_policy.json

#   environment_variables = {
#     DB_TABLE_NAME = module.dynamodb_table.dynamodb_table_id
#     BUCKET_NAME   = module.s3_storage.s3_bucket_id
#   }

#   allowed_triggers = {
#     AllowExecutionFromSNS = {
#       principal  = "sns.amazonaws.com"
#       source_arn = module.sns_topic.topic_arn
#     }
#   }
# }

# module "lambda_getposts" {
#   source = "terraform-aws-modules/lambda/aws"
#   version = "~> 7.0"

#   function_name      = "TTSWebsite_GetPosts"
#   description        = "Lambda function to get posts from database"
#   handler            = "getposts.lambda_handler"
#   runtime            = "python3.12"
#   source_path        = "lambda/python/getposts"
#   timeout            = 10
#   publish            = true
#   attach_policy_json = true
#   policy_json        = data.aws_iam_policy_document.polly_lambda_policy.json

#   environment_variables = {
#     DB_TABLE_NAME = module.dynamodb_table.dynamodb_table_id
#   }

#   allowed_triggers = {
#     AllowExecutionFromAPIGateway = {
#       service    = "apigateway"
#       source_arn = "${module.api_gateway.apigatewayv2_api_execution_arn}/*"
#     }
#   }
# }
