# provider.tf

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0" # AWS Provider 버전 명시 (필요에 따라 조정)
    }
  }
}

# AWS Provider 설정 - LocalStack 타겟팅
provider "aws" {
  region                      = "us-east-1" # LocalStack과 호환성이 좋은 리전 사용
  access_key                  = "test"      # LocalStack은 기본적으로 아무 값이나 허용
  secret_key                  = "test"      # 실제 AWS 키를 사용하지 않도록 주의!
  skip_credentials_validation = true        # 자격증명 유효성 검사 건너뛰기
  skip_metadata_api_check     = true        # EC2 메타데이터 API 체크 건너뛰기
  skip_requesting_account_id  = true        # 계정 ID 요청 건너뛰기

  # === 중요: LocalStack 엔드포인트 설정 ===
  # "endpoint" 인자는 provider 블록 내에서 직접 지원되지 않습니다. (사용자 오류 원인)
  # 대신, 각 서비스별로 엔드포인트를 지정해야 합니다.
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