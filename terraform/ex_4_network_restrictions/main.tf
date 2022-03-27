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

resource "aws_vpc" "ex_restrictions" {
  cidr_block = "10.42.0.0/16"

  tags = {
    Name = "ex_restrictions"
  }
}

resource "aws_internet_gateway" "ex_restrictions_igw" {
  vpc_id = aws_vpc.ex_restrictions.id

  tags = {
    Name = "ex_restrictions_igw"
  }
}

resource "aws_route" "default_route_with_igw" {
  route_table_id         = aws_vpc.ex_restrictions.default_route_table_id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.ex_restrictions_igw.id
}

resource "aws_subnet" "ex_restrictions_public" {
  vpc_id     = aws_vpc.ex_restrictions.id
  cidr_block = "10.42.1.0/24"
  
  map_public_ip_on_launch = true


  tags = {
    Name = "ex_restrictions_public"
  }
}

resource "aws_network_acl" "ex_restrictions_acl" {
  vpc_id     = aws_vpc.ex_restrictions.id
  subnet_ids = [aws_subnet.ex_restrictions_public.id]

  ingress {
    rule_no = 100
    from_port = 0
    to_port = 0
    protocol = "tcp"
    cidr_block = "0.0.0.0/0"
    action = "allow"
  }

  ingress {
    rule_no = 101
    from_port = 0
    to_port = 0
    icmp_type = 8 # echo request type
    protocol = "icmp"
    cidr_block = "0.0.0.0/0"
    action = "allow"
  }

  egress {
    rule_no = 100
    from_port = 0
    to_port = 0
    protocol = "tcp"
    cidr_block = "0.0.0.0/0"
    action = "allow"
  }

  egress {
    rule_no = 101
    from_port = 0
    to_port = 0
    protocol = "icmp"
    cidr_block = "0.0.0.0/0"
    action = "allow"
  }

  tags = {
    Name = "ex_restrictions_acl"
  }
}

resource "aws_security_group" "ex_restrictions_sg" {
  name        = "ex_restrictions_sg"
  description = "Allow SSH for Devs"
  vpc_id      = aws_vpc.ex_restrictions.id

  ingress {
    from_port = 22
    to_port = 22
    protocol = "tcp"
    description = "for ssh"
    cidr_blocks = ["x.x.x.x/32"] # replace with own ip address
  }

  ingress {
    from_port = 8 # icmp type as from_port for some reason
    to_port = 0
    protocol = "icmp"
    description = "for pinging"
    cidr_blocks = ["x.x.x.x/32"] # replace with own ip address
  }

  egress {
    from_port = 0
    protocol = "-1"
    to_port = 0
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "ex_restrictions_sg"
  }
}

resource "aws_instance" "ex_restrictions" {
  count = "1"

  ami           = "ami-0a8b4cd432b1c3063"
  instance_type = "t2.micro"

  subnet_id              = aws_subnet.ex_restrictions_public.id
  vpc_security_group_ids = [aws_security_group.ex_restrictions_sg.id]

  tags = {
    Name = "ex_restrictions"
  }
}
