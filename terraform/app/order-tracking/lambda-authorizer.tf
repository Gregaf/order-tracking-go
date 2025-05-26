module "lambda_authorizer" {
  source = "../../modules/lambda-function"

  name_prefix = var.project_name
  environment = var.environment

  function_name = "fake-authorizer"
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  source_path   = "${path.module}/../../../dist/auth_service/authorizers/fake.zip"
}
