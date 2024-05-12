module "lambda_newpost_amplify" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "~> 7.0"

  function_name          = "tts-website-newpost-amplify"
  description            = "Lambda function to handle new posts from DynamoDB Streams"
  handler                = "bootstrap"
  runtime                = "provided.al2023"
  architectures          = ["x86_64"]
  create_package         = false
  local_existing_package = "lambda/golang/newpost-amplify/bin/newpost-amplify.zip"

  timeout = 10
  publish = true

  environment_variables = {
    SNS_TOPIC = module.sns_topic.topic_arn
  }

  event_source_mapping = {
    dynamodb_public = {
      event_source_arn           = data.aws_dynamodb_table.dynamodb_public_table.stream_arn
      starting_position          = "LATEST"
      destination_arn_on_failure = module.sqs_failure.queue_arn
      filter_criteria = [
        {
          pattern = jsonencode({
            eventName : ["INSERT"]
          })
        },
      ]
    }
    dynamodb_private = {
      event_source_arn           = data.aws_dynamodb_table.dynamodb_private_table.stream_arn
      starting_position          = "LATEST"
      destination_arn_on_failure = module.sqs_failure.queue_arn
      filter_criteria = [
        {
          pattern = jsonencode({
            eventName : ["INSERT"]
          })
        },
      ]
    }
  }

  allowed_triggers = {
    dynamodb_public = {
      principal  = "dynamodb.amazonaws.com"
      source_arn = data.aws_dynamodb_table.dynamodb_public_table.stream_arn
    },
    dynamodb_private = {
      principal  = "dynamodb.amazonaws.com"
      source_arn = data.aws_dynamodb_table.dynamodb_private_table.stream_arn
    },
  }

  attach_policy_statements = true
  policy_statements = {
    dynamodb = {
      effect = "Allow",
      actions = [
        "dynamodb:GetRecords",
        "dynamodb:GetShardIterator",
        "dynamodb:DescribeStream",
        "dynamodb:ListStreams",
      ],
      resources = [
        data.aws_dynamodb_table.dynamodb_public_table.stream_arn,
        data.aws_dynamodb_table.dynamodb_private_table.stream_arn,
      ]
    },
    sns = {
      effect    = "Allow",
      actions   = ["sns:Publish"],
      resources = [module.sns_topic.topic_arn]
    },
    sqs_failure = {
      effect    = "Allow",
      actions   = ["sqs:SendMessage"],
      resources = [module.sqs_failure.queue_arn]
    },
  }

  depends_on = [
    module.dynamodb_table.dynamodb_table_stream_arn
  ]
}

module "sqs_failure" {
  source  = "terraform-aws-modules/sqs/aws"
  version = "~> 4.0"

  name = "tts-website-newpost-amplify-failure"
}

data "aws_dynamodb_table" "dynamodb_public_table" {
  name = "PublicPost-24or6mw5f5ezrerhjtupxwdyzi-NONE"
}

data "aws_dynamodb_table" "dynamodb_private_table" {
  name = "PrivatePost-24or6mw5f5ezrerhjtupxwdyzi-NONE"
}
