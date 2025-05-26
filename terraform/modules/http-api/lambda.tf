module "lambda_functions" {
  for_each = var.lambda_functions

  source = "../lambda-function"

  name_prefix = var.name_prefix
  environment = var.environment

  function_name         = each.key
  handler               = each.value.handler
  runtime               = each.value.runtime
  memory_size           = each.value.memory_size
  timeout               = each.value.timeout
  source_path           = each.value.source_path
  environment_variables = each.value.environment_variables
  policy_statements     = each.value.policy_statements
}
