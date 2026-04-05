# AWS Provider 설정 - LocalStack 타겟팅
provider "aws" {
  region                      = "us-east-1" 
  access_key                  = "test"      
  secret_key                  = "test"      
  skip_credentials_validation = true        
  skip_metadata_api_check     = true        
  skip_requesting_account_id  = true        
 
  endpoints {
    apigateway       = "http://localhost:4566"
    apigatewayv2     = "http://localhost:4566"
    cloudformation   = "http://localhost:4566"
    cloudwatch       = "http://localhost:4566"
    dynamodb         = "http://localhost:4566"
    ec2              = "http://localhost:4566"
    es               = "http://localhost:4566" # Elasticsearch Service
    events           = "http://localhost:4566" # EventBridge
    iam              = "http://localhost:4566"
    kinesis          = "http://localhost:4566"
    kms              = "http://localhost:4566"
    lambda           = "http://localhost:4566"
    logs             = "http://localhost:4566" # CloudWatch Logs
    opensearch       = "http://localhost:4566" # OpenSearch Service
    redshift         = "http://localhost:4566"
    resourcegroups   = "http://localhost:4566"
    route53          = "http://localhost:4566"
    s3               = "http://localhost:4566"
    secretsmanager   = "http://localhost:4566"
    ses              = "http://localhost:4566"
    sns              = "http://localhost:4566"
    sqs              = "http://localhost:4566"
    ssm              = "http://localhost:4566"
    stepfunctions    = "http://localhost:4566"
    sts              = "http://localhost:4566"
    # 필요에 따라 다른 서비스 추가
  }

  # S3 경로 스타일 접근 강제 (LocalStack에 필요)
  s3_use_path_style           = true
}

resource "aws_iam_role" "lambda_exec_role" {
  name = "lambda_exec_role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Action = "sts:AssumeRole",
      Effect = "Allow",
      Principal = {
        Service = "lambda.amazonaws.com"
      }
    }]
  })
}

resource "aws_lambda_function" "my_lambda" {
  function_name = "my-local-lambda"
  handler       = "handler.handler"
  runtime       = "python3.9"
  role          = aws_iam_role.lambda_exec_role.arn
  filename      = "${path.module}/lambda/function.zip"
  source_code_hash = filebase64sha256("${path.module}/lambda/function.zip")
}
