module "http_api" {
  source = "../../modules/http-api"

  name_prefix = var.project_name
  environment = var.environment

  lambda_authorizer_name       = module.lambda_authorizer.function_name
  lambda_authorizer_invoke_arn = module.lambda_authorizer.invoke_arn

  lambda_functions = {
    get_user = {
      handler     = "bootstrap"
      runtime     = "provided.al2023"
      source_path = "${path.module}/../../../dist/user_service/routes/get_user.zip"
      policy_statements = [
        {
          Effect = "Allow"
          Action = [
            "dynamodb:GetItem"
          ]
          Resource = [
            aws_dynamodb_table.user_service_table.arn
          ]
        }
      ]
    },
    create_user = {
      handler     = "bootstrap"
      runtime     = "provided.al2023"
      source_path = "${path.module}/../../../dist/user_service/routes/create_user.zip"
      policy_statements = [
        {
          Effect = "Allow"
          Action = [
            "dynamodb:PutItem"
          ]
          Resource = [
            aws_dynamodb_table.user_service_table.arn
          ]
        }
      ]
    },
    delete_user = {
      handler     = "bootstrap"
      runtime     = "provided.al2023"
      source_path = "${path.module}/../../../dist/user_service/routes/delete_user.zip"
      policy_statements = [
        {
          Effect = "Allow"
          Action = [
            "dynamodb:DeleteItem"
          ]
          Resource = [
            aws_dynamodb_table.user_service_table.arn
          ]
        }
      ]
    },
    create_product = {
      handler     = "bootstrap"
      runtime     = "provided.al2023"
      source_path = "${path.module}/../../../dist/product_service/routes/create_product.zip"
    },
    get_product = {
      handler     = "bootstrap"
      runtime     = "provided.al2023"
      source_path = "${path.module}/../../../dist/product_service/routes/get_product.zip"
    },
    get_products = {
      handler     = "bootstrap"
      runtime     = "provided.al2023"
      source_path = "${path.module}/../../../dist/product_service/routes/get_products.zip"
    }
  }

  api_routes = [
    {
      route_key     = "GET /users/{userID}"
      function_name = "get_user"
    },
    {
      route_key     = "POST /users"
      function_name = "create_user"
    },
    {
      route_key     = "DELETE /users/{userID}"
      function_name = "delete_user"
    },
    {
      route_key     = "POST /products"
      function_name = "create_product"
    },
    {
      route_key     = "GET /products/{productID}"
      function_name = "get_product"
    },
    {
      route_key     = "GET /products"
      function_name = "get_products"
    }
  ]

  cors_configuration = {
    allow_methods = ["GET", "POST", "DELETE"]
    allow_origins = ["*"]
    allow_headers = [
      "Content-Type",
      "Authorization",
      "X-Requested-With",
      "Accept",
      "Origin"
    ]
    expose_headers    = []
    allow_credentials = false
    max_age           = 60
  }

  log_retention_days = 14
}

# amazonq-ignore-next-line
resource "aws_dynamodb_table" "user_service_table" {
  name           = "UserServiceTable"
  billing_mode   = "PROVISIONED"
  read_capacity  = 5
  write_capacity = 5
  hash_key       = "Pk"
  range_key      = "Sk"

  attribute {
    name = "Pk"
    type = "S"
  }

  attribute {
    name = "Sk"
    type = "S"
  }

  point_in_time_recovery {
    enabled = false
  }

  tags = {
    Name        = "UserServiceTable"
    Environment = var.environment
  }
}
