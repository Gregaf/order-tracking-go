resource "aws_iam_role" "lambda_role" {
  name = local.role_name

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
      Name        = local.role_name
      Environment = var.environment
    }
  )
}

resource "aws_iam_role_policy" "lambda_custom_inline_policies" {
  count = length(var.policy_statements) > 0 ? 1 : 0

  name = local.inline_policy_name
  role = aws_iam_role.lambda_role.id
  policy = jsonencode({
    Version   = "2012-10-17"
    Statement = var.policy_statements
  })
}

resource "aws_iam_role_policy_attachment" "lambda_basic_execution" {
  role       = aws_iam_role.lambda_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# resource "aws_iam_role_policy_attachment" "lambda_xray" {
#   for_each = {
#     for key, value in var.lambda_functions : key => value
#     if var.enable_xray
#   }

#   role       = aws_iam_role.lambda_roles[each.key].arn
#   policy_arn = "arn:aws:iam::aws:policy/AWSXrayWriteOnlyAccess"
# }

resource "aws_lambda_function" "this" {
  function_name = local.lambda_name
  role          = aws_iam_role.lambda_role.arn
  handler       = var.handler
  runtime       = var.runtime
  memory_size   = var.memory_size
  timeout       = var.timeout

  filename         = var.source_path
  source_code_hash = filebase64sha256(var.source_path)

  #   tracing_config {
  #     mode = var.enable_xray ? "Active" : "PassThrough"
  #   }

  dynamic "environment" {
    for_each = length(var.environment_variables) > 0 ? [1] : []
    content {
      variables = var.environment_variables
    }
  }

  tags = merge(
    var.tags,
    {
      Name        = local.lambda_name
      Environment = var.environment
    }
  )
}

resource "aws_cloudwatch_log_group" "lambda_logs" {
  name              = local.log_group_name
  retention_in_days = var.log_retention_days

  tags = merge(
    var.tags,
    {
      Name        = local.log_group_name
      Environment = var.environment
    }
  )
}
