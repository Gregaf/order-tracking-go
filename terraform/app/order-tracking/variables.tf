variable "aws_region" {
  type    = string
  default = "us-east-1"
}

variable "project_name" {
  type = string
}

variable "environment" {
  description = "Environment"
  type        = string

  validation {
    condition     = contains(["dev", "staging", "prod"], var.environment)
    error_message = "Environment must be one of: dev, staging, prod"
  }
}

variable "google_client_id" {
  description = "Google Client ID for OAuth"
  type        = string
}

variable "google_client_secret" {
  description = "Google Client Secret for OAuth"
  type        = string
  sensitive   = true
}
