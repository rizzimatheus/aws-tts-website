module "lambda_newpost_golang" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "~> 7.0"

  create                 = var.create_golang_lambda
  function_name          = "tts-website-newpost-golang-${random_string.random.result}"
  description            = "Lambda function to handle new posts"
  handler                = "bootstrap"
  runtime                = "provided.al2023"
  architectures          = ["x86_64"]
  create_package         = false
  local_existing_package = "lambda/golang/newpost/bin/newpost.zip"

  timeout = 10
  publish = true

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

  attach_policy_statements = true
  policy_statements = {
    dynamodb = {
      effect = "Allow",
      actions = [
        "dynamodb:PutItem",
      ],
      resources = [
        module.dynamodb_table.dynamodb_table_arn,
      ]
    },
    sns = {
      effect = "Allow",
      actions = [
        "sns:Publish",
      ],
      resources = [
        module.sns_topic.topic_arn,
      ]
    }
  }
}

module "lambda_getposts_golang" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "~> 7.0"

  create                 = var.create_golang_lambda
  function_name          = "tts-website-getposts-golang-${random_string.random.result}"
  description            = "Lambda function to get posts from database"
  handler                = "bootstrap"
  runtime                = "provided.al2023"
  architectures          = ["x86_64"]
  create_package         = false
  local_existing_package = "lambda/golang/getposts/bin/getposts.zip"

  timeout = 10
  publish = true

  environment_variables = {
    DB_TABLE_NAME = module.dynamodb_table.dynamodb_table_id
  }

  allowed_triggers = {
    AllowExecutionFromAPIGateway = {
      service    = "apigateway"
      source_arn = "${module.api_gateway.apigatewayv2_api_execution_arn}/*"
    }
  }

  attach_policy_statements = true
  policy_statements = {
    dynamodb = {
      effect = "Allow",
      actions = [
        "dynamodb:Scan",
        "dynamodb:GetItem",
      ],
      resources = [
        module.dynamodb_table.dynamodb_table_arn,
      ]
    }
  }
}

module "lambda_convertaudio_golang" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "~> 7.0"

  create                 = var.create_golang_lambda
  function_name          = "tts-website-convertaudio-golang-${random_string.random.result}"
  description            = "Lambda function to convert new posts to mp3 files"
  handler                = "bootstrap"
  runtime                = "provided.al2023"
  architectures          = ["x86_64"]
  create_package         = false
  local_existing_package = "lambda/golang/convertaudio/bin/convertaudio.zip"

  timeout = 900
  publish = true

  environment_variables = {
    BUCKET_NAME = module.s3_storage.s3_bucket_id
  }

  allowed_triggers = {
    AllowExecutionFromSNS = {
      principal  = "sns.amazonaws.com"
      source_arn = module.sns_topic.topic_arn
    }
  }

  attach_policy_statements = true
  policy_statements = {
    polly = {
      effect = "Allow",
      actions = [
        "polly:SynthesizeSpeech",
      ],
      resources = [
        "*"
      ]
    },
    dynamodb = {
      effect = "Allow",
      actions = [
        "dynamodb:GetItem",
        "dynamodb:UpdateItem",
      ],
      resources = [
        module.dynamodb_table.dynamodb_table_arn,
      ]
    },
    s3 = {
      effect = "Allow",
      actions = [
        "s3:PutObject",
        "s3:PutObjectAcl",
        "s3:GetBucketLocation",
      ],
      resources = [
        module.s3_storage.s3_bucket_arn,
        "${module.s3_storage.s3_bucket_arn}/*",
      ]
    }
  }
}
