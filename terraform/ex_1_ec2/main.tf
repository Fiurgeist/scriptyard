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
  region = "us-east-1"
}

resource "aws_instance" "test_t2" {
  count = "4"

  ami           = "ami-0a8b4cd432b1c3063"
  instance_type = "t2.micro"

  tags = {
    Name = "Test T2 Terraform"
  }
}

resource "aws_instance" "test_m4" {
  count = "2"

  ami           = "ami-0a8b4cd432b1c3063"
  instance_type = "m4.large"

  tags = {
    Name = "Test M4 Terraform"
  }
}

