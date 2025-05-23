variable "project_name" {
  description = "Name of the project for storing Terraform remote statefile."
  type        = string
}

variable "environment" {
  description = "Environment"
  type        = string

  validation {
    condition     = contains(["dev", "staging", "prod"], var.environment)
    error_message = "Environment must be one of: dev, staging, prod"
  }
}

variable "tags" {
  description = "A map of tags to assign to resources"
  type        = map(string)
  default     = {}
}
