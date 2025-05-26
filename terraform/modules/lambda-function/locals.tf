locals {
  base_name          = replace(lower("${var.name_prefix}-${var.function_name}"), "/[\\s_]/", "-")
  lambda_name        = replace(lower("${local.base_name}-${var.environment}"), "/[\\s_]/", "-")
  role_name          = replace(lower("${local.base_name}-lambda-role-${var.environment}"), "/[\\s_]/", "-")
  log_group_name     = replace(lower("${local.base_name}-logs-${var.environment}"), "/[\\s_]/", "-")
  inline_policy_name = replace(lower("${local.base_name}-inline-policy"), "/[\\s_]/", "-")
}
