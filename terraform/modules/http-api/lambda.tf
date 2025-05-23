resource "aws_iam_role" "lambda_role" {
  name = "${var.name_prefix}-lambda-role-${var.environment}"

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
      Name        = "${var.name_prefix}-lambda-role-${var.environment}"
      Environment = var.environment
    }
  )
}

resource "aws_iam_role_policy_attachment" "lambda_basic_execution" {
  role       = aws_iam_role.lambda_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy_attachment" "lambda_xray" {
  count      = var.enable_xray ? 1 : 0
  role       = aws_iam_role.lambda_role.name
  policy_arn = "arn:aws:iam::aws:policy/AWSXrayWriteOnlyAccess"
}

resource "aws_lambda_function" "functions" {
  for_each = var.lambda_functions

  function_name = "${var.name_prefix}-${each.key}-${var.environment}"
  role          = aws_iam_role.lambda_role.arn
  handler       = each.value.handler
  runtime       = each.value.runtime
  memory_size   = each.value.memory_size
  timeout       = each.value.timeout
  layers        = each.value.layers

  filename         = "${path.module}/lambda_${each.key}.zip"
  source_code_hash = filebase64sha256("${path.module}/lambda_${each.key}.zip")

  tracing_config {
    mode = var.enable_xray ? "Active" : "PassThrough"
  }

  dynamic "environment" {
    for_each = length(each.value.environment_variables) > 0 ? [1] : []
    content {
      variables = each.value.environment_variables
    }
  }

  depends_on = [null_resource.lambda_zip_files]

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

  name              = "aws/lambda/${var.name_prefix}-${each.key}-logs-${var.environment}"
  retention_in_days = var.log_retention_days

  tags = merge(
    var.tags,
    {
      Name        = "${var.name_prefix}-${each.key}-${var.environment}"
      Environment = var.environment
    }
  )
}

resource "null_resource" "lambda_zip_files" {
  for_each = var.lambda_functions

  triggers = {
    source_code_hash = filemd5(each.value.source_path)
  }

  provisioner "local-exec" {
    command = "zip -j ${path.module}/lambda_${each.key}.zip ${each.value.source_path}"
  }
}
