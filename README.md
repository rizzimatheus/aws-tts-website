# AudioNotes V1

Welcome to AudioNotes V1, a web application where you can transform text into speech using various voices. Write your text, select a voice, and get an audio file read by your chosen voice.

## Demo
Try out the live demo [here](https://tts.matheusrizzi.com).

## Documentation
For detailed information, visit the [documentation](https://matheusrizzi.com/projects/en/audionotes).

## Technology Stack
![HTML5](https://img.shields.io/badge/html5-%23E34F26.svg?style=for-the-badge&logo=html5&logoColor=white) ![JavaScript](https://img.shields.io/badge/javascript-%23323330.svg?style=for-the-badge&logo=javascript&logoColor=%23F7DF1E) ![CSS3](https://img.shields.io/badge/css3-%231572B6.svg?style=for-the-badge&logo=css3&logoColor=white) ![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white) ![Python](https://img.shields.io/badge/python-3670A0?style=for-the-badge&logo=python&logoColor=ffdd54) ![AWS](https://img.shields.io/badge/AWS-%23FF9900.svg?style=for-the-badge&logo=amazon-aws&logoColor=white) ![Terraform](https://img.shields.io/badge/terraform-%235835CC.svg?style=for-the-badge&logo=terraform&logoColor=white)
- **Frontend**: HTML, CSS, JavaScript
- **Hosting**: AWS S3 with CloudFront distribution
- **Backend**: AWS API Gateway, DynamoDB, S3, Lambda (Golang & Python), SNS, Amazon Polly
- **IaC**: Terraform

## Installation
1. Ensure you have an AWS account and have set up AWS credentials on your machine.
2. If you prefer not to host the Terraform backend on S3, you can delete the `backend.tf` file. Alternatively, you can manually create the DynamoDB table and S3 bucket and specify them in this file. For automated resource creation, use the CloudFormation script available [here](https://github.com/rizzimatheus/terraform-backend-aws).
3. Edit `variables.tf` to update the domain names to ones you own.
4. For SSL certificate validation, change the `validation_method` in `ssl-cert.tf` to `EMAIL` or `DNS`.
5. Run `terraform init` followed by `terraform apply`.
6. Validate the SSL certificate using the chosen method.
7. Optionally, edit `variables.tf` to select between Golang and Python for Lambda functions. Update `sns.tf` and `apigw.tf` accordingly.
8. CloudFront distribution may take a few minutes to deploy. Monitor the progress in the AWS console.

## Usage
Access the application using the specified domain or the CloudFront URL provided after executing `terraform apply`.

## Deletion
To remove the deployed resources, execute `terraform destroy`.

## AudioNotes V2

For the second version of AudioNotes, please visit the [AudioNotes V2 repository](https://github.com/rizzimatheus/next-amplify-tts-website).


## License
This project is licensed under the MIT License - see the LICENSE file for details.
