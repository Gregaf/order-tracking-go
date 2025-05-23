resource "aws_iam_role" "lambda_roles" {
  for_each = var.lambda_functions

  name = "${var.name_prefix}-${each.key}-lambda-role-${var.environment}"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })

  tags = merge(
    var.tags,
    {
      Name        = "${var.name_prefix}-${each.key}-lambda-role-${var.environment}"
      Environment = var.environment
    }
  )
}

resource "aws_iam_role_policy" "lambda_custom_inline_policies" {
  for_each = {
    for key, lambda in var.lambda_functions : key => lambda
    if length(lookup(lambda, "policy_statements", [])) > 0
  }

  name = "${var.name_prefix}-${each.key}-inline-policy"
  role = aws_iam_role.lambda_roles[each.key].id
  policy = jsonencode({
    Version   = "2012-10-17"
    Statement = each.value.policy_statements
  })
}

resource "aws_iam_role_policy_attachment" "lambda_basic_execution" {
  for_each = var.lambda_functions

  role       = aws_iam_role.lambda_roles[each.key].name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy_attachment" "lambda_xray" {
  for_each = {
    for key, value in var.lambda_functions : key => value
    if var.enable_xray
  }

  role       = aws_iam_role.lambda_roles[each.key].arn
  policy_arn = "arn:aws:iam::aws:policy/AWSXrayWriteOnlyAccess"
}

resource "aws_lambda_function" "functions" {
  for_each = var.lambda_functions

  function_name = "${var.name_prefix}-${each.key}-${var.environment}"
  role          = aws_iam_role.lambda_roles[each.key].arn
  handler       = each.value.handler
  runtime       = each.value.runtime
  memory_size   = each.value.memory_size
  timeout       = each.value.timeout
  layers        = each.value.layers

  filename         = each.value.source_path
  source_code_hash = filebase64sha256(each.value.source_path)

  tracing_config {
    mode = var.enable_xray ? "Active" : "PassThrough"
  }

  dynamic "environment" {
    for_each = length(each.value.environment_variables) > 0 ? [1] : []
    content {
      variables = each.value.environment_variables
    }
  }

  tags = merge(
    var.tags,
    {
      Name        = "${var.name_prefix}-${each.key}-${var.environment}"
      Environment = var.environment
    }
  )
}

resource "aws_cloudwatch_log_group" "lambda_logs" {
  for_each = var.lambda_functions

  name              = "${var.name_prefix}-${each.key}-logs-${var.environment}"
  retention_in_days = var.log_retention_days

  tags = merge(
    var.tags,
    {
      Name        = "${var.name_prefix}-${each.key}-${var.environment}"
      Environment = var.environment
    }
  )
}
