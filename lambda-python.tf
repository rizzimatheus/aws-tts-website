# module "lambda_newpost" {
#   source  = "terraform-aws-modules/lambda/aws"
#   version = "~> 7.0"

#   function_name = "TTSWebsite_NewPost"
#   description   = "Lambda function to handle new posts"
#   handler       = "newpost.lambda_handler"
#   runtime       = "python3.12"
#   source_path   = "lambda/python/newpost"
#   timeout       = 10
#   publish       = true

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

#   attach_policy_statements = true
#   policy_statements = {
#     dynamodb = {
#       effect = "Allow",
#       actions = [
#         "dynamodb:PutItem",
#       ],
#       resources = [
#         module.dynamodb_table.dynamodb_table_arn,
#       ]
#     },
#     sns = {
#       effect = "Allow",
#       actions = [
#         "sns:Publish",
#       ],
#       resources = [
#         module.sns_topic.topic_arn,
#       ]
#     }
#   }
# }

# module "lambda_convertaudio" {
#   source  = "terraform-aws-modules/lambda/aws"
#   version = "~> 7.0"

#   function_name = "TTSWebsite_ConvertAudio"
#   description   = "Lambda function to convert new posts to mp3 files"
#   handler       = "convertaudio.lambda_handler"
#   runtime       = "python3.12"
#   source_path   = "lambda/python/convertaudio"
#   timeout       = 900
#   publish       = true

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

#   attach_policy_statements = true
#   policy_statements = {
#     polly = {
#       effect = "Allow",
#       actions = [
#         "polly:SynthesizeSpeech",
#       ],
#       resources = [
#         "*"
#       ]
#     },
#     dynamodb = {
#       effect = "Allow",
#       actions = [
#         "dynamodb:GetItem",
#         "dynamodb:UpdateItem",
#       ],
#       resources = [
#         module.dynamodb_table.dynamodb_table_arn,
#       ]
#     },
#     s3 = {
#       effect = "Allow",
#       actions = [
#         "s3:PutObject",
#         "s3:PutObjectAcl",
#         "s3:GetBucketLocation",
#       ],
#       resources = [
#         module.s3_storage.s3_bucket_arn,
#         "${module.s3_storage.s3_bucket_arn}/*",
#       ]
#     }
#   }
# }

# module "lambda_getposts" {
#   source  = "terraform-aws-modules/lambda/aws"
#   version = "~> 7.0"

#   function_name = "TTSWebsite_GetPosts"
#   description   = "Lambda function to get posts from database"
#   handler       = "getposts.lambda_handler"
#   runtime       = "python3.12"
#   source_path   = "lambda/python/getposts"
#   timeout       = 10
#   publish       = true

#   environment_variables = {
#     DB_TABLE_NAME = module.dynamodb_table.dynamodb_table_id
#   }

#   allowed_triggers = {
#     AllowExecutionFromAPIGateway = {
#       service    = "apigateway"
#       source_arn = "${module.api_gateway.apigatewayv2_api_execution_arn}/*"
#     }
#   }

#   attach_policy_statements = true
#   policy_statements = {
#     dynamodb = {
#       effect = "Allow",
#       actions = [
#         "dynamodb:Scan",
#         "dynamodb:GetItem",
#       ],
#       resources = [
#         module.dynamodb_table.dynamodb_table_arn,
#       ]
#     }
#   }
# }
