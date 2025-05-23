terraform {
  required_version = "1.12.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.97.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

module "backend" {
  source = "../../modules/backend"

  project_name = var.project_name
  environment  = var.environment
}
