provider "aws" {
  region = var.region
  access_key = var.access_key
  s3_force_path_style = var.s3_force_path_style
  secret_key = var.secret_key
  skip_credentials_validation = var.skip_credentials_validation
  skip_metadata_api_check = var.skip_metadata_api_check
  skip_requesting_account_id = var.skip_requesting_account_id

  endpoints {
    iam = var.localstack_url
    sns = var.localstack_url
  }
}

variable "project" {
  default = "products"
}

resource "aws_iam_user" "products" {
  name = "products"
  tags = {
    project = var.project
    env = var.account
  }
}

resource "aws_iam_policy" "products_sns_policy" {
  name = "products-sns-policy"

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "sns",
            "Effect": "Allow",
            "Action": [
                "sns:Publish"
            ],
            "Resource": [
                "${aws_sns_topic.products_events.arn}"
            ]
        }
    ]
}
EOF
}

resource "aws_iam_policy_attachment" "sns_policy_attach" {
  name = "products-sns-attachment"
  users = [
    aws_iam_user.products.name]
  policy_arn = aws_iam_policy.products_sns_policy.arn
}