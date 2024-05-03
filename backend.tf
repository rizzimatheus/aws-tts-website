terraform {
  backend "s3" {
    dynamodb_table = "terraform-backend-TerraformBackendDynamoDBTable-160ZTFCXW5MLD"
    bucket         = "terraform-backend-terraformbackends3bucket-kkfee8oxkcas"
    region         = "us-east-1"
    key            = "aws-tts-website"
  }
}