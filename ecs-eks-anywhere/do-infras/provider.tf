terraform {
  required_providers {
    digitalocean = {
      source = "digitalocean/digitalocean"
      version = "1.22.2"
    }
  }
  backend "s3" {
    bucket = "ecs-anywhere-terraform-state"
    key = "global/s3/terraform.tfstate"
    region = "us-east-1"
    encrypt = true
  }
}

resource "aws_s3_bucket" "terraform_state" {
  bucket = "ecs-anywhere-terraform-state"
  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        sse_algorithm = "AES256"
      }
    }
  }
}

variable "do_token" {}
variable "pvt_key" {}

provider "digitalocean" {
  token = var.do_token
}

provider "aws" {
  region = "us-east-1"
}