output "api_endpoint" {
  description = "The HTTP API Gateway endpoint URL"
  value       = aws_apigatewayv2_stage.default.invoke_url
}

output "api_id" {
  description = "The ID of the HTTP API Gateway"
  value       = aws_apigatewayv2_api.http_api.id
}

output "lambda_function_arns" {
  description = "Map of created Lambda function arns"
  value       = { for k, v in module.lambda_functions : k => v.arn }
}

output "api_execution_arn" {
  description = "The execution ARN of the HTTP API Gateway"
  value       = aws_apigatewayv2_api.http_api.execution_arn
}
