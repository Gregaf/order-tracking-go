resource "aws_cognito_user_pool" "this" {
  name = "${var.project_name}-${var.environment}"

  username_attributes      = ["email"]
  auto_verified_attributes = ["email"]

  password_policy {
    minimum_length    = 8
    require_lowercase = true
    require_numbers   = true
    require_symbols   = true
    require_uppercase = true
  }

  lambda_config {
    pre_token_generation_config {
      lambda_arn     = module.pre_token_lambda.arn
      lambda_version = "V2_0"
    }
  }
}

resource "aws_cognito_user_pool_client" "name" {
  name         = "${var.project_name}-${var.environment}"
  user_pool_id = aws_cognito_user_pool.this.id

  generate_secret = false

  refresh_token_validity = 30
  access_token_validity  = 1
  id_token_validity      = 1

  token_validity_units {
    refresh_token = "days"
    access_token  = "hours"
    id_token      = "hours"
  }

  allowed_oauth_flows          = ["implicit", "code"]
  allowed_oauth_scopes         = ["phone", "email", "openid", "profile"]
  callback_urls                = ["https://TODO/callback", "http://localhost:3000/callback"]
  logout_urls                  = ["https://TODO/logout", "http://localhost:3000/logout"]
  supported_identity_providers = ["COGNITO"]

  explicit_auth_flows = [
    "ALLOW_USER_SRP_AUTH",
    "ALLOW_REFRESH_TOKEN_AUTH",
    "ALLOW_USER_PASSWORD_AUTH"
  ]
}

module "pre_token_lambda" {
  source = "../../modules/lambda-function"

  name_prefix = var.project_name
  environment = var.environment

  function_name = "pre_token_generator"
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  source_path   = "${path.module}/../../../dist/auth_service/flow/pre_token.zip"
}

resource "aws_lambda_permission" "allow_cognito_invoke" {
  statement_id  = "AllowCognitoInvoke"
  action        = "lambda:InvokeFunction"
  function_name = module.pre_token_lambda.function_name
  principal     = "cognito-idp.amazonaws.com"
  source_arn    = aws_cognito_user_pool.this.arn
}
