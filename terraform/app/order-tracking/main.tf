terraform {
  required_version = "1.12.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.97.0"
    }
  }

  backend "s3" {
    bucket       = "order-tracking-tfstate-dev"
    key          = "app/terraform.tfstate"
    region       = "us-east-1"
    use_lockfile = true
  }
}

provider "aws" {
  region = var.aws_region
}

module "hello_world_api" {
  source = "../../modules/http-api"

  name_prefix = var.project_name
  environment = var.environment

  lambda_functions = {
    hello_world = {
      handler     = "bootstrap"
      runtime     = "provided.al2023"
      source_path = "${path.module}/../../../dist/hello-world.zip"
    }
  }

  api_routes = [
    {
      route_key     = "GET /hello"
      function_name = "hello_world"
    }
  ]

  log_retention_days = 14
}

output "api_endpoint" {
  description = "The HTTP API Gateway endpoint URL"
  value       = module.hello_world_api.api_endpoint
}
