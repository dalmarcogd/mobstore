terraform {
  required_providers {
    aws = {
      source = "aws"
    }
    mysql = {
      source = "terraform-providers/mysql"
      version = "~> 1.9.0"
    }
  }
  required_version = ">= 0.13"
}
