variable "name_prefix" {
  description = "Prefix to use for resource naming"
  type        = string
}

variable "environment" {
  description = "Environment name (e.g., dev, staging, prod)"
  type        = string
}

variable "tags" {
  description = "Additional tags to apply to all resources"
  type        = map(string)
  default     = {}
}

variable "lambda_authorizer_name" {
  description = ""
  type        = string
}

variable "lambda_authorizer_invoke_arn" {
  description = ""
  type        = string
}

variable "lambda_functions" {
  description = "Map of Lambda functions to create"
  type = map(object({
    handler               = string
    runtime               = string
    source_path           = string
    memory_size           = optional(number, 128)
    timeout               = optional(number, 30)
    environment_variables = optional(map(string), {})
    layers                = optional(list(string), [])
    policy_statements = optional(list(object({
      Effect   = string
      Action   = list(string)
      Resource = list(string)
    })), [])
    managed_policy_arns = optional(list(string), [])
  }))
}

variable "api_routes" {
  description = "API Gateway routes configuration"
  type = list(object({
    route_key     = string # HTTP method and path, e.g., "GET /users"
    function_name = string
  }))
}

variable "log_retention_days" {
  description = "Cloudwatch logs retention period in days"
  type        = number
  default     = 30
}

variable "enable_xray" {
  description = "Enable AWS X-Ray tracing"
  type        = bool
  default     = false
}

variable "stage_name" {
  description = "Name of the API Gateway stage"
  type        = string
  default     = "$default"
}

variable "cors_configuration" {
  description = "CORS configuration for the API Gateway"
  type = object({
    allow_origins     = list(string)
    allow_methods     = list(string)
    allow_headers     = list(string)
    expose_headers    = list(string)
    allow_credentials = bool
    max_age           = number
  })
  default = null
}
