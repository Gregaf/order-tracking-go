resource "aws_apigatewayv2_api" "http_api" {
  name          = "${var.name_prefix}-http-api-${var.environment}"
  protocol_type = "HTTP"

  dynamic "cors_configuration" {
    for_each = var.cors_configuration != null ? [var.cors_configuration] : []
    content {
      allow_origins     = cors_configuration.value.allow_origins
      allow_methods     = cors_configuration.value.allow_methods
      allow_headers     = cors_configuration.value.allow_headers
      expose_headers    = cors_configuration.value.expose_headers
      allow_credentials = cors_configuration.value.allow_credentials
      max_age           = cors_configuration.value.max_age
    }
  }

  tags = merge(
    var.tags,
    {
      Name        = "${var.name_prefix}-http-api-${var.environment}"
      Environment = var.environment
    }
  )
}

resource "aws_apigatewayv2_stage" "default" {
  api_id      = aws_apigatewayv2_api.http_api.id
  name        = var.stage_name
  auto_deploy = true

  #   access_log_settings {
  #     destination_arn = aws
  #   }
}

resource "aws_cloudwatch_log_group" "api_logs" {
  name              = "/aws/apigateway/${var.name_prefix}-http-api-${var.environment}"
  retention_in_days = var.log_retention_days

  tags = merge(
    var.tags,
    {
      Name        = "${var.name_prefix}-api-logs-${var.environment}"
      Environment = var.environment
    }
  )
}

resource "aws_apigatewayv2_integration" "lambda_integrations" {
  for_each = { for k, v in var.lambda_functions : k => v }

  api_id                 = aws_apigatewayv2_api.http_api.id
  integration_type       = "AWS_PROXY"
  integration_uri        = aws_lambda_function.functions[each.key].invoke_arn
  payload_format_version = "2.0"
  timeout_milliseconds   = 30000

  response_parameters {
    status_code = 200
    mappings = var.enable_xray ? {
      "append:header.X-Ray-Trace-Id" = "$context.xrayTraceId"
    } : {}
  }
}

resource "aws_apigatewayv2_route" "routes" {
  for_each = { for idx, route in var.api_routes : idx => route }

  api_id    = aws_apigatewayv2_api.http_api.id
  route_key = each.value.route_key
  target    = "integrations/${aws_apigatewayv2_integration.lambda_integrations[each.value.function_name].id}"
}

resource "aws_lambda_permission" "api_gateway_permissions" {
  for_each = { for idx, route in var.api_routes : idx => route }

  statement_id  = "AllowExecutionFromAPIGateway-${each.key}"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.functions[each.value.function_name].function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.http_api.execution_arn}/*/*"
}
