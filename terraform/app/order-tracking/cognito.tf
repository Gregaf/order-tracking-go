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

  schema {
    attribute_data_type = "String"
    name                = "name"
    mutable             = true
    required            = true
    string_attribute_constraints {
      min_length = 1
      max_length = 64
    }
  }

  schema {
    attribute_data_type = "String"
    name                = "email"
    mutable             = true
    required            = true
    string_attribute_constraints {
      min_length = 1
      max_length = 256
    }
  }

  schema {
    attribute_data_type = "String"
    name                = "family_name"
    mutable             = true
    required            = true
    string_attribute_constraints {
      min_length = 1
      max_length = 64
    }
  }

  lambda_config {
    pre_token_generation_config {
      lambda_arn     = module.pre_token_lambda.arn
      lambda_version = "V2_0"
    }
    post_authentication = module.post_auth_lambda.arn
  }
}

resource "aws_cognito_user_pool_client" "this" {
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
  allowed_oauth_flows_user_pool_client = true
  allowed_oauth_flows                  = ["code"]
  allowed_oauth_scopes                 = ["phone", "email", "openid", "profile"]
  callback_urls                        = ["http://localhost:5173/callback"]
  logout_urls                          = ["http://localhost:5173/logout"]
  supported_identity_providers         = ["COGNITO"]

  explicit_auth_flows = [
    "ALLOW_USER_SRP_AUTH",
    "ALLOW_REFRESH_TOKEN_AUTH",
    "ALLOW_USER_PASSWORD_AUTH"
  ]
}

resource "aws_cognito_identity_provider" "google_provider" {
  user_pool_id  = aws_cognito_user_pool.this.id
  provider_name = "Google"
  provider_type = "Google"

  provider_details = {
    authorize_scopes = "email"
    client_id        = var.google_client_id
    client_secret    = var.google_client_secret
  }

  attribute_mapping = {
    email       = "email"
    family_name = "family_name"
    name        = "name"
  }
}

resource "aws_cognito_user_pool_domain" "this" {
  domain       = "${var.project_name}-${var.environment}"
  user_pool_id = aws_cognito_user_pool.this.id

  managed_login_version = 2
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

resource "aws_lambda_permission" "allow_cognito_invoke_post_auth" {
  statement_id  = "AllowCognitoInvoke"
  action        = "lambda:InvokeFunction"
  function_name = module.post_auth_lambda.function_name
  principal     = "cognito-idp.amazonaws.com"
  source_arn    = aws_cognito_user_pool.this.arn
}

module "post_auth_lambda" {
  source = "../../modules/lambda-function"

  name_prefix = var.project_name
  environment = var.environment

  function_name = "post_auth"
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  source_path   = "${path.module}/../../../dist/auth_service/flow/post_auth.zip"
  policy_statements = [
    {
      Effect = "Allow"
      Action = ["dynamodb:PutItem"]
      Resource = [
        aws_dynamodb_table.user_service_table.arn
      ]
    }
  ]
}
