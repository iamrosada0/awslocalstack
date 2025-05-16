terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}
variable "access_key" {
  type    = string
  default = "dummy"
}
variable "secret_key" {
  type    = string
  default = "dummy"
}
variable "region" {
  type    = string
  default = "us-west-2"
}
variable "localstack_endpoint" {
  type    = string
  default = "http://localhost:4566"
}
variable "bucket_name" {
  type    = string
  default = "my-test-bucket"
}

variable "sqs_queue_name" {
  type    = string
  default = "my-custom-sqs-queue"
}
provider "aws" {
  access_key                  = var.access_key
  secret_key                  = var.secret_key
  region                      = var.region

  s3_use_path_style           = true
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true

  endpoints {
    s3 = var.localstack_endpoint
    sqs= var.localstack_endpoint
  }
}

resource "aws_s3_bucket" "test-bucket" {
  bucket = var.bucket_name
}

resource "aws_sqs_queue" "terraform_queue" {
  name                      = var.sqs_queue_name
  delay_seconds             = 90
  max_message_size          = 2048
  message_retention_seconds = 86400
  receive_wait_time_seconds = 10
 
}