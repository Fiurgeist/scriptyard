terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.27"
    }
  }

  required_version = ">= 0.14.9"
}

provider "aws" {
  profile = "default"
  region = var.region
}

resource "aws_vpc" "lambda_test" {
  cidr_block = "10.42.0.0/16"
}

resource "aws_subnet" "lambda_public" {
  vpc_id     = aws_vpc.lambda_test.id
  cidr_block = "10.42.1.0/24"
  
  map_public_ip_on_launch = true

  tags = {
    Name = "Lambda Public"
  }
}

resource "aws_iam_role" "iam_for_lambda" {
  name = "iam_for_lambda"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "lambda_basic" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
  role       = aws_iam_role.iam_for_lambda.name
}

resource "aws_lambda_function" "lambda_ex" {
  filename      = "lambda.zip"
  function_name = "udacity_project_2"
  role          = aws_iam_role.iam_for_lambda.arn
  handler       = "lambda.lambda_handler"

  source_code_hash = filebase64sha256("lambda.zip")

  runtime = "python3.8"

  environment {
    variables = {
      greeting = "Konbanwa"
    }
  }
}

resource "aws_lambda_invocation" "lambda_ex" {
  function_name = aws_lambda_function.lambda_ex.function_name

  input = ""
}
